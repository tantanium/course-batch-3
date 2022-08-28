package exercise

import (
	"course/internal/model"
)

type Repository interface {
	GetExercise(int) (*model.Exercise, error)
	GetQuestion(id int) (*model.Question, error)
	GetScore(id, uid int) (*[]model.Answer, error)
	CreateExercise(title, description string) (*model.Exercise, error)
	CreateQuestion(idexer, userid int, body, optia, optib, optic, optid, answer string) (*model.Question, error)
	CreateAnswer(idexer, idques, iduser int, answer string) (*model.Answer, error)
}
