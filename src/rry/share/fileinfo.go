package share

import (
	"encoding/json"
)

type FileInfo struct {
	Deleted bool
	Sha1    string
}

func (p *FileInfo) Equals(x string) bool {
	value, err := p.ToString()
	if err != nil {
		return false
	}
	return value == x
}

func NewFileInfo(data string) (p *FileInfo, err error) {
	p = &FileInfo{}
	err = json.Unmarshal([]byte(data), p)
	return
}

func (p *FileInfo) ToString() (value string, err error) {
	b, err := json.Marshal(p)
	if err != nil {
		return
	}

	value = string(b)
	return
}
