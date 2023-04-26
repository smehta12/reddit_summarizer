package main

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smehta12/reddit_summarizer/inference"
	"github.com/smehta12/reddit_summarizer/reddit"
)

//TODO: Change to proper package name

func getSummary(c *gin.Context) {
	var subredditName string
	var postId string

	redditURL := c.Query("reddit_url")
	// TODO: Get from auth method
	redditUName := c.Query("reddit_username")
	redditPwd := c.Query("reddit_password")

	if redditURL == "" || redditUName == "" || redditPwd == "" {
		c.JSON(400, gin.H{"error": "missing parameters"})
	}

	parsedUrl, err := url.Parse(redditURL)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid Reddit URL"})
		return
	}

	uri := parsedUrl.Path
	splittedURI := strings.Split(uri, "/")

	if splittedURI[1] != "r" {
		c.JSON(400, gin.H{"error": "invalid Reddit URL"})
		return
	} else {
		subredditName = splittedURI[2]
	}

	if splittedURI[3] != "comments" {
		c.JSON(400, gin.H{"error": "invalid Reddit URL"})
		return
	} else {
		postId = splittedURI[4]
	}

	sortingMethod := "top"
	depth := 1

	bearerToken := reddit.GetUserToken(redditUName, redditPwd)
	comments := reddit.LoadComments(subredditName, postId, sortingMethod, depth, bearerToken)
	summarizedText := inference.GetSummarizedText(comments)

	c.JSON(200, strings.TrimSpace(summarizedText))
}

func main() {

	// Generate gin code endpoint for post request

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Reddit Summarizer",
		})
	})
	router.GET("/summary", getSummary)

	router.Run() // listen and serve on 0.0.0.0:8080
}
