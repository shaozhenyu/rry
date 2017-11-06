package share

import (
	"fmt"
	"os"
	"sync"
)

type Service struct {
	dir  string
	sep  string
	file *os.File
	lock sync.Mutex
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

	return
}

func (p *Service) Save(path, clocks, finfo string) (err error) {
	if p.file == nil {
		return
	}

	_, err = p.file.WriteString(fmt.Sprintf("%s\t%s\t%s\n", path, clocks, finfo))
	return
}

func NewService(dir string) (server *Service, err error) {
	server = &Service{sep: "/"}
	if len(dir) != 0 {
		err = server.load(dir)
	}
	return
}
