package inference

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

const OPEN_AI_COMPLETION_ENDPOINT = "https://api.openai.com/v1/completions"
const SUMMARY_SIZE = 500
const SUMMRY_SUFFIX = "\ntldr"
const MAX_TOKENS = 4096 - SUMMARY_SIZE - len(SUMMRY_SUFFIX)
const GPT_MODEL = "text-davinci-003"
const MODEL_TOKENIZER_ENCODING = "p50k_base"

type SummryResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		Logprobs     any    `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func GetSummarizedText(comments []string) string {
	cleanupComments(comments)
	summerizedText := summrizeText(comments)
	return summerizedText
}

func summrizeText(comments []string) string {

	// base case.
	if len(comments) == 1 && len(comments[0]) < SUMMARY_SIZE {
		return comments[0]
	}

	i := 0
	totalNumOfTokens := 0
	var paragraph string
	var summrizedText []string
	summerizedTextIdx := 0
	for i < len(comments) {
		numOfTokens := getNumberOfTokens(comments[i])

		// TODO: What if the numOfTokens in the sentence is more than 4096?

		totalNumOfTokens += numOfTokens
		if totalNumOfTokens <= MAX_TOKENS {
			paragraph += comments[i]
		} else {

			// Save the summries in the array
			summrizedText[summerizedTextIdx] = requestSummary(paragraph)
			summerizedTextIdx++
			i--
			totalNumOfTokens = 0
			paragraph = ""
		}
		i++
	}
	return summrizeText(summrizedText)
}

func requestSummary(paragraph string) string {

	type Body struct {
		Model     string `json:"model"`
		Prompt    string `json:"prompt"`
		MaxTokens int    `json:"max_tokens"`
		Suffix    string `json:"suffix"`
	}

	var b Body
	b.Model = GPT_MODEL
	b.Prompt = paragraph
	b.MaxTokens = SUMMARY_SIZE
	b.Suffix = SUMMRY_SUFFIX

	requestBody, err := json.Marshal(b)

	if err != nil {
		log.Println("Error while creating request body json")
		panic(err)
	}

	request, err := http.NewRequest("POST", OPEN_AI_COMPLETION_ENDPOINT, bytes.NewBuffer(requestBody))

	if err != nil {
		log.Println("Error while creating new request")
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+os.Getenv("OPEN_AI_BEARER"))

	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Println("Error while getting summary response")
		panic(err)
	}

	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		log.Println("Error in reading raw comments from response")
		log.Fatal(err)
	}
	fmt.Println("Response Status:", response.Status)

	var resJson SummryResponse

	err = json.Unmarshal(responseData, &resJson)

	if err != nil {
		log.Println("Error in unmarshalling")
		log.Fatal(err)
	}

	return resJson.Choices[0].Text
}

func cleanupComments(comments []string) {
	for i, s := range comments {
		comments[i] = strings.TrimSpace(s)
	}
}

func getNumberOfTokens(comment string) int {
	encoding, err := tiktoken.GetEncoding(MODEL_TOKENIZER_ENCODING)

	if err != nil {
		log.Println("Error when getting tokenizer encoding")
		log.Fatal(err)
	}

	return len(encoding.Encode(comment, nil, nil))

	// Using rule of thumb from https://platform.openai.com/docs/introduction/key-concepts
	// As a rough rule of thumb, 1 token is approximately 4 characters or 0.75 words for English text.

	// numOfChars := len(comment)
	// numOfTokens := math.Ceil(float64(numOfChars) / 4)
	// return int(numOfTokens)
}
