package note

import (
	"fmt"
	"time"
)

// Duration is the duration relative to the whole note
type Duration float64

const (
	// Double is 𝅜
	Double Duration = 2
	// Whole is 𝅝
	Whole Duration = 1
	// Half is 𝅗𝅥
	Half Duration = 0.5
	// Quarter is ♩
	Quarter Duration = 0.25
	// Eighth is ♪
	Eighth Duration = 0.125
	// Sixteenth is 𝅘𝅥𝅯
	Sixteenth Duration = 0.0625
)

// ToSeconds returns the note value as time.Duration. tempoNote is the note in
// the tempo marking, or the _beat_ in _beats per minute_.
func (d Duration) ToSeconds(tempoNote Duration, bpm int) time.Duration {
	// duration for each tempoNote
	timeBeat := float64(time.Minute) / float64(bpm)
	// convert to the d note value
	t := float64(d) * (timeBeat / float64(tempoNote))

	return time.Duration(t)
}

// String returns the musical note symbol
func (d Duration) String() string {
	switch d {
	case Double:
		return "𝅜"
	case Whole:
		return "𝅝"
	case Half:
		return "𝅗𝅥"
	case Quarter:
		return "♩"
	case Eighth:
		return "♪"
	case Sixteenth:
		return "𝅘𝅥𝅯"
	default:
		var s string
		f := float64(d)
		if f < 1 {
			s = fmt.Sprintf("1/%v note", 1/f)
		} else {
			s = fmt.Sprintf("%v note", f)
		}
		return s
	}
}

// StringRest returns the rest note symbol
func (d Duration) StringRest() string {
	switch d {
	case Double:
		return "𝄺"
	case Whole:
		return "𝄻"
	case Half:
		return "𝄼"
	case Quarter:
		return "𝄽"
	case Eighth:
		return "𝄾"
	case Sixteenth:
		return "𝄿"
	default:
		var s string
		f := float64(d)
		if f < 1 {
			s = fmt.Sprintf("1/%v rest", 1/f)
		} else {
			s = fmt.Sprintf("%v rest", f)
		}
		return s
	}
}
