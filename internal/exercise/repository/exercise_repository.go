package repository

import (
	"course/internal/exercise"
	"course/internal/model"

	"gorm.io/gorm"
)

type ExerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) exercise.Repository {
	return &ExerciseRepository{
		db: db,
	}
}

func (e *ExerciseRepository) GetExercise(id int) (*model.Exercise, error) {
	var ex model.Exercise
	err := e.db.Debug().Where("id = ?", id).Preload("Questions").Take(&ex).Error
	if err != nil {
		return nil, err
	}
	return &ex, err
}

func (e *ExerciseRepository) GetQuestion(id int) (*model.Question, error) {
	var ex model.Question
	err := e.db.Debug().Where("id = ?", id).Take(&ex).Error
	if err != nil {
		return nil, err
	}
	return &ex, err
}

func (e *ExerciseRepository) GetScore(id, uid int) (*[]model.Answer, error) {
	var an []model.Answer
	err := e.db.Where("exercise_id = ? AND user_id = ?", id, uid).Find(&an).Error
	if err != nil {
		return nil, err
	}
	return &an, err
}

func (e *ExerciseRepository) CreateExercise(title, description string) (*model.Exercise, error) {
	var ex model.Exercise
	ex.Title = title
	ex.Description = description
	err := e.db.Debug().Create(&ex).Error
	if err != nil {
		return nil, err
	}
	return &ex, err
}

func (e *ExerciseRepository) CreateQuestion(idexer, userid int, body, optia, optib, optic, optid, answer string) (*model.Question, error) {
	ex := model.Question{
		ExerciseID:    idexer,
		Body:          body,
		OptionA:       optia,
		OptionB:       optib,
		OptionC:       optic,
		OptionD:       optid,
		CorrectAnswer: answer,
		CreatorID:     userid,
	}
	err := e.db.Create(&ex).Error
	if err != nil {
		return nil, err
	}
	return &ex, err
}

func (e *ExerciseRepository) CreateAnswer(idexer, idques, iduser int, answer string) (*model.Answer, error) {
	var ex model.Answer
	ex.Answer = answer
	ex.ExerciseID = idexer
	ex.QuestionID = idques
	ex.UserID = iduser
	err := e.db.Debug().Create(&ex).Error
	if err != nil {
		return nil, err
	}
	return &ex, err
}
