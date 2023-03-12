package x12

import (
	"os"
	"strings"
)

type Segment struct {
	Line    int
	Segid   string
	Element [][]string
}

type Isa struct {
	EdiOptions
	data string
}

// ISA*00*          *00*          *ZZ*131628401      *ZZ*AETNA_POS      *230217*0953*^*00501*000003892*0*P*>~
func NewIsa(d string) *Isa {
	r := &Isa{data: d}
	r.Sender.Key = d[32:34]
	r.Sender.Value = d[35:50]
	r.Receiver.Key = d[51:53]
	r.Receiver.Value = d[54:69]
	r.Sdelim = d[105:106]
	r.Cdelim = d[104:105]
	r.Edelim = d[3:4]
	r.Rdelim = d[82:83]
	return r
}

func (s *Segment) El(n int) string {
	return s.Element[n][0]
}

func (s *Segment) Set(n int, value string) {
	s.Element[n] = []string{value}
}

func (s *Segment) CopyTo(w *EdiStream) {
	for i, v := range s.Element {
		if i > 0 {
			w.s.WriteString(w.Edelim)
		}
		for j, o := range v {
			if j > 0 {
				w.s.WriteString(w.Cdelim)
			}
			w.s.WriteString(o)
		}
	}
	w.s.WriteString(w.Sdelim + "\r\n")
}
func ReadEdi(f string, fn func(s Segment) error) error {
	b, e := os.ReadFile(f)
	if e != nil {
		return e
	}
	s := string(b)
	isa := NewIsa(s)
	lns := strings.Split(s[106:], isa.Sdelim)
	var seg Segment
	for i, ln := range lns {
		seg.Line = i
		el := strings.Split(ln, isa.Edelim)
		seg.Element = make([][]string, len(el))
		for i := range el {
			seg.Element[i] = strings.Split(el[i], isa.Cdelim)
		}
		seg.Segid = seg.Element[0][0]
		fn(seg)
	}
	return nil
}
