package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// Response models the top-level JSON
type Response struct {
	RemoteLoginURLs []struct {
		URL  string `json:"url"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	} `json:"remote_login_urls"`
}

func loginSSO(email string, password string) (sessionCookie *http.Cookie, err error) {
	// build multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// add post fields
	writer.WriteField("username", email)
	writer.WriteField("password", password)
	writer.WriteField("ajax", "1")
	writer.Close()

	req, _ := http.NewRequest("POST", LOGIN_URL, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	LoginResponse, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer LoginResponse.Body.Close()

	respBody, err := io.ReadAll(LoginResponse.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if LoginResponse.StatusCode != 200 {
		fmt.Printf("Status: %s\nResponse Body:\n%s", LoginResponse.Status, string(respBody))
		return
	}

	for _, c := range LoginResponse.Cookies() {
		if c.Name == "ssohls" {
			sessionCookie = c
		}
	}

	if sessionCookie == nil {
		fmt.Printf("Cookie not found for %s", email)
		return nil, nil
	}

	return sessionCookie, nil
	/*
		var resp Response
		if err := json.Unmarshal(respBody, &resp); err != nil {
			log.Fatalf("failed to parse JSON: %v", err)
		}

		var token string
		for _, entry := range resp.RemoteLoginURLs {
			if strings.Contains(entry.URL, "heise.de") {
				token = entry.Data.Token
				fmt.Println("Login Token:", token)
				loginShop(entry.URL, token)
				break
			}
		}
	*/
}
