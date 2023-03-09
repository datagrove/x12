package x12

import (
	"fmt"
	"time"
)

type EdiOptions struct {
	Sender     KeyValue
	Receiver   KeyValue
	Gs02       string
	Gs03       string
	Sdelim     string
	Edelim     string
	Cdelim     string
	Rdelim     string
	Production string
}

const (
	CCYYMMDD = "20060102"
	HHMM     = "1504"
)

func NewEdiOptions(send, receiver string) *EdiOptions {
	sendx := KeyValue{
		Key:   send[0:2],
		Value: send[2:],
	}
	receiverx := KeyValue{
		Key:   receiver[0:2],
		Value: receiver[2:],
	}
	r := &EdiOptions{
		Sender:     sendx,
		Receiver:   receiverx,
		Gs03:       receiverx.Value,
		Gs02:       sendx.Value,
		Sdelim:     "~",
		Edelim:     "*",
		Cdelim:     ":",
		Rdelim:     "^",
		Production: "P",
	}

	return r
}
func (op *EdiOptions) NewEdiWriter(path string, controlNumber int) (*EdiWriter, error) {
	tm := time.Now()
	ccyymmdd := tm.Format(CCYYMMDD)
	hhmm := tm.Format(HHMM)
	r := &EdiWriter{
		EdiOptions:    op,
		path:          path,
		ccyymmdd:      ccyymmdd,
		hhmm:          hhmm,
		groupCount:    0,
		controlNumber: fmt.Sprintf("%09d", controlNumber),
	}
	x := op
	v := []string{"ISA", "00", Pad("", 10), "00", Pad("", 10),
		x.Sender.Key, Pad(x.Sender.Value, 15),
		x.Receiver.Key, Pad(x.Receiver.Value, 15),
		ccyymmdd[2:], hhmm, x.Rdelim, "00501", r.controlNumber, "0", "P", x.Cdelim}

	for i, v := range v {
		if i > 0 {
			r.s.WriteString(op.Edelim)
		}
		r.s.WriteString(v)
	}
	r.s.WriteString(op.Sdelim + "\r\n")
	return r, nil
}
