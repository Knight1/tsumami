package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

func loginShop(url string, token string) {
	// build multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// add post fields
	writer.WriteField("token", token)

	// finalize
	writer.Close()

	// create request
	req, _ := http.NewRequest("POST", url, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Status: %s\nresp Body:\n%s", resp.Status, string(respBody))
		return
	}
	fmt.Printf("Headers:\n%s", resp.Header)
	fmt.Printf("Body:\n%s", string(respBody))

	var sessionCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if len(c.Name) >= 8 && c.Name[:8] == "session-" {
			fmt.Println("Session cookie:", c)
			sessionCookie = c
		}
	}
	fmt.Println(sessionCookie)
	downloadMagazine(sessionCookie, "ct", "2025", "9")
}
