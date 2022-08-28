package main

import (
	"course/internal/database"
	"course/internal/exercise/repository"
	"course/internal/exercise/usecase"
	"course/internal/middleware"
	userUc "course/internal/user/usecase"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db := database.NewDabataseConn(dsn)
	eRepo := repository.NewExerciseRepository(db)

	exerciseUcs := usecase.NewExerciseUsecase(eRepo)
	userUcs := userUc.NewUserUsecase(db)

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"message": "hello world",
		})
	})
	// exercise
	r.GET("/exercises/:id", middleware.WithAuthentication(userUcs), exerciseUcs.GetExercise)
	r.GET("/exercises/:id/scores", middleware.WithAuthentication(userUcs), exerciseUcs.CalculateScore)

	r.POST("/exercises", middleware.WithAuthentication(userUcs), exerciseUcs.CreateExercise)
	r.POST("/exercises/:id/questions", middleware.WithAuthentication(userUcs), exerciseUcs.CreateQuestion)
	r.POST("/exercises/:id/questions/:idb/answer", middleware.WithAuthentication(userUcs), exerciseUcs.CreateAnswer)

	// user
	r.POST("/register", userUcs.Register)
	r.POST("/login", userUcs.Login)
	r.Run(":1234")
}
