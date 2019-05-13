package synth

import (
	"io"
	"time"
)

// NewSquareWave returns a SquareWave io.Reader. The Reader returns int16 samples
// of a square wave with maximum amplitude.
// Byte ordering is little endian. The format is:
//     [sample 0 byte 0] [sample 0 byte 1] [sample 1 byte 0] [sample 1 byte 1]...
func NewSquareWave(sampleRate int, freq float64, duration time.Duration) io.Reader {
	return &SquareWave{
		freq: freq,
		// 2 bytes per sample
		nSamples:   int64(float64(sampleRate) * duration.Seconds()),
		sampleRate: sampleRate,
	}
}

// SquareWave is an io.Reader that returns int16 samples. The Reader returns
// int16 samples of a square wave with maximum amplitude.
// Byte ordering is little endian. The format is:
//     [sample 0 byte 0] [sample 0 byte 1] [sample 1 byte 0] [sample 1 byte 1]...
type SquareWave struct {
	freq float64
	// nSamples is the total number of samples that can be read
	nSamples int64
	// offset is measured in number of samples read so far
	offset int64

	sampleRate int
}

func (s *SquareWave) Read(buf []byte) (int, error) {
	if s.offset >= s.nSamples {
		return 0, io.EOF
	}

	samplesPeriod := int64(float64(s.sampleRate) / s.freq)
	samplesHalfPeriod := samplesPeriod / 2

	var i int
	for i = 0; i < len(buf)/2 && s.offset < s.nSamples; i += 2 {
		pos := s.offset % samplesPeriod

		var value int16
		if pos <= samplesHalfPeriod {
			value = equilibrium + max
		} else {
			value = equilibrium - max
		}

		// int16 to 2 bytes, little-endian
		buf[i] = byte(value)
		buf[i+1] = byte(value >> 8)

		s.offset++
	}

	if s.offset >= s.nSamples {
		return i, io.EOF
	}

	return i, nil
}
