package x12_writer

import (
	"fmt"
	"os"
	"strings"
)

type EdiWriter struct {
	*EdiOptions
	controlNumber string
	path          string
	ccyymmdd      string
	hhmm          string
	groupCount    int
	stCount       int
	segCount      int
	s             strings.Builder
}

func (w *EdiWriter) Write(r ...string) {
	w.segCount++
	for len(r) > 0 && len(r[len(r)-1]) == 0 {
		r = r[0 : len(r)-1]
	}
	for i, v := range r {
		if i > 0 {
			w.s.WriteString(w.Edelim)
		}
		w.s.WriteString(v)
	}
	w.s.WriteString(w.Sdelim + "\r\n")
}

func (w *EdiWriter) Date(qual, ccyymmdd string) {
	if len(ccyymmdd) > 0 {
		w.Write("DTP", qual, "D8", ccyymmdd)
	}
}
func (w *EdiWriter) Ref(qual, val string) {
	if len(val) > 0 {
		w.Write("REF", qual, val)
	}
}

func (w *EdiWriter) Close() {
	w.Write("IEA", fmt.Sprintf("%d", w.groupCount), w.controlNumber)
	os.WriteFile(w.path, []byte(w.s.String()), 0666)
}

func (w *EdiWriter) BeginGroup() {
	w.groupCount++
	w.stCount = 0
	w.Write("GS", "BE",
		w.EdiOptions.Sender.Value,
		w.EdiOptions.Receiver.Value,
		w.ccyymmdd,
		w.hhmm,
		fmt.Sprintf("%d", w.groupCount),
		"X", "005010X220A1")
}
func (w *EdiWriter) EndGroup() {
	w.Write("GE", fmt.Sprintf("%d", w.stCount), fmt.Sprintf("%d", w.groupCount))
}

func (w *EdiWriter) BeginTransaction(transactionSet string) {
	w.stCount++
	w.segCount = 0
	w.Write("ST",
		transactionSet,
		fmt.Sprintf("%09d", w.stCount),
		"005010X220A1")
}
func (w *EdiWriter) EndTransaction() {
	w.segCount++
	w.Write("SE", fmt.Sprintf("%d", w.segCount), fmt.Sprintf("%09d", w.stCount))
}
