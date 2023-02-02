package x12_writer

func Pad(s string, length int) string {
	var lx = make([]byte, length-len(s))
	for i := range lx {
		lx[i] = 32
	}
	return s + string(lx)
}
