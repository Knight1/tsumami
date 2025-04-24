package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func getAccessRights(sessionCookie *http.Cookie) {
	req, _ := http.NewRequest("GET", "https://www.heise.de/api/accountservice/subscriptions/access-rights", nil)
	req.AddCookie(sessionCookie)

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Headers:\n%s", resp.Header)
		fmt.Printf("Body:\n%s", string(respBody))
		return
	}
	//fmt.Printf("Status: %s\n", resp.Status)
	//fmt.Printf("Headers:\n%s", resp.Header)
	//fmt.Printf("Body:\n%s", string(respBody))
}
