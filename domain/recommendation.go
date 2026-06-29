package domain

import "context"

type Prediction struct {
	BengkelID        uint    `json:"bengkel_id"`
	PredictedRating  float64 `json:"predicted_rating"`
	Distance         float64 `json:"distance"`
	HaversinePassed  bool    `json:"haversine_passed"`
}

type RecommendationUsecase interface {
	GetRecommendations(ctx context.Context, userUID string, limit int) ([]Bengkel, error)
}
