// copy file
package main

import (
	"io"
	"log"
	"os"
)

func main() {
	src, err := os.Open("test")
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	dst, err := os.OpenFile("test1", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()

	io.Copy(dst, src)
}
