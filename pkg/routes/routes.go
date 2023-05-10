package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"tg-bot/pkg/db"
	"tg-bot/pkg/producer"
	"tg-bot/pkg/types"
	"tg-bot/pkg/utils"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func Start(salt, botName string, database db.IDataBase) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

		var params types.QueryParams

		if err := c.BindQuery(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		chatSign := utils.HashWithSalt(params.ChatId, salt)

		if chatSign != params.ChatSign {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userSign := utils.HashWithSalt(params.UserId, salt)

		if userSign != params.UserSign {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		chaId, err := strconv.ParseInt(params.ChatId, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "field 'chat_id' cant parse to integer",
			})
			return
		}
		userId, err := strconv.ParseInt(params.UserId, 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "field 'user_id' cant parse to integer",
			})
			return
		}

		if err := database.CreateAuthSign(c, int(userId), int(chaId)); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		botLink := fmt.Sprintf("https://t.me/%s", botName)
		c.Redirect(http.StatusSeeOther, botLink)
	}
}

func KickWorker(prod *producer.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg types.KickMessage

		if err := c.BindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := prod.SendKafkaMessage(msg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}
