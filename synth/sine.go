package synth

import (
	"io"
	"math"
	"time"
)

// NewSineWave returns a SineWave io.Reader. The Reader returns int16 samples
// of a sine wave with maximum amplitude.
// Byte ordering is little endian. The format is:
//     [sample 0 byte 0] [sample 0 byte 1] [sample 1 byte 0] [sample 1 byte 1]...
func NewSineWave(sampleRate int, freq float64, duration time.Duration) io.Reader {
	return &SineWave{
		freq: freq,
		// 2 bytes per sample
		nSamples:   int64(float64(sampleRate) * duration.Seconds()),
		sampleRate: sampleRate,
	}
}

// SineWave is an io.Reader that returns int16 samples. The Reader returns
// int16 samples of a sine wave with maximum amplitude.
// Byte ordering is little endian. The format is:
//     [sample 0 byte 0] [sample 0 byte 1] [sample 1 byte 0] [sample 1 byte 1]...
type SineWave struct {
	freq float64
	// nSamples is the total number of samples that can be read
	nSamples int64
	// offset is measured in number of samples read so far
	offset int64

	sampleRate int
}

func (s *SineWave) Read(buf []byte) (int, error) {
	if s.offset >= s.nSamples {
		return 0, io.EOF
	}

	samplesPeriod := int64(float64(s.sampleRate) / s.freq)

	var i int
	for i = 0; i < len(buf)/2 && s.offset < s.nSamples; i += 2 {
		radian := float64(s.offset) / float64(samplesPeriod) * 2 * math.Pi
		value := equilibrium + int16(float64(max)*math.Sin(radian))

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
