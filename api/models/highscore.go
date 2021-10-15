package models

import (
	"crate-run-api/utils"
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

// HighScore represents the highest score of a single user on a single level.
type HighScore struct {
	ID             uint      `json:"-" binding:"-" gorm:"primary_key"`
	UserID         uint      `json:"-" binding:"-" gorm:"column:userId"`
	LevelNumber    uint      `json:"-" binding:"-" gorm:"column:levelNumber"`
	CompletionTime uint      `json:"completionTime" gorm:"column:completionTime"`
	Moves          uint      `json:"moves"`
	AchievedAt     time.Time `json:"-" binding:"-" gorm:"column:achievedAt"`
	User           User      `json:"user" binding:"-" `
}

// TableName is the name of the SQL table used by GORM to store the HighScore objects.
func (HighScore) TableName() string {
	return "highScore"
}

// UserTotalScore bundles a user with their global high score.
type UserTotalScore struct {
	User       User      `json:"user"`
	TotalScore uint      `json:"totalScore"`
	AchievedAt time.Time `json:"-"`
}

// GetGlobalHighScores returns a list of user's global scores in ascending order.
//
// The global score is aggregated for each user from their position on every level's top 10 high
// score list. For first place they get 10 points, for 10th place they get 1 point.
func GetGlobalHighScores(db *gorm.DB) ([]UserTotalScore, error) {
	levelHighScores := make(map[int][]HighScore)
	levelNumbers, err := getLevelsWithHighscore(db)

	if err != nil {
		return nil, fmt.Errorf("could not get global high scores: %w", err)
	}

	for _, levelNumber := range levelNumbers {
		scores, err := GetLevelHighScores(db, levelNumber)

		if err != nil {
			return nil, fmt.Errorf("could not get global high scores: %w", err)
		}

		levelHighScores[int(levelNumber)] = scores
	}

	return mergeLevelScoresToGlobal(levelHighScores), nil
}

// GetLevelHighScores returns a list of user high scores for a given level, in ascending order.
//
// A score is higher than another if it's number of moves are lower. When moves are the same the
// lower completion time wins. When both are same the older score wins.
func GetLevelHighScores(db *gorm.DB, levelNumber uint) ([]HighScore, error) {
	var highScores []HighScore

	err := db.Preload("User").
		Where(HighScore{LevelNumber: levelNumber}).
		Order("moves asc, completionTime asc, achievedAt asc").
		Find(&highScores).Error

	if err != nil {
		return nil, fmt.Errorf("could not get level high scores: %w", err)
	}

	return highScores, nil
}

func getLevelsWithHighscore(db *gorm.DB) ([]uint, error) {
	rows, err := db.Table("highScore").Select("DISTINCT levelNumber").Rows()

	if err != nil {
		return nil, fmt.Errorf("could not get levels: %w", err)
	}

	var levels []uint

	defer rows.Close()
	for rows.Next() {
		var result uint
		err := rows.Scan(&result)

		if err != nil {
			return nil, fmt.Errorf("could not get levels: %w", err)
		}

		levels = append(levels, result)
	}

	return levels, nil
}

func mergeLevelScoresToGlobal(levelHighScores map[int][]HighScore) []UserTotalScore {
	const maxScore int = 10

	globalScores := make(map[uint]UserTotalScore)

	for _, userHighScores := range levelHighScores {
		for i := 0; i < utils.Min(maxScore, len(userHighScores)); i++ {
			userId := userHighScores[i].User.ID

			_, exist := globalScores[userId]

			if !exist {
				globalScores[userId] = UserTotalScore{
					User:       userHighScores[i].User,
					TotalScore: 0,
					AchievedAt: userHighScores[i].AchievedAt,
				}
			}

			achievedAt := globalScores[userId].AchievedAt
			if achievedAt.Unix() < userHighScores[i].AchievedAt.Unix() {
				achievedAt = userHighScores[i].AchievedAt
			}

			globalScores[userId] = UserTotalScore{
				User:       globalScores[userId].User,
				TotalScore: globalScores[userId].TotalScore + uint(maxScore) - uint(i),
				AchievedAt: achievedAt,
			}
		}
	}

	ret := make([]UserTotalScore, 0, len(globalScores))

	for _, score := range globalScores {
		ret = append(ret, score)
	}

	sort.Slice(ret, func(a, b int) bool {
		if ret[a].TotalScore == ret[b].TotalScore {
			return ret[a].AchievedAt.Unix() < ret[b].AchievedAt.Unix()
		}

		return ret[a].TotalScore > ret[b].TotalScore
	})

	return ret
}
