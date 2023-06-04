package inference

import (
	"log"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

type SummarizerRequester interface {
	requestSummary() string
}

type PyServiceRequestSummary struct {
	Paragraph string
}

func GetSummarizedText(sr SummarizerRequester, comments []string, summarySize int, totalMaxTokens int) string {
	cleanupComments(comments)
	summarizedText := summarizeTextRecursive(sr, comments, summarySize, totalMaxTokens)
	summarizedText = formatSummary(summarizedText)
	return summarizedText
}

// comments: data to summarize
// summarySize: How long summary can be.
// totalTokens: Total Max tokens. e.g. 4096 in text-davinci-003
func summarizeTextRecursive(sr SummarizerRequester, comments []string, summarySize int, totalMaxTokens int) string {
	// base case.
	if len(comments) == 1 && len(comments[0]) < summarySize {
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
		if totalNumOfTokens <= totalMaxTokens {
			paragraph += comments[i]
		} else {
			assignParagraph(sr, paragraph)
			summarizedText = append(summarizedText, sr.requestSummary())
			i--
			totalNumOfTokens = 0
			paragraph = ""
		}
		i++
	}

	// for last paragraph
	assignParagraph(sr, paragraph)
	summarizedText = append(summarizedText, sr.requestSummary())
	return summarizeTextRecursive(sr, summarizedText, summarySize, totalMaxTokens)
}

func assignParagraph(sr SummarizerRequester, paragraph string) {
	switch s := sr.(type) {
	case *PyServiceRequestSummary:
		s.Paragraph = paragraph
	case *OpenAIRequestSummary:
		s.Paragraph = paragraph
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

func (rs PyServiceRequestSummary) requestSummary() string {
	return ""
}
