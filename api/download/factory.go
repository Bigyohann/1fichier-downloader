// Package download provides a factory for creating download clients.
package download

import (
	"errors"

	"bigyohann/apidownloader/internal/service/downloader"
	"bigyohann/apidownloader/pkg/alldebrid"
	"bigyohann/apidownloader/pkg/onefichier"
)

// GetClient returns a download client for the given provider.
func GetClient(provider string) (downloader.DownloadClient, error) {
	switch provider {
	case "1fichier":
		return &onefichier.Client{}, nil
	case "alldebrid":
		return &alldebrid.Client{}, nil
	default:
		return nil, errors.New("invalid provider")
	}
}
