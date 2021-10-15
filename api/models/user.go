package models

import (
	"crate-run-api/utils"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// User represents a single registered player.
type User struct {
	ID           uint        `json:"id" binding:"required" gorm:"primary_key"`
	Username     string      `json:"username" binding:"required"`
	PasswordHash string      `json:"-" gorm:"column:passwordHash"`
	Nickname     string      `json:"nickname" binding:"required"`
	HighScores   []HighScore `json:"-" gorm:"foreignkey:UserID;not null"`
}

var (
	// ErrUsernameTaken is returned when a user with the given username already exists.
	ErrUsernameTaken = errors.New("username taken")

	// ErrInvalidLevel is returned when the given level number is too low or too high.
	ErrInvalidLevel = errors.New("invalid level number")

	// ErrUserAlreadyExists is returned when a given user is already in the database.
	ErrUserAlreadyExists = errors.New("user already exists")
)

// TableName is the name of the SQL table used by GORM to store the User objects
func (User) TableName() string {
	return "user"
}

// Register tries to save a User to the database with the given password. If the user already
// has a database ID it returns ErrUserAlreadyExists. If the username field is empty it
// generates a random username before saving. If the username is taken it returns
// ErrUsernameTaken.
func (u *User) Register(db *gorm.DB, password string) error {
	if u.ID > 0 {
		return ErrUserAlreadyExists
	}

	if u.Username == "" {
		u.Username = utils.GenerateUsername()
	}

	var count int64
	if err := db.Model(u).Where(u).Count(&count).Error; err != nil {
		return fmt.Errorf("could not query users: %w", err)
	}

	if count > 0 {
		return ErrUsernameTaken
	}

	u.PasswordHash = utils.HashPassword(password)

	if err := db.Create(u).Error; err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	return nil
}

// CheckPassword returns whether the given password matches the stored hash or not.
func (u *User) CheckPassword(db *gorm.DB, password string) (bool, error) {
	u.PasswordHash = utils.HashPassword(password)

	err := db.Where(u).First(u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("could not check password: %w", err)
	}

	return true, nil
}

// ChangeNickname updates the user's nickname in the database.
func (u *User) ChangeNickname(db *gorm.DB, nickname string) error {
	err := db.Model(u).Update(User{Nickname: nickname}).Error

	if err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}

	return nil
}

// GetHighScores returns all the high scores for the user.
func (u *User) GetHighScores(db *gorm.DB) ([]HighScore, error) {
	var highScores []HighScore

	err := db.Model(u).Association("HighScores").Find(&highScores).Error

	if err != nil {
		return nil, fmt.Errorf("could not get user high scores: %w", err)
	}

	return highScores, nil
}

// RegisterHighScore saves a new high score for the user. If no high score existed for the given
// level it creates a new entry in the database. If a high score existed, but it's considered a
// lower score it updates it. If the old score is higher it's a no-op.
func (u *User) RegisterHighScore(db *gorm.DB, registeredScore HighScore) error {
	const maxLevel uint = 50

	if registeredScore.LevelNumber < 1 || registeredScore.LevelNumber > maxLevel {
		return ErrInvalidLevel
	}

	registeredScore.UserID = u.ID
	registeredScore.AchievedAt = time.Now()

	var highScores []HighScore
	err := db.Where(HighScore{
		UserID:      u.ID,
		LevelNumber: registeredScore.LevelNumber,
	}).Find(&highScores).Error

	if err != nil {
		return fmt.Errorf("could not register high score: %w", err)
	}

	err = nil
	if len(highScores) == 0 {
		err = db.Create(&registeredScore).Error
	} else if highScores[0].Moves > registeredScore.Moves ||
		(highScores[0].Moves == registeredScore.Moves &&
			highScores[0].CompletionTime > registeredScore.CompletionTime) {

		err = db.Model(&HighScore{}).
			Where(HighScore{ID: highScores[0].ID}).Updates(registeredScore).Error
	}

	if err != nil {
		return fmt.Errorf("could not register high score: %w", err)
	}

	return nil
}
