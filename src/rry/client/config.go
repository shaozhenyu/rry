package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	LocalConfigFile = "local.json"
)

type LocalConfig struct {
	Hid  uint64
	Name string
	Path string
}

func InitLocalConfig(path, sync, user, dpath string) error {
	fmt.Println(path, sync, user, dpath)
	if strings.HasPrefix(dpath, "/") || strings.HasPrefix(dpath, "\\") {
		return errors.New("path can not start with '/' or '\\': " + dpath)
	}

	os.MkdirAll(path+Sep+sync, 0777)

	rand.Seed(time.Now().UnixNano())
	hid := uint64(rand.Uint32())

	local := LocalConfig{
		Hid:  hid,
		Name: user,
		Path: dpath,
	}

	b, err := json.MarshalIndent(&local, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	f, err := os.OpenFile(path+Sep+sync+Sep+LocalConfigFile, os.O_CREATE|os.O_EXCL|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	_, err = f.WriteString("\n")
	if err != nil {
		return err
	}

	return nil
}
