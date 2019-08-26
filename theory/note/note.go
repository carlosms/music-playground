package note

import (
	"fmt"
)

// Note is a musical Pitch and Duration
type Note struct {
	Pitch
	Duration
}

// NewNote creates a new musical Note
func NewNote(p Pitch, d Duration) Note {
	return Note{
		Pitch:    p,
		Duration: d,
	}
}

// String returns a human readable representation of this note
func (n Note) String() string {
	return fmt.Sprintf("%v %v", n.Duration, n.Pitch)
}

// Add returns a new Note adding an interval to this note's pitch, making it higher
func (n Note) Add(i Interval) Note {
	return Note{
		Pitch:    n.Pitch.Add(i),
		Duration: n.Duration,
	}
}

// Subtract returns a new Note subtracting an interval to this note's pitch, making it lower
func (n Note) Subtract(i Interval) Note {
	return Note{
		Pitch:    n.Pitch.Subtract(i),
		Duration: n.Duration,
	}
}
