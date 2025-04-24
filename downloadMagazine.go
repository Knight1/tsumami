package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type PdfResponse struct {
	DownloadURL string `json:"downloadUrl"`
	WaitSec     *int   `json:"waitSec"`
}

// heiseplus = ix, ct, make, tr, mac-and-i, ct-foto,
// ct-wissen, ix-special, heise-online-sonderhefte
func downloadMagazine(sessionCookie *http.Cookie, magazine string, year string, issue string) {
	downloadThumbnail(magazine, year, issue)
	outputPath := "magazines/" + magazine + "/" + year + "/"

	var res PdfResponse
	var resp *http.Response
	var err error
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest("GET", "https://www.heise.de/select/"+magazine+"/archiv/"+year+"/"+issue+"/download", nil)
		req.AddCookie(sessionCookie)

		client := &http.Client{
			Timeout: time.Second * 30,
		}

		req.Header.Add("Accept", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			fmt.Println("Error getting Magazine:", err)
			time.Sleep(time.Second * 30)
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * 30)
		}

		if resp.StatusCode != 200 {
			fmt.Printf("Status: %s\n", resp.Status)
			fmt.Printf("Headers:\n%s", resp.Header)
			fmt.Printf("Body:\n%s", string(respBody))
			return
		}

		err = json.Unmarshal(respBody, &res)
		if err != nil {
			panic(err)
		}

		if res.WaitSec != nil {
			fmt.Printf("WaitSec: %d\n", *res.WaitSec)
			time.Sleep(time.Duration(*res.WaitSec) * time.Second)
		}

		if res.DownloadURL != "" && res.WaitSec == nil {
			break
		}
	}

	// Rewrite the URL to include IPv6 Support
	const newPrefix = "https://s3.dualstack.eu-west-1.amazonaws.com/pdf-abo/"
	const oldPrefix = "https://pdf-abo.s3.amazonaws.com/"
	if strings.HasPrefix(res.DownloadURL, oldPrefix) {
		res.DownloadURL = newPrefix + strings.TrimPrefix(res.DownloadURL, oldPrefix)
	}

	req, err := http.NewRequest("GET", res.DownloadURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if verbose {
		fmt.Println(res.DownloadURL)
	}

	client := &http.Client{
		Timeout: time.Second * 120,
	}

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error downloading %s: %s\n", res.DownloadURL, err)
		return
	}

	if !isValidPDF(respBody) {
		fmt.Printf("Download failed: %s\n", resp.Status)
		return
	}

	filename := path.Base(res.DownloadURL)

	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		fmt.Println(err)
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

func isValidPDF(data []byte) bool {
	return bytes.HasPrefix(data, []byte("%PDF-")) &&
		bytes.Contains(data, []byte("%%EOF"))
}
