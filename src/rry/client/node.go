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

	"rry/server"
	"rry/share"

	"ember/http/rpc"
)

const MetaPath = "meta"

type Node struct {
	config  *Config
	service *share.Service
	client  *server.Client
	cache   *share.Service
	rpc     *rpc.Client
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

	// cache local file info
	service, err := share.NewService(meta)
	if err != nil {
		return
	}

	client := &server.Client{}
	rpc := rpc.NewClient(config.Remote.Address)
	err = rpc.Reg(client)
	if err != nil {
		return
	}

	// cache remote file info
	cache, err := share.NewService("")
	if err != nil {
		return
	}

	p = &Node{
		config:  config,
		service: service,
		client:  client,
		cache:   cache,
		rpc:     rpc,
	}
	return
}

func (p *Node) Sync(path string) (err error) {
	dpath, abs, err := p.config.paths(path)
	if err != nil {
		return
	}

	leaves, err := p.client.Gets(dpath)
	if err != nil {
		return
	}
	done := map[string]bool{}

	for k, v := range leaves {
		related := k[len(p.config.Local.Path)+1:]

		// TODO just sync path data now
		if !strings.HasPrefix(related, "data") {
			done[k] = true
			return nil
		}

		var clocks share.Clocks
		clocks, err = share.NewClocksFromString(v.Clocks)
		if err != nil {
			return errors.New("cache.Generating: " + err.Error())
		}
		err = p.cache.Put(k, clocks, v.Value)
		if err != nil {
			return errors.New("cache.Put: " + err.Error())
		}

		if done[k] {
			continue
		}

		err = p.sync(k)
		if err != nil {
			return err
		}
		done[k] = true
	}

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

		// filter config path
		if strings.HasPrefix(related, ConfigDir) {
			return nil
		}

		// TODO just sync path data now
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
	abs := p.config.Path + dpath[len(p.config.Local.Path):len(dpath)]
	//related := dpath[len(p.config.Local.Path)+1 : len(dpath)]

	rclocks, rvalue := p.cache.Get(dpath)
	var rinfo *share.FileInfo
	if len(rvalue) != 0 {
		rinfo, err = share.NewFileInfo(rvalue)
		if err != nil {
			return errors.New(rvalue + " parse rvalue: " + err.Error())
		}
	}

	lclocks, lvalue := p.service.Get(dpath)
	var linfo *share.FileInfo
	if len(lvalue) != 0 {
		linfo, err = share.NewFileInfo(lvalue)
		if err != nil {
			return errors.New(lvalue + " parse lvalue: " + err.Error())
		}
	}

	status := rclocks.Compare(lclocks)
	ss := "local: " + lclocks.Sig(":") + ", remote: " + rclocks.Sig(":")
	if status == share.Smaller {
		return errors.New(dpath + " local version greater than remote version, unknown error. " + ss)
	}
	if status == share.Conflicted {
		return errors.New(dpath + " conflicted from multiple source. " + ss)
	}

	finfo, err := p.info(abs)
	if err != nil {
		return errors.New("node.GetInfo: " + err.Error())
	}

	if status == share.Greater {
		fmt.Println("remote > local:", ss)
		if !finfo.Equals(lvalue) && !finfo.Deleted {
			if finfo.Equals(rvalue) {
				fmt.Println("file no changed, updating db(local)")
				err = p.service.Put(dpath, rclocks, rvalue)
				if err != nil {
					return errors.New("local.Put: " + err.Error())
				}
			} else {
				fmt.Println("local also edited, confilicted")
				return errors.New(dpath + " conflicted. " + ss)
			}
		} else {
			if rinfo.Deleted {
				fmt.Println("remote removed, removing local file")
				err = os.Remove(abs)
				if err != nil {
					if os.IsNotExist(err) {
						fmt.Println("local file not exists, skip")
						err = nil
					} else {
						return errors.New("os.Remove: " + err.Error())
					}
				}
			} else {
				tmp := abs + ".downloading"
				err = os.MkdirAll(filepath.Dir(tmp), 0777)
				if err != nil {
					return errors.New("os.Mkdir: " + err.Error())
				}
				err = Download(tmp, dpath)
				if err != nil {
					return errors.New("download: " + err.Error())
				}
				err = os.Rename(tmp, abs)
				if err != nil {
					return errors.New("os.Rename: " + err.Error())
				}
			}

			err = p.service.Put(dpath, rclocks, rvalue)
			if err != nil {
				return errors.New("local.Put: " + err.Error())
			}
		}
		return
	}

	if status != share.Equal {
		panic("gg")
	}

	if finfo.Deleted && (linfo == nil || linfo.Deleted) || finfo.Equals(lvalue) {
		//fmt.Println(abs, " nothing happen")
		return
	}

	if !finfo.Deleted {
		err = Upload(abs, dpath)
		if err != nil {
			return errors.New("upload: " + err.Error())
		}
	}

	clocks := rclocks.Copy()
	clocks.Absorb(lclocks)
	clocks.Edit(p.config.Local.Hid)
	value, err := finfo.ToString()
	if err != nil {
		return errors.New("finfo.ToString: " + err.Error())
	}

	err = p.client.Put(dpath, clocks.ToString(), value)
	if err != nil {
		return errors.New("remote.Put: " + err.Error())
	}

	err = p.cache.Put(dpath, clocks, value)
	if err != nil {
		return errors.New("cache.Put: " + err.Error())
	}

	err = p.service.Put(dpath, clocks, value)
	if err != nil {
		return errors.New("local.Put: " + err.Error())
	}

	return
}

func (p *Node) info(abs string) (info share.FileInfo, err error) {
	finfo, err := os.Stat(abs)
	if os.IsNotExist(err) {
		err = nil
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

	info = share.FileInfo{false, base64.StdEncoding.EncodeToString(sha.Sum(nil))}
	return

}
