package app

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"crate-run-api/models"

	"github.com/gin-gonic/gin"
)

func (a Application) getGlobalHighScores(c *gin.Context) {
	scores, err := models.GetGlobalHighScores(a.DB)

	if err != nil {
		log.Printf("Handler error: %s", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, scores)
}

func (a Application) getLevelHighScores(c *gin.Context) {
	levelNumber, err := strconv.Atoi(c.Param("number"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	scores, err := models.GetLevelHighScores(a.DB, uint(levelNumber))

	if err != nil {
		log.Printf("Handler error: %s", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, scores)
}

func (a Application) getUserHighScores(c *gin.Context) {
	user := getUser(c)

	highScores, err := user.GetHighScores(a.DB)

	if err != nil {
		log.Printf("Handler error: %s", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	ret := make(map[uint]models.HighScore)

	for _, highScore := range highScores {
		ret[highScore.LevelNumber] = highScore
	}

	c.JSON(http.StatusOK, ret)
}

func (a Application) registerHighScore(c *gin.Context) {
	user := getUser(c)
	var registeredScore models.HighScore

	if c.ShouldBindJSON(&registeredScore) != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	levelNumber, err := strconv.Atoi(c.Param("number"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	registeredScore.LevelNumber = uint(levelNumber)

	err = user.RegisterHighScore(a.DB, registeredScore)

	if err != nil {
		if errors.Is(err, models.ErrInvalidLevel) {
			c.Status(http.StatusBadRequest)
		} else {
			log.Printf("Handler error: %s", err.Error())
			c.Status(http.StatusInternalServerError)
		}

		return
	}

	c.Status(http.StatusOK)
}
