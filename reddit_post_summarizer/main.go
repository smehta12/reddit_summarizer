package main

import (
	"fmt"
	"os"

	"./reddit"
)

func main() {
	bearerToken := "22403373384748-__MNvQmkd1ZdSsnQnzVNEpOauH9mDg" //reddit.GetUserToken("appsummrize", "app123**")
	subredditName := "personalfinance"
	postId := "11yqgoo"
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

}
