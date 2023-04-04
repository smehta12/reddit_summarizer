package reddit

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// TODO: getting {"error": "unsupported_grant_type"} when running code below
func GetUserToken(username string, password string) string {
	type User struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		GrantType string `json:"grant_type"`
	}

	u := User{
		Username:  username,
		Password:  password,
		GrantType: "password",
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(u)

	if err != nil {
		log.Fatal(err)
	}

	// var jsonData = []byte(`{
	// 	"username": "appsummrize",
	// 	"password": "app123**"
	// 	"grant_type": "password"
	// }`)

	r, err := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", &buf) //bytes.NewBuffer(jsonData)) //&buf)

	if err != nil {
		panic(err)
	}

	r.Header.Add("User-Agent", "posts_summarizer/0.0.1")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	app_id := "ot7-YcZcb5b8DzATpM9QPg"
	secret := "0FMS283Nyzb32Csmb05FRsgoaXKftg"
	// auth := app_id + ":" + secret
	// encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	// r.Header.Set("Authorization", "Basic "+encodedAuth)

	r.SetBasicAuth(app_id, secret)

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
	log.Println(string(responseData))

	return string(responseData)
}
