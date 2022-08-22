package main

import (
	"course/internal/database"
	"course/internal/exercise/usecase"
	"course/internal/middleware"
	userUc "course/internal/user/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db := database.NewDabataseConn()
	exerciseUcs := usecase.NewExerciseUsecase(db)
	userUcs := userUc.NewUserUsecase(db)
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"message": "hello world",
		})
	})
	// exercise
	r.GET("/exercises/:id", middleware.WithAuthentication(userUcs), exerciseUcs.GetExercise)
	r.GET("/exercises/:id/scores", middleware.WithAuthentication(userUcs), exerciseUcs.CalculateScore)

	// user
	r.POST("/register", userUcs.Register)
	r.POST("/login", userUcs.Login)
	r.Run(":1234")
}
