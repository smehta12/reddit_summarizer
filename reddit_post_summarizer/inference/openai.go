package inference

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

const OPEN_AI_COMPLETION_ENDPOINT = "https://api.openai.com/v1/completions"
const SUMMARY_SIZE = 500
const SUMMARY_SUFFIX = "\ntldr"
const MAX_TOKENS = 4096 - SUMMARY_SIZE - len(SUMMARY_SUFFIX)
const GPT_MODEL = "text-davinci-003"
const MODEL_TOKENIZER_ENCODING = "p50k_base"

type SummaryResponse struct {
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

type SummaryCleanupResponse struct {
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Text  string `json:"text"`
		Index int    `json:"index"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func GetSummarizedText(comments []string) string {
	cleanupComments(comments)
	summarizedText := summarizeText(comments)
	summarizedText = cleanupSummary(summarizedText)
	return summarizedText
}

func summarizeText(comments []string) string {

	// base case.
	if len(comments) == 1 && len(comments[0]) < SUMMARY_SIZE {
		return comments[0]
	}

	i := 0
	totalNumOfTokens := 0
	var paragraph string
	var summarizedText []string
	for i < len(comments) {
		numOfTokens := getNumberOfTokens(comments[i])

		// TODO: What if the numOfTokens in the sentence is more than 4096?

		totalNumOfTokens += numOfTokens
		if totalNumOfTokens <= MAX_TOKENS {
			paragraph += comments[i]
		} else {
			summarizedText = append(summarizedText, requestSummary(paragraph))
			i--
			totalNumOfTokens = 0
			paragraph = ""
		}
		i++

	}
	// for last paragraph
	summarizedText = append(summarizedText, requestSummary(paragraph))
	return summarizeText(summarizedText)
}

func requestSummary(paragraph string) string {

	type ModelParameters struct {
		Model           string  `json:"model"`
		Prompt          string  `json:"prompt"`
		MaxTokens       int     `json:"max_tokens"`
		Suffix          string  `json:"suffix"`
		Temperature     float32 `json:"temperature"`
		Top_p           float32 `json:"top_p"`
		N               int     `json:"n"`
		PresencePenalty float32 `json:"presence_penalty"`
	}

	var mp ModelParameters
	mp.Model = GPT_MODEL
	mp.Prompt = paragraph + SUMMARY_SUFFIX
	mp.MaxTokens = SUMMARY_SIZE
	mp.Suffix = ""
	mp.N = 1
	mp.Temperature = 0.7
	mp.Top_p = 1 //Default
	mp.PresencePenalty = 1

	requestBody, err := json.Marshal(mp)

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

	var resJson SummaryResponse

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
}

func cleanupSummary(summary string) string {
	cleanupInstruction := "cleanup this text"

	modelParameters := make(map[string]interface{})

	modelParameters["model"] = "text-davinci-edit-001"
	modelParameters["input"] = summary
	modelParameters["instruction"] = cleanupInstruction
	modelParameters["top_p"] = 1
	modelParameters["temperature"] = 0

	requestBody, err := json.Marshal(modelParameters)

	if err != nil {
		log.Println("Error while creating request body json in summary cleanup")
		panic(err)
	}

	request, err := http.NewRequest("POST", "https://api.openai.com/v1/edits", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Println("Error while creating new request in summary cleanup")
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+os.Getenv("OPEN_AI_BEARER"))

	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Println("Error while getting summary response in summary cleanup")
		panic(err)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		log.Println("Error in reading raw comments from response in summary cleanup")
		panic(err)
	}

	var resJson SummaryCleanupResponse

	json.Unmarshal(responseData, &resJson)

	return resJson.Choices[0].Text
}
