package client

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	LocalConfigFile     = "local.json"
	RemoteConfigFile    = "remote.json"
	DefaultRemoteAdress = "http://127.0.0.1:8080"
)

type Config struct {
	Path   string
	Local  LocalConfig
	Remote RemoteConfig
}

type LocalConfig struct {
	Hid  uint64
	Name string
	Path string
}

type RemoteConfig struct {
	Address string
}

func NewConfig(path, sync string) (config *Config, err error) {
	config = &Config{Path: path}

	lfile, err := os.Open(path + Sep + sync + Sep + LocalConfigFile)
	if err != nil {
		return
	}
	defer lfile.Close()

	err = json.NewDecoder(lfile).Decode(&config.Local)
	if err != nil {
		return
	}

	// may not init remote file
	rfile, err := os.Open(path + Sep + sync + Sep + RemoteConfigFile)
	if err != nil {
		return
	}
	defer rfile.Close()

	err = json.NewDecoder(rfile).Decode(&config.Remote)
	if err != nil {
		return
	}
	if len(config.Remote.Address) == 0 {
		config.Remote.Address = DefaultRemoteAdress
	}

	return
}

func InitLocalConfig(path, sync, user, dpath string) error {
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

	return writeFile(b, path+Sep+sync+Sep+LocalConfigFile)
}

func InitRemoteConfig(path, sync, dpath string) error {
	if strings.HasPrefix(dpath, "/") || strings.HasPrefix(dpath, "\\") {
		return errors.New("path can not start with '/' or '\\': " + dpath)
	}

	os.MkdirAll(path+Sep+sync, 0777)

	remote := RemoteConfig{}
	b, err := json.MarshalIndent(&remote, "", "\t")
	if err != nil {
		return err
	}
	return writeFile(b, path+Sep+sync+Sep+RemoteConfigFile)
}

func writeFile(data []byte, filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	_, err = f.WriteString("\n")
	if err != nil {
		return err
	}

	return nil
}

func (p *Config) paths(path string) (dpath, abs string, err error) {
	if len(path) == 0 {
		dpath = p.Local.Path
		abs = p.Path
	}

	abs, err = filepath.Abs(path)
	if err != nil {
		return
	}

	if !strings.HasPrefix(abs, p.Path) {
		err = errors.New("path not in managed directory(" + p.Path + "): " + abs)
		return
	}

	if p.Path == abs {
		dpath = p.Local.Path
	} else {
		dpath = p.Local.Path + "/" + abs[len(p.Path)+1:len(abs)]
	}
	return
}
