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

	// TODO: Add parallel summarization from different algorithms. Loop over different algorithms and
	// get the best one.
	bearerToken := reddit.GetUserToken(redditUName, redditPwd)
	comments := reddit.LoadComments(subredditName, postId, sortingMethod, depth, bearerToken)

	config := getYamlConfig()
	summarizedText := make(map[string]string, len(config))

	emptyStr := ""
	var sr inference.SummarizerRequester
	for model_name := range config {
		con, ok := config[model_name].(map[string]interface{})
		if !ok {
			panic("could not find key in the model config")
		}

		summarysize := con["min_new_tokens"].(int)
		var totalMaxTokens int
		if con["model_type"] == "openai" {
			ors := inference.OpenAIRequestSummary{Paragraph: &emptyStr}
			sr = &ors
			totalMaxTokens = con["max_tokens"].(int) - summarysize - len(con["summary_suffix"].(string))
		} else if con["model_type"] == "py_service" {
			psrs := inference.PyServiceRequestSummary{Paragraph: &emptyStr}
			sr = &psrs
			totalMaxTokens = con["max_tokens"].(int) - con["min_new_tokens"].(int)
		} else {
			panic("Invalid Model Type")
		}

		// TODO: Add channel for parallel summarization
		summarizedText[model_name] = inference.GetSummarizedText(sr, comments, summarysize, totalMaxTokens,
			con["model_name"].(string))
		log.Println("Got summary from model" + model_name)
	}

	fmt.Println(summarizedText)

	c.JSON(200, strings.TrimSpace(summarizedText["text-davinci-003"]))
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
