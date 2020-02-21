package convert

func Uint32SliceToByte(in []uint32) []byte {
	out := make([]byte, len(in)*4)
	for k, v := range in {
		out[k*4] = byte(v)
		out[k*4+1] = byte(v >> 8)
		out[k*4+2] = byte(v >> 16)
		out[k*4+3] = byte(v >> 24)
	}
	return out
}
