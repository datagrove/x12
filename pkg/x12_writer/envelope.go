package x12_writer

import (
	"fmt"
	"time"
)

type EdiOptions struct {
	Sender     KeyValue
	Receiver   KeyValue
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

func NewEdiOptions(send, receiver KeyValue) *EdiOptions {
	r := &EdiOptions{
		Sender:     send,
		Receiver:   receiver,
		Sdelim:     "~",
		Edelim:     "*",
		Cdelim:     ":",
		Rdelim:     "^",
		Production: "P",
	}

	return r
}
func (op *EdiOptions) Open(path string, controlNumber int) *EdiWriter {
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
	r.Write("ISA", "00", Pad("", 10), "00", Pad("", 10),
		x.Sender.Key, Pad(x.Sender.Value, 15),
		x.Receiver.Key, Pad(x.Receiver.Value, 15),
		ccyymmdd[2:], hhmm, x.Rdelim, "00501", r.controlNumber, "0", "P", x.Cdelim,
	)
	return r
}
