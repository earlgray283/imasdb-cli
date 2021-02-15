package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func DownloadFile(name, dirpath, url string) error {
	fmt.Print("Downloading url: " + url + "...")
	time.Sleep(2 * time.Second)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := os.MkdirAll(dirpath, 0755); err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", dirpath, name)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := io.Copy(file, resp.Body)
	if n == 0 {
		return errors.New("it seems url is 0 byte")
	}

	fmt.Println("done.")

	return err
}
