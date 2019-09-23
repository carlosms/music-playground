package note_test

import (
	"fmt"
	"testing"

	"github.com/carlosms/music-playground/theory/note"
	"github.com/stretchr/testify/assert"
)

func TestNoteString(t *testing.T) {
	n := note.NewNote(note.Fsharp5, note.Quarter)

	assert.Equal(t, "♩ F#5", fmt.Sprint(n))
	assert.Equal(t, "♩ F#5", n.String())
}

func TestRestString(t *testing.T) {
	n := note.NewRest(note.Eighth)

	assert.Equal(t, "𝄾", fmt.Sprint(n))
	assert.Equal(t, "𝄾", n.String())

	n = note.NewRest(note.Half)
	assert.Equal(t, "𝄼", fmt.Sprint(n))
	assert.Equal(t, "𝄼", n.String())

	n = note.NewRest(note.Sixteenth / 2)
	assert.Equal(t, "1/32 rest", fmt.Sprint(n))
	assert.Equal(t, "1/32 rest", n.String())

	n = note.NewRest(note.Whole * 4)
	assert.Equal(t, "4 rest", fmt.Sprint(n))
	assert.Equal(t, "4 rest", n.String())

	assert.Equal(t, "", n.Pitch.String())
}

func TestNoteFrequency(t *testing.T) {
	n := note.NewNote(note.C4, note.Half)

	assert.Equal(t, 261.6255653005986, n.Frequency())

	c5 := n.Add(note.Octave)
	assert.Equal(t, 523.2511306011972, c5.Frequency())
	assert.Equal(t, note.Half, c5.Duration)

	a3 := n.Subtract(note.Tone + note.Semitone)
	assert.Equal(t, 220.0, a3.Frequency())
	assert.Equal(t, note.Half, a3.Duration)
}

func TestNoteDuration(t *testing.T) {
	n := note.NewNote(note.C4, note.Half)
	assert.Equal(t, note.Half, n.Duration)

	shorter := n
	shorter.Duration /= 2
	assert.Equal(t, n.Frequency(), shorter.Frequency())
	assert.Equal(t, note.Quarter, shorter.Duration)

	longer := n
	longer.Duration *= 2
	assert.Equal(t, n.Frequency(), longer.Frequency())
	assert.Equal(t, note.Whole, longer.Duration)
}

func TestRestFrequency(t *testing.T) {
	n := note.NewRest(note.Sixteenth)

	assert.Equal(t, 0.0, n.Frequency())

	higher := n.Add(note.Octave)
	assert.Equal(t, 0.0, higher.Frequency())

	lower := n.Subtract(note.Tone + note.Semitone)
	assert.Equal(t, 0.0, lower.Frequency())
}

func TestRestDuration(t *testing.T) {
	n := note.NewRest(note.Whole)
	assert.Equal(t, note.Whole, n.Duration)

	shorter := n
	shorter.Duration /= 2
	assert.Equal(t, n.Frequency(), shorter.Frequency())
	assert.Equal(t, note.Half, shorter.Duration)

	longer := n
	longer.Duration *= 2
	assert.Equal(t, n.Frequency(), longer.Frequency())
	assert.Equal(t, note.Double, longer.Duration)
}
