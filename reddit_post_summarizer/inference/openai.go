package inference

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const OPEN_AI_COMPLETION_ENDPOINT = "https://api.openai.com/v1/completions"
const OPEN_AI_EDIT_ENDPOINT = "https://api.openai.com/v1/edits"

const SUMMARY_SUFFIX = "\ntldr"
const FORMATTING_GPT_MODEL = "text-davinci-edit-001"
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

type OpenAIRequestSummary struct {
	Paragraph *string
	ModelName string
}

func (rs OpenAIRequestSummary) requestSummary() string {

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
	mp.Model = rs.ModelName
	mp.Prompt = *rs.Paragraph + SUMMARY_SUFFIX
	mp.MaxTokens = 500
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

	responseData := openaiPostRequest(bytes.NewBuffer(requestBody), OPEN_AI_COMPLETION_ENDPOINT, "summarization")

	var resJson SummaryResponse

	err = json.Unmarshal(responseData, &resJson)

	if err != nil {
		log.Println("Error in unmarshalling")
		log.Fatal(err)
	}

	return resJson.Choices[0].Text
}

func formatSummary(summary *string) string {
	cleanupInstruction := "cleanup this text"

	modelParameters := make(map[string]interface{})

	modelParameters["model"] = FORMATTING_GPT_MODEL
	modelParameters["input"] = *summary
	modelParameters["instruction"] = cleanupInstruction
	modelParameters["top_p"] = 1
	modelParameters["temperature"] = 0

	requestBody, err := json.Marshal(modelParameters)

	if err != nil {
		log.Println("Error while creating request body json in summary cleanup")
		panic(err)
	}

	responseData := openaiPostRequest(bytes.NewBuffer(requestBody), OPEN_AI_EDIT_ENDPOINT, "Summary Cleanup")

	var resJson SummaryCleanupResponse

	json.Unmarshal(responseData, &resJson)

	return resJson.Choices[0].Text
}

func openaiPostRequest(body io.Reader, requestURL string, intent string) []byte {
	// body: Post request body
	// requestURL: OpenAI endpoint
	// intent: Purpose of the request. It is used in logs

	request, err := http.NewRequest("POST", requestURL, body)

	if err != nil {
		log.Println("Error while creating new request for " + intent)
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+os.Getenv("OPEN_AI_BEARER"))

	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.Println("Error while getting response for" + intent)
		panic(err)
	}

	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		log.Println("Error in reading from response for " + intent)
		panic(err)
	}

	return responseData
}
