package main

import (
	"fmt"
	"os"

	"./inference" //TODO: Change to proper package name
	"./reddit"
)

func main() {
	bearerToken := reddit.GetUserToken(os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	subredditName := "personalfinance"
	postId := "12bpmx3"
	sortingMethod := "top"
	depth := 1

	comments := reddit.LoadComments(subredditName, postId, sortingMethod, depth, bearerToken)

	f, err := os.Create("comments.txt")
	if err != nil {
		fmt.Println("Issue in file creation")
	}

	defer f.Close()

	for _, val := range comments {
		f.WriteString(val)
		// f.WriteString("\n--------\n")
	}
	inference.GetSummarizedText(comments)
}
