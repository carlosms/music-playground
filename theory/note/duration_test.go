package note_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/carlosms/music-playground/theory/note"
	"github.com/stretchr/testify/assert"
)

func TestToSeconds(t *testing.T) {
	n := note.Note{Duration: note.Quarter}

	assert.Equal(t, 500*time.Millisecond, n.ToSeconds(note.Quarter, 120))
	assert.Equal(t, time.Second, n.ToSeconds(note.Eighth, 120))
	assert.Equal(t, 250*time.Millisecond, n.ToSeconds(note.Half, 120))
}

func TestDurationString(t *testing.T) {
	assert.Equal(t, `𝅝`, fmt.Sprint(note.Whole))
	assert.Equal(t, `♪`, fmt.Sprint(note.Quarter/2))
	assert.Equal(t, `♩`, fmt.Sprint(note.Eighth*2))

	assert.Equal(t, `1/32 note`, fmt.Sprint(note.Sixteenth/2))
	assert.Equal(t, `8 note`, fmt.Sprint(note.Double*4))

	assert.Equal(t, `𝅜`, fmt.Sprint(note.Double))
	assert.Equal(t, `𝅝`, fmt.Sprint(note.Whole))
	assert.Equal(t, `𝅗𝅥`, fmt.Sprint(note.Half))
	assert.Equal(t, `♩`, fmt.Sprint(note.Quarter))
	assert.Equal(t, `♪`, fmt.Sprint(note.Eighth))
	assert.Equal(t, `𝅘𝅥𝅯`, fmt.Sprint(note.Sixteenth))
}

func TestDurationStringRest(t *testing.T) {
	assert.Equal(t, `𝄻`, note.Whole.StringRest())
	assert.Equal(t, `𝄾`, (note.Quarter / 2).StringRest())
	assert.Equal(t, `𝄽`, (note.Eighth * 2).StringRest())

	assert.Equal(t, `1/32 rest`, (note.Sixteenth / 2).StringRest())
	assert.Equal(t, `8 rest`, (note.Double * 4).StringRest())

	assert.Equal(t, `𝄺`, note.Double.StringRest())
	assert.Equal(t, `𝄻`, note.Whole.StringRest())
	assert.Equal(t, `𝄼`, note.Half.StringRest())
	assert.Equal(t, `𝄽`, note.Quarter.StringRest())
	assert.Equal(t, `𝄾`, note.Eighth.StringRest())
	assert.Equal(t, `𝄿`, note.Sixteenth.StringRest())
}
