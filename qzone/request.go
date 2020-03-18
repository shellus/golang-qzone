package qzone

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// DownloadFile will downloadPhoto a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string, t time.Time) error {
	fmt.Printf("downloadPhoto file into: %s \n from url: %s\n", filepath, url)
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = os.Chtimes(filepath, t, t)
	if err != nil {
		fmt.Printf("chtimes err: %s\n", err.Error())
	}

	return err
}
