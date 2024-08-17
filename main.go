package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FilenameInput struct {
	Filename string `json:"filename" binding:"required"`
}

func main() {
	fmt.Println("Go VM alert rule Validator")


	router := gin.Default()

	// Handle POST request to /upload
	router.POST("/upload", func(c *gin.Context) {
		var input FilenameInput
		// Bind the JSON payload to the FilenameInput struct
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// You can now use the `input.Filename` to process the file or any other logic
		fmt.Printf("Received filename: %s\n", input.Filename)

		// Simulate processing and return a success response
		c.JSON(http.StatusOK, gin.H{
			"message":  "File processed successfully",
			"filename": input.Filename,
		})
	})

	// Start the server on port 8080
	router.Run(":8080")

}