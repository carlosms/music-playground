package synth

import (
	"io"
	"time"
)

const (
	// equilibrium is the middle point of a wave
	equilibrium int16 = 0
	// max is the maximum amplitude of a wave
	max = int16(32767) // (2^16 - 1) / 2
)

// A WaveGenerator returns an io.Reader. The Reader returns int16 samples
// of a wave with maximum amplitude.
// Byte ordering is little endian. The format is:
//     [sample 0 byte 0] [sample 0 byte 1] [sample 1 byte 0] [sample 1 byte 1]...
type WaveGenerator = func(sampleRate int, freq float64, duration time.Duration) io.Reader
