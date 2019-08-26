package note

//go:generate go run ../../gen/main.go

// Interval is the distance between pitches, measured in semitones
type Interval = uint8

const (
	// Semitone is the smallest interval between pitches, a 12th of an octave
	Semitone Interval = 1
	// Tone is 2 semitones
	Tone Interval = 2
	// Octave is the distance between a pitch and another with double frequency
	Octave Interval = 12
)

// Pitch represents a musical frequency
type Pitch uint8

// String returns a human readable name for this pitch
func (p Pitch) String() string {
	return pitchValues[p].name
}

// Frequency returns the hertz value for this pitch
func (p Pitch) Frequency() float64 {
	return pitchValues[p].frequency
}

// Add returns a new Pitch adding an interval to this pitch, making it higher
func (p Pitch) Add(i Interval) Pitch {
	return Pitch(uint8(p) + uint8(i))
}

// Subtract returns a new Pitch subtracting an interval to this pitch, making it lower
func (p Pitch) Subtract(i Interval) Pitch {
	return Pitch(uint8(p) - uint8(i))
}
