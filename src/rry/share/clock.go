package share

import (
	"fmt"
	"time"
)

type Clocks map[uint64]Version

type Version struct {
	Ver  uint64
	Time int64
}

func NewClocks() Clocks {
	return make(Clocks)
}

func (p Clocks) Edit(hid uint64) {
	m := uint64(0)
	for _, v := range p {
		if v.Ver > m {
			m = v.Ver
		}
	}
	p[hid] = Version{Ver: m + 1, Time: time.Now().UnixNano() / 1e9}
}

func (p Clocks) ToString() (result string) {
	if len(p) == 0 {
		return
	}

	for k, v := range p {
		result += fmt.Sprintf("%d:%d:%d,", k, v.Ver, v.Time)
	}

	return result[:len(result)-1]
}
