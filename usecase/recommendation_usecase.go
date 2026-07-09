package usecase

import (
	"context"
	"math"
	"sort"

	"be-ac-mobil-ku/domain"
)

type recommendationUsecase struct {
	userRepo    domain.UserRepository
	bengkelRepo domain.BengkelRepository
	ratingRepo  domain.RatingRepository
	layananRepo domain.LayananRepository
}

func NewRecommendationUsecase(
	userRepo domain.UserRepository,
	bengkelRepo domain.BengkelRepository,
	ratingRepo domain.RatingRepository,
	layananRepo domain.LayananRepository,
) domain.RecommendationUsecase {
	return &recommendationUsecase{
		userRepo:    userRepo,
		bengkelRepo: bengkelRepo,
		ratingRepo:  ratingRepo,
		layananRepo: layananRepo,
	}
}

// Calculate Haversine Distance in Kilometers
func CalculateHaversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0 // kilometers

	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	radLat1 := lat1 * math.Pi / 180.0
	radLat2 := lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(radLat1)*math.Cos(radLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func (u *recommendationUsecase) GetRecommendations(ctx context.Context, userUID string, limit int) ([]domain.Bengkel, error) {
	// 1. Get current user profile and coordinates
	currentUser, err := u.userRepo.GetByUID(ctx, userUID)
	if err != nil {
		return nil, err
	}

	// 2. Get all bengkels
	allBengkels, err := u.bengkelRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Apply Constraint Filtering: Calculate distances & populate average ratings
	var candidateBengkels []domain.Bengkel
	maxDistanceKm := 50.0 // Constraint parameter: 50 km radius limit

	for _, b := range allBengkels {
		dist := 0.0
		if currentUser.Latitude != 0.0 && currentUser.Longitude != 0.0 {
			dist = CalculateHaversine(currentUser.Latitude, currentUser.Longitude, b.Latitude, b.Longitude)
		}
		b.Distance = dist

		// Populate ratings
		avgKualitas, avgHarga, _ := u.ratingRepo.GetAverageRatingByBengkelID(ctx, b.ID)
		b.AvgRatingKualitas = avgKualitas
		b.AvgRatingHarga = avgHarga
		if avgKualitas > 0 || avgHarga > 0 {
			b.AvgRatingKeseluruhan = (avgKualitas + avgHarga) / 2.0
		} else {
			b.AvgRatingKeseluruhan = 0.0
		}

		// Constraint Filtering: only include bengkels within spatial range,
		// unless user location is not set (0.0, 0.0) in which case we include all
		if currentUser.Latitude == 0.0 && currentUser.Longitude == 0.0 {
			candidateBengkels = append(candidateBengkels, b)
		} else if dist <= maxDistanceKm {
			candidateBengkels = append(candidateBengkels, b)
		}
	}

	// Fetch min harga for candidate bengkels
	for i := range candidateBengkels {
		b := &candidateBengkels[i]
		layanans, err := u.layananRepo.GetByBengkelID(ctx, b.ID)
		if err == nil && len(layanans) > 0 {
			minHarga := layanans[0].EstimasiHarga
			for _, l := range layanans {
				if l.EstimasiHarga < minHarga {
					minHarga = l.EstimasiHarga
				}
			}
			b.MinHarga = minHarga
		}
	}

	// 4. Collaborative Filtering (User-Based)
	// Fetch all ratings to construct User-Item Matrix
	allRatings, err := u.ratingRepo.GetAll(ctx)
	if err != nil || len(allRatings) == 0 {
		// Fallback: Sort candidates by distance and average ratings if no rating matrix exists
		sort.Slice(candidateBengkels, func(i, j int) bool {
			if candidateBengkels[i].AvgRatingKeseluruhan == candidateBengkels[j].AvgRatingKeseluruhan {
				return candidateBengkels[i].Distance < candidateBengkels[j].Distance
			}
			return candidateBengkels[i].AvgRatingKeseluruhan > candidateBengkels[j].AvgRatingKeseluruhan
		})
		if len(candidateBengkels) > limit {
			return candidateBengkels[:limit], nil
		}
		return candidateBengkels, nil
	}

	// Construct Matrix: UserUID -> BengkelID -> Rating Value (average of Kualitas & Harga)
	userRatings := make(map[string]map[uint]float64)
	userAvgRating := make(map[string]float64)
	userRatingCount := make(map[string]int)

	for _, r := range allRatings {
		if _, ok := userRatings[r.PelangganID]; !ok {
			userRatings[r.PelangganID] = make(map[uint]float64)
		}
		val := float64(r.RatingKualitas+r.RatingHarga) / 2.0
		userRatings[r.PelangganID][r.BengkelID] = val
		userAvgRating[r.PelangganID] += val
		userRatingCount[r.PelangganID]++
	}

	// Calculate averages
	for uID := range userAvgRating {
		if count := userRatingCount[uID]; count > 0 {
			userAvgRating[uID] = userAvgRating[uID] / float64(count)
		}
	}

	// Calculate Cosine Similarity between target user (userUID) and all other users
	similarities := make(map[string]float64)
	targetUserRatings := userRatings[userUID]
	targetAvg := userAvgRating[userUID]

	if len(targetUserRatings) > 0 {
		for otherUID, otherRatings := range userRatings {
			if otherUID == userUID {
				continue
			}

			// Cosine Similarity calculation
			var dotProduct, targetNorm, otherNorm float64
			hasCommon := false

			for bID, targetVal := range targetUserRatings {
				otherVal, exists := otherRatings[bID]
				if exists {
					hasCommon = true
					dotProduct += (targetVal - targetAvg) * (otherVal - userAvgRating[otherUID])
				}
				targetNorm += (targetVal - targetAvg) * (targetVal - targetAvg)
			}

			for _, otherVal := range otherRatings {
				otherNorm += (otherVal - userAvgRating[otherUID]) * (otherVal - userAvgRating[otherUID])
			}

			if hasCommon && targetNorm > 0 && otherNorm > 0 {
				similarities[otherUID] = dotProduct / (math.Sqrt(targetNorm) * math.Sqrt(otherNorm))
			}
		}
	}

	// Predict ratings for candidates
	type scoredBengkel struct {
		bengkel         domain.Bengkel
		predictedRating float64
	}
	var scoredList []scoredBengkel

	for _, b := range candidateBengkels {
		pred := 0.0
		ratedByTarget := false
		if targetUserRatings != nil {
			_, ratedByTarget = targetUserRatings[b.ID]
		}

		if ratedByTarget {
			// If already rated, we can just use the rating value as predicted
			pred = targetUserRatings[b.ID]
		} else if len(similarities) > 0 {
			// Predict using User-Based CF formula
			var weightedSum, simSum float64
			for otherUID, sim := range similarities {
				if rVal, exists := userRatings[otherUID][b.ID]; exists {
					weightedSum += sim * (rVal - userAvgRating[otherUID])
					simSum += math.Abs(sim)
				}
			}

			if simSum > 0 {
				pred = targetAvg + (weightedSum / simSum)
			} else {
				pred = b.AvgRatingKeseluruhan
			}
		} else {
			pred = b.AvgRatingKeseluruhan
		}

		// Ensure prediction is within 1-5 range
		if pred < 1.0 && pred > 0 {
			pred = 1.0
		} else if pred > 5.0 {
			pred = 5.0
		}

		scoredList = append(scoredList, scoredBengkel{
			bengkel:         b,
			predictedRating: pred,
		})
	}

	// Sort scored list by predicted rating descending, then by distance ascending
	sort.Slice(scoredList, func(i, j int) bool {
		if math.Abs(scoredList[i].predictedRating-scoredList[j].predictedRating) < 0.001 {
			return scoredList[i].bengkel.Distance < scoredList[j].bengkel.Distance
		}
		return scoredList[i].predictedRating > scoredList[j].predictedRating
	})

	// Format output
	var result []domain.Bengkel
	for i, sb := range scoredList {
		if i >= limit {
			break
		}
		b := sb.bengkel
		// We can inject the predicted rating value back to average rating or custom field if needed
		if sb.predictedRating > 0 {
			b.AvgRatingKeseluruhan = sb.predictedRating
		}
		result = append(result, b)
	}

	return result, nil
}
