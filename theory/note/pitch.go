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

// Pitch is an interface for a note musical frequency
type Pitch interface {
	// String returns a human readable name for this pitch
	String() string
	// Frequency returns the hertz value for this pitch
	Frequency() float64
	// Add returns a new Pitch adding an interval to this pitch, making it higher
	Add(i Interval) Pitch
	// Subtract returns a new Pitch subtracting an interval to this pitch, making it lower
	Subtract(i Interval) Pitch
}

// restPitch represents the (absence of) musical frequency for a rest note
type restPitch struct{}

// Frequency returns the 0 hertz value
func (r restPitch) Frequency() float64 {
	return 0
}

// Add is a no operation, returns a RestPitch
func (r restPitch) Add(i Interval) Pitch {
	return r
}

// Subtract is a no operation, returns a RestPitch
func (r restPitch) Subtract(i Interval) Pitch {
	return r
}

// String is a no operation, returns an empty string
func (r restPitch) String() string {
	return ""
}

// pitchValue represents a musical frequency
type pitchValue uint8

// String returns a human readable name for this pitch
func (p pitchValue) String() string {
	return pitchValues[p].name
}

// Frequency returns the hertz value for this pitch
func (p pitchValue) Frequency() float64 {
	return pitchValues[p].frequency
}

// Add returns a new Pitch adding an interval to this pitch, making it higher
func (p pitchValue) Add(i Interval) Pitch {
	return pitchValue(uint8(p) + uint8(i))
}

// Subtract returns a new Pitch subtracting an interval to this pitch, making it lower
func (p pitchValue) Subtract(i Interval) Pitch {
	return pitchValue(uint8(p) - uint8(i))
}
