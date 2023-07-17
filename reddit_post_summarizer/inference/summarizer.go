package inference

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

type SummarizerRequester interface {
	requestSummary(c chan string)
}

type PyServiceRequestSummary struct {
	Paragraph string
	ModelName string
}

type SummarizedTextReturn struct {
	ModelName string
	Text      string
}

func GetSummarizedText(sr SummarizerRequester, comments []string, summarySize int, totalMaxTokens int,
	model_name string, c chan SummarizedTextReturn) {
	cleanupComments(comments)
	summarizedText := summarizeTextRecursive(sr, comments, summarySize, totalMaxTokens, model_name)
	summarizedText = formatSummary(&summarizedText)
	s := SummarizedTextReturn{ModelName: model_name, Text: summarizedText}
	c <- s
}

// comments: data to summarize
// summarySize: How long summary can be.
// totalTokens: Total Max tokens. e.g. 4096 in text-davinci-003
func summarizeTextRecursive(sr SummarizerRequester, comments []string, summarySize int, totalMaxTokens int,
	model_name string) string {
	// base case.
	if len(comments) == 1 {
		return comments[0]
	}

	i := 0
	totalNumOfTokens := 0
	var paragraph string
	var summarizedText []string
	channel := make(chan string, 4)
	usedChannel := false

	for i < len(comments) {
		numOfTokens := getNumberOfTokens(comments[i])
		totalNumOfTokens += numOfTokens
		if totalNumOfTokens <= totalMaxTokens {
			paragraph += comments[i]
		} else if numOfTokens > totalMaxTokens {
			// TODO: What if the numOfTokens in the sentence is more than totalMaxTokens?
			// TODO: Devide the comment into small comments like less than max tokens allowed.
		} else {
			// Can't pass the pointer of paragraph here because sometimes it gets reset before subsequant
			// calls processes it.
			assignParagraph(sr, paragraph, model_name)
			go sr.requestSummary(channel)
			usedChannel = true
			i--
			totalNumOfTokens = 0
			paragraph = ""
		}
		i++
	}

	if usedChannel {
		summStr := <-channel
		summarizedText = append(summarizedText, summStr)
	}

	// for last paragraph
	assignParagraph(sr, paragraph, model_name)
	go sr.requestSummary(channel)
	summarizedText = append(summarizedText, <-channel)
	return summarizeTextRecursive(sr, summarizedText, summarySize, totalMaxTokens, model_name)
}

func assignParagraph(sr SummarizerRequester, paragraph string, model_name string) {
	switch s := sr.(type) {
	case *PyServiceRequestSummary:
		s.Paragraph = paragraph
		s.ModelName = model_name
	case *OpenAIRequestSummary:
		s.Paragraph = paragraph
		s.ModelName = model_name
	}
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

func (rs PyServiceRequestSummary) requestSummary(c chan string) {
	values := map[string]string{"model_name": rs.ModelName, "prompt": rs.Paragraph}

	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post("http://0.0.0.0:8000/summarize", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Println("Error when requsting summary from python service for model" + rs.ModelName)
	}

	// fmt.Println(resp)
	respData, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error in reading from response for" + rs.ModelName)
		panic(err)
	}

	c <- string(respData)
}
