package main

import (
	"fmt"
	"io/ioutil"
	"log"
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

	bearerToken := reddit.GetUserToken(redditUName, redditPwd)
	comments := reddit.LoadComments(subredditName, postId, sortingMethod, depth, bearerToken)

	config := getYamlConfig()
	summaries := make(map[string]string, len(config))

	log.Println("Starting to get summaries")
	emptyStr := ""
	var sr inference.SummarizerRequester
	channel := make(chan inference.SummarizedTextReturn, len(config))
	for model_name := range config {
		log.Println("Getting summary from " + model_name)
		con, ok := config[model_name].(map[string]interface{})
		if !ok {
			panic("could not find key in the model config")
		}

		summarysize := con["min_new_tokens"].(int)
		var totalMaxTokens int
		if con["model_type"] == "openai" {
			ors := inference.OpenAIRequestSummary{Paragraph: emptyStr}
			sr = &ors
			totalMaxTokens = con["max_tokens"].(int) - summarysize - len(con["summary_suffix"].(string))
		} else if con["model_type"] == "py_service" {
			psrs := inference.PyServiceRequestSummary{Paragraph: emptyStr}
			sr = &psrs
			totalMaxTokens = con["max_tokens"].(int) - con["min_new_tokens"].(int)
		} else {
			panic("Invalid Model Type")
		}

		go inference.GetSummarizedText(sr, comments, summarysize, totalMaxTokens, con["model_name"].(string), channel)
	}

	for len(summaries) != len(config) {
		s := <-channel
		summaries[s.ModelName] = strings.TrimSpace(s.Text)
	}

	fmt.Println(summaries)
	log.Println("Completed gettting summaries")

	// Rank summarization
	highestRankedSummary := inference.GetHighestRankedSummary(postId, comments, summaries)

	log.Println("Sending highest ranked summary")
	c.JSON(200, highestRankedSummary)
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
