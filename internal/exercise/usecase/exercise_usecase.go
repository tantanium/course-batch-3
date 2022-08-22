package usecase

import (
	"course/internal/domain"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExerciseUsecase struct {
	db *gorm.DB
}

func NewExerciseUsecase(db *gorm.DB) *ExerciseUsecase {
	return &ExerciseUsecase{
		db: db,
	}
}

func (exerUsecase ExerciseUsecase) GetExercise(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, map[string]string{
			"message": "invalid exercise id",
		})
		return
	}

	var exercise domain.Exercise
	err = exerUsecase.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(404, map[string]string{
			"message": "exercise not found",
		})
		return
	}
	c.JSON(200, exercise)
}

func (exerUsecase ExerciseUsecase) CalculateScore(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, map[string]string{
			"message": "invalid exercise id",
		})
		return
	}

	var exercise domain.Exercise
	err = exerUsecase.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		c.JSON(404, map[string]string{
			"message": "exercise not found",
		})
		return
	}

	userID := int(c.Request.Context().Value("user_id").(float64))
	var answers []domain.Answer
	err = exerUsecase.db.Where("exercise_id = ? AND user_id = ?", id, userID).Find(&answers).Error
	if err != nil || len(answers) == 0 {
		c.JSON(200, map[string]interface{}{
			"score": 0,
		})
		return
	}
	mapQA := make(map[int]domain.Answer)
	for _, answer := range answers {
		mapQA[answer.QuestionID] = answer
	}

	var score ScoreCount
	wg := new(sync.WaitGroup)
	for _, question := range exercise.Questions {
		newQuestion := question
		wg.Add(1)
		go func() {
			defer wg.Done()
			if strings.EqualFold(newQuestion.CorrectAnswer, mapQA[newQuestion.ID].Answer) {
				score.Inc(newQuestion.Score)
			}
		}()
	}
	wg.Wait()
	c.JSON(200, map[string]interface{}{
		"score": score.score,
	})
}

type ScoreCount struct {
	score int
	mu    sync.Mutex
}

func (sc *ScoreCount) Inc(value int) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.score += value
}
