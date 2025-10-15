package downloader

import (
	"encoding/json"
	"os"
	"time"

	"bigyohann/apidownloader/internal/database"
	"bigyohann/apidownloader/internal/database/models"
	"bigyohann/apidownloader/internal/service"

	log "github.com/sirupsen/logrus"
)

func downloadFile(client DownloadClient, url string, file models.File) {
	// download file
	resp, err := client.DownloadFile(url, os.Getenv("DOWNLOAD_PATH"))
	file.Status = "Downloading"
	db := database.GetDB()
	db.Save(&file)
	if err != nil {
		log.Fatal(err)
		return
	}

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
Loop:
	for {
		select {
		case <-t.C:
			// calculate download speed
			speed := resp.BytesPerSecond() / (1024 * 1024)
			if speed > 0 {
				log.Printf("  %vMB/s\n", speed)
			}
			log.Printf("  transferred %v/%vMB (%.2f%%)\n",
				int(resp.BytesComplete())/(1024*1024),
				int(resp.Size())/(1024*1024),
				100*resp.Progress())
			// object to json string
			jsonStr, _ := json.Marshal(map[string]any{
				"fileId":   file.ID,
				"speed":    int(speed),
				"progress": int(100 * resp.Progress()),
				"status":   "Downloading",
			})
			service.PushToMercureHub("file_update", string(jsonStr))

		case <-resp.Done:
			// download is complete
			file.Downloaded = true
			file.Status = "Downloaded"
			db := database.GetDB()
			db.Save(&file)
			jsonStr, _ := json.Marshal(map[string]any{
				"fileId": file.ID,
				"status": "Downloaded",
			})
			service.PushToMercureHub("file_update", string(jsonStr))
			break Loop
		}
	}
	// handle download Error
	if err := resp.Err(); err != nil {
		file.Status = "Error"
		db := database.GetDB()
		db.Save(&file)
	}
}

func HandleDownloadFile(client DownloadClient, url string) models.File {
	fileData, err := client.GetFileData(url)
	if err != nil {
		log.Error(err)
	}

	file := models.File{}
	db := database.GetDB()
	db.Where("filename = ?", fileData.Filename).First(&file)

	if file.ID != 0 && file.Status != "Error" {
		log.Info("File already downloaded / in download")
		return file
	}

	file = models.File{
		Filename:    fileData.Filename,
		Size:        fileData.Size,
		URL:         fileData.URL,
		Downloaded:  false,
		ContentType: fileData.ContentType,
		Checksum:    fileData.Checksum,
		Date:        fileData.Date,
		Status:      "Created",
	}
	db.Create(&file)

	go downloadFile(client, url, file)

	return file
}
