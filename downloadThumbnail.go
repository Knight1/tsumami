package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

func downloadThumbnail(magazine string, year string, issue string) {
	outputPath := "magazines/" + magazine + "/" + year + "/"
	req, _ := http.NewRequest("GET", "https://heise.cloudimg.io/v7/_www-heise-de_/select/thumbnail/"+magazine+"/"+year+"/"+issue+".jpg", nil)

	// If it does not exists abort
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

	filename := path.Base(req.RequestURI)

	respBody, _ = io.ReadAll(resp.Body)

	if len(respBody) > 0 {
		err = os.WriteFile(outputPath+filename, respBody, 0644)
		if err != nil {
			fmt.Println("Error writing file", err)
			return
		}
	}
}
