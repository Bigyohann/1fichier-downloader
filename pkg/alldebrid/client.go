// Package alldebrid
package alldebrid

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"bigyohann/apidownloader/internal/service/downloader"

	"github.com/cavaliergopher/grab/v3"
	log "github.com/sirupsen/logrus"
)

// Client implements the downloader.DownloadClient interface for Alldebrid.
type (
	Client       struct{}
	PostDownload struct {
		Link []string `json:"link"`
	}
)

type Info struct {
	Link       string `json:"link"`
	Filename   string `json:"filename"`
	Size       int64  `json:"size"`
	Host       string `json:"host"`
	HostDomain string `json:"hostDomain"`
}

type Response struct {
	Status string `json:"status"`
	Data   struct {
		Infos []Info `json:"infos"`
	} `json:"data"`
}

type ResponseDownload struct {
	Status string `json:"status"`
	Data   *struct {
		Link       string `json:"link"`
		Host       string `json:"host"`
		Filename   string `json:"filename"`
		Paws       bool   `json:"paws"`
		Filesize   int64  `json:"filesize"`
		ID         string `json:"id"`
		Path       string `json:"path"`
		HostDomain string `json:"hostDomain"`
	} `json:"data"`
}

func (c *Client) getDownloadLink(url string) (string, error) {
	// call the API to get the download link

	client := &http.Client{}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	if err := writer.WriteField("link", url); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := getRequest("POST", "https://api.alldebrid.com/v4/link/unlock", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Warnf("Error closing response body: %v", cerr)
		}
	}()

	var responseDownload ResponseDownload

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(respBody, &responseDownload); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}
	fmt.Println(string(respBody))

	if responseDownload.Status != "success" || responseDownload.Data == nil {
		return "", errors.New("failed to get download link from alldebrid")
	}

	return responseDownload.Data.Link, nil
}

func getRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("ALLDEBRID_TOKEN"))
	return req, nil
}

func (c *Client) DownloadFile(url string, destinationPath string) (*grab.Response, error) {
	downloadLink, err := c.getDownloadLink(url)
	fmt.Println("Download link:", downloadLink)
	if err != nil {
		return nil, fmt.Errorf("error getting download link: %w", err)
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
	var fileInfo downloader.FileInfo

	client := &http.Client{}

	// Prepare multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	if err := writer.WriteField("link[]", url); err != nil {
		return fileInfo, err
	}
	if err := writer.Close(); err != nil {
		return fileInfo, err
	}

	req, err := getRequest("POST", "https://api.alldebrid.com/v4/link/infos", &buf)
	if err != nil {
		return fileInfo, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return fileInfo, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Warnf("Error closing response body: %v", cerr)
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return fileInfo, err
	}

	var response Response

	if err := json.Unmarshal(respBody, &response); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return fileInfo, err
	}

	if response.Status != "success" || response.Data.Infos[0].Size == 0 {
		return fileInfo, errors.New("failed to get file data from alldebrid")
	}

	fileInfo = downloader.FileInfo{
		Filename:    response.Data.Infos[0].Filename,
		Size:        int(response.Data.Infos[0].Size),
		URL:         response.Data.Infos[0].Link,
		ContentType: "",
		Checksum:    "",
		Date:        "",
	}

	return fileInfo, nil
}
