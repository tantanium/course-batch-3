package exercise

import (
	"github.com/gin-gonic/gin"
)

type Usecase interface {
	GetExercise(c *gin.Context)
	CalculateScore(c *gin.Context)
	CreateExercise(c *gin.Context)
	CreateQuestion(c *gin.Context)
	CreateAnswer(c *gin.Context)
}
