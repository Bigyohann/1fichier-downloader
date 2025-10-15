// Package downloader provides interfaces and structures for downloading files from various providers.
package downloader

import "github.com/cavaliergopher/grab/v3"

// FileInfo represents generic file information from a download provider.
type FileInfo struct {
	Filename    string
	Size        int
	URL         string
	ContentType string
	Checksum    string
	Date        string
}

// DownloadClient defines the contract for a download client.
type DownloadClient interface {
	GetFileData(url string) (FileInfo, error)
	DownloadFile(url string, destinationPath string) (*grab.Response, error)
}
