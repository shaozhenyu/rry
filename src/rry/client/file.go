package client

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var (
	UploadAddress   = "http://127.0.0.1:8081/upload"
	DownloadAddress = "http://127.0.0.1:8081/download/"
)

func Download(filename string) error {
	url := DownloadAddress + filename

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	lfile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	_, err = io.Copy(lfile, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func Upload(filename string) error {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	fWriter, err := writer.CreateFormFile("uploadfile", filename)
	if err != nil {
		return err
	}

	lfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer lfile.Close()

	_, err = io.Copy(fWriter, lfile)
	if err != nil {
		return err
	}

	contentType := writer.FormDataContentType()
	err = writer.Close()
	if err != nil {
		return err
	}

	resp, err := http.Post(UploadAddress, contentType, buf)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("http.post file error:" + filename)
	}

	return nil
}
