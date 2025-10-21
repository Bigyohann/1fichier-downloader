// Package onefichier implements functions to interact with the 1fichier API.
package onefichier

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"bigyohann/apidownloader/internal/service/downloader"

	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
)

// Client implements the downloader.DownloadClient interface for 1fichier.
type Client struct{}

type PostDownload struct {
	URL    string `json:"url"`
	Pretty int    `json:"pretty"`
}

type ResponseDownload struct {
	URL    string `json:"url"`
	Status string `json:"status"`
}

type ResponseFileData struct {
	Pass        int    `json:"pass"`
	Description string `json:"description"`
	ACL         int    `json:"acl"`
	Cdn         int    `json:"cdn"`
	Inline      int    `json:"inline"`

	URL         string `json:"url"`
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

func (c *Client) getDownloadLink(url string) (string, error) {
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
		URL:    url,
		Pretty: 1,
	}
	postDownloadJSON, _ := json.Marshal(postDownload)
	bodyReader := bytes.NewReader([]byte(postDownloadJSON))

	req, _ := getRequest(
		"POST",
		"https://api.1fichier.com/v1/download/get_token.cgi",
		bodyReader,
	)

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("error getting download link")
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("error reading response body")
	}
	responseDownload := ResponseDownload{}
	json.Unmarshal(resBody, &responseDownload)

	if responseDownload.Status != "OK" {
		return "", errors.New("error getting download link")
	}

	return responseDownload.URL, nil
}

func (c *Client) DownloadFile(url string, destinationPath string) (*grab.Response, error) {
	downloadLink, err := c.getDownloadLink(url)
	if err != nil {
		log.Error("Error getting download link")
	}

	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(destinationPath, downloadLink)

	// start download
	log.Infof("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	return resp, nil
}

func (c *Client) GetFileData(url string) (downloader.FileInfo, error) {
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
		URL:    url,
		Pretty: 1,
	}
	postDownloadJSON, _ := json.Marshal(postDownload)
	bodyReader := bytes.NewReader([]byte(postDownloadJSON))

	req, _ := getRequest(
		"POST",
		"https://api.1fichier.com/v1/file/info.cgi",
		bodyReader)

	resp, err := client.Do(req)
	if err != nil {
		return downloader.FileInfo{}, errors.New("error getting file data")
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return downloader.FileInfo{}, errors.New("error reading response body")
	}
	responseFileData := ResponseFileData{}
	json.Unmarshal(resBody, &responseFileData)

	fileInfo := downloader.FileInfo{
		Filename:    responseFileData.Filename,
		Size:        responseFileData.Size,
		URL:         responseFileData.URL,
		ContentType: responseFileData.ContentType,
		Checksum:    responseFileData.Checksum,
		Date:        responseFileData.Date,
	}

	return fileInfo, nil
}
