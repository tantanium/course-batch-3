package model

import "time"

type Exercise struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions" gorm:"foreignKey:ExerciseID"`
}

// gorm:"foreignKey:ID
type Question struct {
	ID            int       `json:"id"`
	ExerciseID    int       `json:"exerciseid"`
	Body          string    `json:"body"`
	OptionA       string    `json:"option_a"`
	OptionB       string    `json:"option_b"`
	OptionC       string    `json:"option_c"`
	OptionD       string    `json:"option_d"`
	CorrectAnswer string    `json:"correct_answer"`
	Score         int       `json:"score" gorm:"default:10"`
	CreatorID     int       `json:"creator_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Answer struct {
	ID         int
	ExerciseID int
	QuestionID int
	UserID     int
	Answer     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
