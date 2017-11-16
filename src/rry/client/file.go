package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var (
	UploadAddress   = "http://127.0.0.1:8081/upload"
	DownloadAddress = "http://127.0.0.1:8081/download/"
)

func Download(lfilename string, rfilename string) error {
	fmt.Printf("Download remote %s -> local %s\n", rfilename, lfilename)
	url := DownloadAddress + rfilename

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Download error :" + rfilename)
	}

	lfile, err := os.OpenFile(lfilename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	_, err = io.Copy(lfile, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func Upload(lfilename string, rfilename string) error {
	fmt.Printf("Upload local %s -> remote %s\n", lfilename, rfilename)
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	fWriter, err := writer.CreateFormFile("uploadfile", rfilename)
	if err != nil {
		return err
	}

	lfile, err := os.Open(lfilename)
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
		return errors.New("http.post file error:" + lfilename)
	}

	return nil
}
