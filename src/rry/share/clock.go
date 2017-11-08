package share

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	Equal      = 0
	Smaller    = -1
	Greater    = 1
	Conflicted = -2
)

type Clocks map[uint64]Version

type Version struct {
	Ver  uint64
	Time int64
}

func NewClocks() Clocks {
	return make(Clocks)
}

func NewClocksFromString(str string) (clocks Clocks, err error) {
	clocks = make(Clocks)
	if len(str) == 0 {
		return
	}

	segs := strings.Split(str, ",")
	for _, seg := range segs {
		fields := strings.Split(seg, ":")
		if len(fields) != 3 {
			err = errors.New("bad clocks format:" + str)
		}
		var key, ver, tim int
		key, err = strconv.Atoi(fields[0])
		if err != nil {
			return
		}
		ver, err = strconv.Atoi(fields[1])
		if err != nil {
			return
		}
		tim, err = strconv.Atoi(fields[2])
		if err != nil {
			return
		}
		clocks[uint64(key)] = Version{
			Ver:  uint64(ver),
			Time: int64(tim),
		}
	}
	return
}

func (p Clocks) Compare(x Clocks) int {
	l1 := len(p)
	l2 := len(x)

	if l1 == l2 {
		c := Equal
		for k, v1 := range p {
			v2, ok := x[k]
			if !ok {
				return Conflicted
			}

			if v1.Ver == v2.Ver {
				continue
			}

			if v1.Ver < v2.Ver {
				if c == Greater {
					return Conflicted
				} else {
					c = Smaller
				}
			}

			if v1.Ver > v2.Ver {
				if c == Smaller {
					return Conflicted
				} else {
					c = Greater
				}
			}
		}
		return c
	}

	if l1 < l2 {
		for k, v1 := range p {
			v2, ok := x[k]
			if !ok {
				return Conflicted
			}
			if v1.Ver == v2.Ver {
				continue
			}
			if v1.Ver > v2.Ver {
				return Conflicted
			}
		}
		return Smaller
	}

	if l1 > l2 {
		for k, v2 := range x {
			v1, ok := p[k]
			if !ok {
				return Conflicted
			}
			if v1.Ver == v2.Ver {
				continue
			}
			if v1.Ver < v2.Ver {
				return Conflicted
			}
		}
		return Greater
	}
	panic("unexpected")
}

func (p Clocks) Sig(sep string) string {
	mk, mv := p.max()
	return fmt.Sprintf("%v"+sep+"%v"+sep+"%x", len(p), mv, mk)
}

func (p Clocks) max() (uint64, uint64) {
	var mk uint64
	var mv uint64
	for k, v := range p {
		if v.Ver > mv {
			mk = k
			mv = v.Ver
		}
	}
	return mk, mv
}

func (p Clocks) Absorb(x Clocks) {
	for k, v2 := range x {
		v1, ok := p[k]
		if !ok || v2.Ver > v1.Ver {
			p[k] = v2
		}
	}
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

func (c Clocks) Copy() Clocks {
	m := make(Clocks)
	for k, v := range c {
		m[k] = v
	}
	return m
}

func CompareResultString(result int) string {
	switch result {
	case Equal:
		return "equal"
	case Smaller:
		return "smaller"
	case Greater:
		return "greater"
	case Conflicted:
		return "conflicted"
	}
	return "unknown"
}
