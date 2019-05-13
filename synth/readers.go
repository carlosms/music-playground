package synth

import (
	"io"
	"math"
)

// Combine takes an arbitrary number of Readers that return int16 samples, and
// returns a Reader that combines those samples into one single sample value
func Combine(readers ...io.Reader) io.Reader {
	return &CombinedReader{readers}
}

// CombinedReader takes an arbitrary number of Readers that return int16 samples,
// and combines those samples into one single sample value
type CombinedReader struct {
	readers []io.Reader // underlying readers
}

func (m *CombinedReader) Read(p []byte) (int, error) {
	var eofErr error
	var n int
	for n = 0; n < len(p)-1; n += 2 {
		var total int16
		eofErr = io.EOF

		for _, r := range m.readers {
			// read 1 sample (2 bytes)
			buf := make([]byte, 2)
			n, err := io.ReadFull(r, buf)

			if n == 2 {
				eofErr = nil

				// Convert 2 bytes to int16, little-endian
				v := int16(buf[0]) + int16(buf[1])<<8
				total = add(total, v)
			}

			// ErrUnexpectedEOF means the reader had less than 2 bytes, we can't
			// use that as a sample so it is also discarded gracefully
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return n - 2, err
			}
		}

		// int16 back to to 2 bytes, little-endian
		p[n] = byte(total)
		p[n+1] = byte(total >> 8)
	}

	return n, eofErr
}

func add(a, b int16) int16 {
	v := a + b

	signA := a < 0
	signB := b < 0
	signV := v < 0

	// negative + negative = positive, overflow
	if signA && signB && !signV {
		return math.MinInt16
	}

	// positive + positive = negative, overflow
	if !signA && !signB && signV {
		return math.MaxInt16
	}

	return v
}
