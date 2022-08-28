package usecase

import (
	"course/internal/exercise"
	"course/internal/model"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ExerciseUsecase struct {
	eRepo exercise.Repository
}

func NewExerciseUsecase(eRepo exercise.Repository) exercise.Usecase {
	return &ExerciseUsecase{
		eRepo: eRepo,
	}
}

func (exerUsecase *ExerciseUsecase) GetExercise(c *gin.Context) {
	type QuestionResponse struct {
		ID        int       `json:"id"`
		Body      string    `json:"body"`
		OptionA   string    `json:"option_a"`
		OptionB   string    `json:"option_b"`
		OptionC   string    `json:"option_c"`
		OptionD   string    `json:"option_d"`
		Score     int       `json:"score"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	type ExerciseResponse struct {
		ID          int                `json:"id"`
		Title       string             `json:"title"`
		Description string             `json:"description"`
		Question    []QuestionResponse `json:"question"`
	}

	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, map[string]string{
			"message": "invalid exercise id",
		})
		return
	}

	var exercise *model.Exercise
	exercise, err = exerUsecase.eRepo.GetExercise(id)
	if err != nil {
		c.JSON(404, map[string]string{
			"message": "exercise not found",
		})
		return
	}

	var qData []QuestionResponse

	if len(exercise.Questions) > 0 {
		for _, value := range exercise.Questions {
			ques := QuestionResponse{
				ID:        value.ID,
				Body:      value.Body,
				OptionA:   value.OptionA,
				OptionB:   value.OptionB,
				OptionC:   value.OptionC,
				OptionD:   value.OptionD,
				Score:     value.Score,
				CreatedAt: value.CreatedAt,
				UpdatedAt: value.UpdatedAt,
			}
			qData = append(qData, ques)
		}
	}
	var d ExerciseResponse
	d.ID = exercise.ID
	d.Title = exercise.Title
	d.Description = exercise.Description
	d.Question = qData

	c.JSON(200, d)
}
func (exerUsecase *ExerciseUsecase) CalculateScore(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, map[string]string{
			"message": "invalid exercise id",
		})
		return
	}
	var exercise *model.Exercise
	exercise, err = exerUsecase.eRepo.GetExercise(id)
	if err != nil {
		c.JSON(404, map[string]string{
			"message": "exercise not found",
		})
		return
	}

	userID := int(c.Request.Context().Value("user_id").(float64))
	var answers *[]model.Answer
	answers, err = exerUsecase.eRepo.GetScore(id, userID)

	mapQA := make(map[int]model.Answer)
	for _, answer := range *answers {
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

func (exerUsecase *ExerciseUsecase) CreateExercise(c *gin.Context) {
	type ExerciseRequest struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	var createExerciseRequest ExerciseRequest
	if err := c.ShouldBind(&createExerciseRequest); err != nil {
		c.JSON(400, map[string]string{
			"message": "invalid input",
		})
		return
	}
	if createExerciseRequest.Title == "" {
		c.JSON(400, map[string]string{
			"message": "missing title",
		})
		return
	}
	if createExerciseRequest.Description == "" {
		c.JSON(400, map[string]string{
			"message": "missing description",
		})
		return
	}
	exer, err := exerUsecase.eRepo.CreateExercise(createExerciseRequest.Title, createExerciseRequest.Description)
	if err != nil {
		c.JSON(500, map[string]string{
			"message": "cannot create exercise",
		})
		return
	}

	var data ExerciseRequest
	data.ID = exer.ID
	data.Title = exer.Title
	data.Description = exer.Description

	c.JSON(200, data)
}

func (exerUsecase *ExerciseUsecase) CreateQuestion(c *gin.Context) {
	paramID := c.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		c.JSON(400, map[string]string{
			"message": "invalid exercise id",
		})
		return
	}
	var createQuestionRequest model.Question
	if err := c.ShouldBind(&createQuestionRequest); err != nil {
		c.JSON(400, map[string]string{
			"message": "invalid input",
		})
		return
	}
	// var exercise *model.Exercise
	_, err = exerUsecase.eRepo.GetExercise(id)
	if err != nil {
		c.JSON(404, map[string]string{
			"message": "exercise not found",
		})
		return
	}
	if createQuestionRequest.Body == "" {
		c.JSON(400, map[string]string{
			"message": "missing body",
		})
		return
	}
	if createQuestionRequest.OptionA == "" ||
		createQuestionRequest.OptionA == "" ||
		createQuestionRequest.OptionA == "" ||
		createQuestionRequest.OptionA == "" {
		c.JSON(400, map[string]string{
			"message": "missing option a/b/c/d",
		})
		return
	}
	if createQuestionRequest.CorrectAnswer == "" {
		c.JSON(400, map[string]string{
			"message": "missing correct answer",
		})
		return
	}
	r, _ := regexp.Compile("[abcd]")
	if !r.MatchString(string(createQuestionRequest.CorrectAnswer)) {
		c.JSON(400, map[string]string{
			"message": "invalid correct answer",
		})
		return
	}
	userID := int(c.Request.Context().Value("user_id").(float64))
	_, err = exerUsecase.eRepo.CreateQuestion(
		id,
		userID,
		createQuestionRequest.Body,
		createQuestionRequest.OptionA,
		createQuestionRequest.OptionB,
		createQuestionRequest.OptionC,
		createQuestionRequest.OptionD,
		createQuestionRequest.CorrectAnswer,
	)

	if err != nil {
		c.JSON(500, map[string]string{
			"message": "cannot create question",
		})
		return
	}

	c.JSON(200, map[string]string{
		"message": "success",
	})
}

func (exerUsecase *ExerciseUsecase) CreateAnswer(c *gin.Context) {
	ida, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, map[string]string{
			"message": "invalid exercise id",
		})
		return
	}
	idb, err := strconv.Atoi(c.Param("idb"))
	if err != nil {
		c.JSON(400, map[string]string{
			"message": "invalid question id",
		})
		return
	}

	var createQAnswerRequest model.Answer
	if err := c.ShouldBind(&createQAnswerRequest); err != nil {
		c.JSON(400, map[string]string{
			"message": "invalid input",
		})
		return
	}
	_, err = exerUsecase.eRepo.GetExercise(ida)
	if err != nil {
		c.JSON(404, map[string]string{
			"message": "exercise not found",
		})
		return
	}
	_, err = exerUsecase.eRepo.GetQuestion(idb)
	if err != nil {
		c.JSON(404, map[string]string{
			"message": "exercise not found",
		})
		return
	}
	r, _ := regexp.Compile("[abcd]")
	if !r.MatchString(string(createQAnswerRequest.Answer)) {
		c.JSON(400, map[string]string{
			"message": "missing correct answer",
		})
		return
	}
	userID := int(c.Request.Context().Value("user_id").(float64))
	print(userID)
	_, err = exerUsecase.eRepo.CreateAnswer(ida, idb, userID, createQAnswerRequest.Answer)
	if err != nil {
		c.JSON(500, map[string]string{
			"message": "cannot create answer / answer is exist",
		})
		return
	}

	c.JSON(200, map[string]string{
		"message": "success",
	})

}
