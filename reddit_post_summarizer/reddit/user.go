package reddit

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// TODO: getting {"error": "unsupported_grant_type"} when running code below
func GetUserToken(username string, password string) string {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("grant_type", "password")
	encodedData := data.Encode()

	r, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(encodedData))

	if err != nil {
		panic(err)
	}

	r.Header.Add("User-Agent", "posts_summarizer/0.0.1")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	r.SetBasicAuth(os.Getenv("REDDIT_APP_ID"), os.Getenv("REDDIT_SECRET"))

	client := &http.Client{}

	res, err := client.Do(r)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}
	log.Println("Response Status:", res.Status)

	m := make(map[string]string)
	err = json.Unmarshal(responseData, &m)

	if err != nil {
		panic(err)
	}

	return m["access_token"]
}
