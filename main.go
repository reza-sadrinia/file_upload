package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const (
	defaultBotToken  = "bot_token"
	defaultChannelID = 1 //chat_id
	defaultPort      = 8080
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer fileContent.Close()

		fileData, err := ioutil.ReadAll(fileContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		bot, err := tgbotapi.NewBotAPI(defaultBotToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fileUpload := tgbotapi.FileBytes{
			Name:  file.Filename,
			Bytes: fileData,
		}

		message := tgbotapi.NewDocumentUpload(defaultChannelID, fileUpload)

		_, err = bot.Send(message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File sent successfully!"})
	})

	r.Run(fmt.Sprintf(":%d", defaultPort))
}
