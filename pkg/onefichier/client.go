package onefichier

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/cavaliergopher/grab/v3"
    log "github.com/sirupsen/logrus"
)

type PostDownload struct {
	Url    string `json:"url"`
	Pretty int    `json:"pretty"`
}

type ResponseDownload struct {
	Url    string `json:"url"`
	Status string `json:"status"`
}

type ResponseFileData struct {
	Pass        int    `json:"pass"`
	Description string `json:"description"`
	Acl         int    `json:"acl"`
	Cdn         int    `json:"cdn"`
	Inline      int    `json:"inline"`
	Url         string `json:"url"`
	Filename    string `json:"filename"`
	Size        int    `json:"size"`
	Date        string `json:"date"`
	ContentType string `json:"content-type"`
	Checksum    string `json:"checksum"`
}

func getRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("ONEFICHIER_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func GetDownloadLink(url string) (string, error) {
	// call the API to get the download link

	client := &http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			panic("TODO")
		},
		Jar:     nil,
		Timeout: 0,
	}

	postDownload := PostDownload{
		Url:    url,
		Pretty: 1,
	}
	postDownloadJson, err := json.Marshal(postDownload)
	bodyReader := bytes.NewReader([]byte(postDownloadJson))

	req, err := getRequest(
		"POST",
		"https://api.1fichier.com/v1/download/get_token.cgi",
		bodyReader,
	)

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("Error getting download link")
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Error reading response body")
	}
	responseDownload := ResponseDownload{}
	json.Unmarshal(resBody, &responseDownload)

	if responseDownload.Status != "OK" {
		return "", errors.New("Error getting download link")
	}

	return responseDownload.Url, nil
}

func DownloadFile(url string) (*grab.Response, error) {
	downloadLink, err := GetDownloadLink(url)
	if err != nil {
    log.Error("Error getting download link")
	}

	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(os.Getenv("DOWNLOAD_PATH"), downloadLink)

	// start download
	log.Info("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	return resp, nil
}

func GetFileData(url string) (ResponseFileData, error) {
	client := &http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			log.Warn("Redirecting")
			return nil
		},
		Jar:     nil,
		Timeout: 0,
	}

	postDownload := PostDownload{
		Url:    url,
		Pretty: 1,
	}
	postDownloadJson, err := json.Marshal(postDownload)
	bodyReader := bytes.NewReader([]byte(postDownloadJson))

	req, err := getRequest(
		"POST",
		"https://api.1fichier.com/v1/file/info.cgi",
		bodyReader)

	resp, err := client.Do(req)
	if err != nil {
		return ResponseFileData{}, errors.New("Error getting file data")
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseFileData{}, errors.New("Error reading response body")
	}
	responseFileData := ResponseFileData{}
	json.Unmarshal(resBody, &responseFileData)

	return responseFileData, nil
}
