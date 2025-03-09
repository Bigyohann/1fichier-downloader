package downloader

import (
	"bigyohann/apidownloader/internal/database"
	"bigyohann/apidownloader/internal/database/models"
	"bigyohann/apidownloader/pkg/onefichier"
	"fmt"
	"log"
	"time"
)

func downloadFile(url string, file models.File) {
	// download file
	resp, err := onefichier.DownloadFile(url)
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
			fmt.Printf("  transferred %v/%vMB (%.2f%%)\n",
				int(resp.BytesComplete())/(1024*1024),
				int(resp.Size())/(1024*1024),
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			file.Downloaded = true
			db := database.GetDB()
			db.Save(&file)
			break Loop
		}
	}
}

func HandleDownloadFile(url string) models.File {
	fileData, err := onefichier.GetFileData(url)
	if err != nil {
		fmt.Println(err)
	}

	file := models.File{}
	db := database.GetDB()
	db.Where("filename = ?", fileData.Filename).First(&file)

	if file.ID != 0 {
		fmt.Println("File already downloaded / in download")
		return file
	}

	file = models.File{
		Filename:    fileData.Filename,
		Size:        fileData.Size,
		Url:         fileData.Url,
		Downloaded:  false,
		ContentType: fileData.ContentType,
		Checksum:    fileData.Checksum,
		Date:        fileData.Date,
	}
	db.Create(&file)

	go downloadFile(url, file)

	return file
}
