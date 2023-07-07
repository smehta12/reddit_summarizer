package inference

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/smehta12/reddit_summarizer/utils"
)

func GetSummaryRankings(redditDocId string, comments []string, summaries map[string]string) []float64 {
	doc := strings.Join(comments, " ")
	scores := requestRankings(redditDocId, &doc, utils.GetMapValues(summaries))

	return scores
}

func GetHighestRankedSummary(redditDocId string, comments []string, summaries map[string]string) string {
	var highestRankedSummary int

	scores := GetSummaryRankings(redditDocId, comments, summaries)

	var i int
	var highestAlgoScore string

	for algo := range summaries {
		if scores[i] > float64(highestRankedSummary) {
			highestAlgoScore = algo
		}
	}

	return summaries[highestAlgoScore]
}

func requestRankings(docId string, doc *string, summaries []string) []float64 {
	values := map[string]interface{}{"doc_id": docId, "doc": *doc, "summaries": summaries}

	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post("http://0.0.0.0:8001/rank", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		log.Println("Error when requsting summary from ranking service for" + docId)
	}

	// fmt.Println(resp)
	respData, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error in reading from response for ranking service with id:" + docId)
		panic(err)
	}

	println(respData)

	temp := [1]float64{2.0}

	return temp[:]
}
