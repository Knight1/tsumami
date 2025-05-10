package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func downloadThumbnail(magazine string, year string, issue string) {
	outputPath := "magazines/" + magazine + "/" + year + "/"
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("GET", "https://heise.cloudimg.io/v7/_www-heise-de_/select/thumbnail/"+magazine+"/"+year+"/"+issue+".jpg", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// If it does not exists abort
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Headers:\n%s", resp.Header)
		fmt.Printf("Body:\n%s", string(respBody))
		return
	}

	filename := path.Base(req.URL.Path)

	// Validate the filename to prevent directory traversal or invalid characters
	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") || strings.Contains(filename, "..") {
		fmt.Println("Invalid filename:", filename)
		return
	}

	if len(respBody) > 0 {
		err = os.WriteFile(outputPath+filename, respBody, 0644)
		if err != nil {
			fmt.Println("Error writing file", err)
			return
		}
	}
}
