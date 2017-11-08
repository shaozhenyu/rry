package share

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Service struct {
	dir  string
	sep  string
	file *os.File
	tree *Tree
	lock sync.Mutex
}

func (p *Service) GetAll(path string) (leaves Leaves) {
	p.lock.Lock()
	defer p.lock.Unlock()
	tree := p.get(path)
	return tree.GetLeaves(path)
}

func (p *Service) Get(path string) (clocks Clocks, value string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	tree := p.get(path)
	clocks = tree.Clocks().Copy()
	value = tree.Value()
	return
}

func (p *Service) Put(path string, clocks Clocks, value string) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	tree := p.get(path)
	err = tree.Put(clocks, value)
	p.save(path, clocks.ToString(), value)
	return
}

func (p *Service) get(path string) (tree *Tree) {
	if len(path) == 0 {
		return p.tree
	}
	if strings.HasPrefix(path, p.sep) {
		path = path[len(p.sep):]
	}
	fields := strings.Split(path, p.sep)
	tree = p.tree
	for _, name := range fields {
		tree = tree.Get(name)
	}
	return
}

func (p *Service) save(path, clocks, finfo string) (err error) {
	if p.file == nil {
		return
	}

	_, err = p.file.WriteString(fmt.Sprintf("%s\t%s\t%s\n", path, clocks, finfo))
	return
}

func (p *Service) dfile(dir string) string {
	return dir + "/0"
}

func (p *Service) load(dir string) (err error) {
	fmt.Println("load dir:", dir)
	os.MkdirAll(dir, 0777)

	p.file, err = os.OpenFile(p.dfile(dir), os.O_CREATE|os.O_RDWR|os.O_SYNC, 0644)
	if err != nil {
		return
	}

	r := bufio.NewReader(p.file)
	for {
		data, prefix, err := r.ReadLine()
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				return nil
			}
		}

		if prefix {
			return errors.New("loading: line too long")
		}

		line := string(data)
		fields := strings.Split(line, "\t")
		if len(fields) != 3 {
			return errors.New("wrong line: " + line)
		}
		clocks, err := NewClocksFromString(fields[1])
		if err != nil {
			return errors.New("wrong clock: " + fields[1])
		}
		err = p.get(fields[0]).Put(clocks, fields[2])
		if err != nil {
			return errors.New("unexpected error: " + err.Error())
		}
	}

	return
}

func NewService(dir string) (server *Service, err error) {
	server = &Service{sep: "/", tree: NewTree()}
	if len(dir) != 0 {
		err = server.load(dir)
	}
	return
}
