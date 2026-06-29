package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // Route test awal
    r.GET("/api/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":  "success",
            "message": "Backend AC Mobil Ku siap digunakan!",
        })
    })

    // Jalankan di port 8080
    r.Run(":8080")
}