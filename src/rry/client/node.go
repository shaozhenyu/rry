package client

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"rry/share"
)

const MetaPath = "meta"

type Node struct {
	config  *Config
	service *share.Service
}

func NewNode(path, sync string) (p *Node, err error) {
	config, err := NewConfig(path, sync)
	if err != nil {
		return
	}

	meta := config.Path + Sep + sync + Sep + MetaPath
	err = os.MkdirAll(meta, 0777)
	if err != nil {
		return
	}

	service, err := share.NewService(meta)
	if err != nil {
		return
	}

	p = &Node{
		config:  config,
		service: service,
	}
	return
}

func (p *Node) Sync(path string) (err error) {
	fmt.Println("node sync:", path)
	dpath, abs, err := p.config.paths(path)
	if err != nil {
		return
	}
	fmt.Println("sync dpath:", dpath, "abs:", abs)

	done := map[string]bool{}

	err = filepath.Walk(abs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.New("walking dir: " + err.Error())
		}

		if info.IsDir() {
			return nil
		}

		dpath, _, err = p.config.paths(path)
		if err != nil {
			return err
		}

		related := dpath[len(p.config.Local.Path)+1 : len(dpath)]

		if strings.HasPrefix(related, SzyDir) {
			return nil
		}

		if !strings.HasPrefix(related, "data") {
			return nil
		}

		if done[dpath] {
			return nil
		}

		err = p.sync(dpath)
		if err != nil {
			return err
		}
		done[dpath] = true

		return nil
	})

	return
}

func (p *Node) sync(dpath string) (err error) {
	fmt.Println("file need sync:", dpath)
	abs := p.config.Path + dpath[len(p.config.Local.Path):len(dpath)]
	fmt.Println("file abs :", abs)
	related := dpath[len(p.config.Local.Path)+1 : len(dpath)]
	fmt.Println("file related:", related)

	finfo, err := p.info(abs)
	if err != nil {
		return
	}

	clocks := share.NewClocks()
	clocks.Edit(p.config.Local.Hid)
	value, err := finfo.ToString()
	if err != nil {
		return errors.New("finfo.ToString: " + err.Error())
	}
	p.service.Save(dpath, clocks.ToString(), value)

	return
}

func (p *Node) info(abs string) (info share.FileInfo, err error) {
	finfo, err := os.Stat(abs)
	if os.IsNotExist(err) {
		info = share.FileInfo{true, ""}
		return
	} else if err != nil {
		return
	} else if finfo.IsDir() {
		err = errors.New("ignore directory: " + abs)
		return
	}

	f, err := os.Open(abs)
	if err != nil {
		return
	}

	sha := sha1.New()
	_, err = io.Copy(sha, f)
	if err != nil {
		return
	}

	info = share.FileInfo{true, base64.StdEncoding.EncodeToString(sha.Sum(nil))}
	return

}
