package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const (
	defaultBotToken  = "bot_token"
	defaultChannelID = 1 // chat_id
	defaultPort      = 8080
	maxFileSize      = 10 << 30 // 10 GB in bytes
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
			return
		}
		defer file.Close()

		if header.Size > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the maximum allowed size"})
			return
		}

		bot, err := tgbotapi.NewBotAPI(defaultBotToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Telegram bot"})
			return
		}

		fileUpload := tgbotapi.FileBytes{
			Name:  header.Filename,
			Bytes: make([]byte, header.Size), // pre-allocate memory
		}

		// Use io.LimitedReader to limit the number of bytes read from the file
		limitedReader := io.LimitedReader{R: file, N: maxFileSize}
		if _, err := io.ReadFull(&limitedReader, fileUpload.Bytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}

		message := tgbotapi.NewDocumentUpload(defaultChannelID, fileUpload)
		_, err = bot.Send(message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send file to Telegram"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File sent successfully!"})
	})

	r.Run(fmt.Sprintf(":%d", defaultPort))
}
