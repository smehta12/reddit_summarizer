package main

import (
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smehta12/reddit_summarizer/inference"
	"github.com/smehta12/reddit_summarizer/reddit"
	"gopkg.in/yaml.v3"
)

func getSummary(c *gin.Context) {
	var subredditName string
	var postId string

	redditURL := c.Query("reddit_url")
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

	// TODO: Add parallel summarization from different algorithms. Loop over different algorithms and
	// get the best one.
	bearerToken := reddit.GetUserToken(redditUName, redditPwd)
	comments := reddit.LoadComments(subredditName, postId, sortingMethod, depth, bearerToken)
	var sr inference.SummarizerRequester
	emptyStr := ""
	ors := inference.OpenAIRequestSummary{Paragraph: &emptyStr}
	sr = &ors

	config := getYamlConfig()
	con, ok := config["text-davinci-003"].(map[string]interface{})
	if !ok {
		panic("could not find key in the model config")
	}
	summarysize := con["summary_size"].(int)
	totalMaxTokens := con["max_tokens"].(int) - summarysize - len(con["summary_suffix"].(string))
	summarizedText := inference.GetSummarizedText(sr, comments, summarysize, totalMaxTokens)

	c.JSON(200, strings.TrimSpace(summarizedText))
}

func getYamlConfig() map[string]interface{} {

	yamlFile, err := ioutil.ReadFile("../model_configs.yaml")
	if err != nil {
		panic(err)
	}

	// Unmarshal the YAML file into a map.
	config := make(map[string]interface{})
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return config
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
