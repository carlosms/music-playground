package synth

import (
	"io"
)

// Sustain takes a Reader that returns int16 samples, and returns a Reader
// that multiplies each value by the given percentage. The percentage must
// be a value between 0 and 1
func Sustain(r io.Reader, percentage float64) io.Reader {
	if percentage < 0 || percentage > 1 {
		panic("percentage must be between 0 and 1")
	}
	return &SustainedReader{r, percentage}
}

// SustainedReader takes a Reader that returns int16 samples, and multiplies
// each value by the given percentage. The percentage must be a value between
// 0 and 1
type SustainedReader struct {
	r          io.Reader // underlying reader
	percentage float64
}

func (s *SustainedReader) Read(p []byte) (n int, err error) {
	n, err = s.r.Read(p)

	for i := 0; i < n-1; i += 2 {
		// Convert 2 bytes to int16, little-endian
		v := int16(p[i]) + int16(p[i+1])<<8

		v = int16(float64(v) * s.percentage)

		// int16 back to to 2 bytes, little-endian
		p[i] = byte(v)
		p[i+1] = byte(v >> 8)
	}

	return
}
