package share

import (
	"encoding/json"
)

type FileInfo struct {
	Deleted bool
	Sha1    string
}

func (p *FileInfo) ToString() (value string, err error) {
	b, err := json.Marshal(p)
	if err != nil {
		return
	}

	value = string(b)
	return
}
