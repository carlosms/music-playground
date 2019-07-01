package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/carlosms/music-playground/synth"
	"github.com/carlosms/music-playground/theory/note"

	"github.com/hajimehoshi/oto"
)

const (
	sampleRate        = 44100
	channelNum        = 1
	bitDepthInBytes   = 2
	bufferSizeInBytes = 5120
)

func majorScale(n note.Pitch) []note.Pitch {
	// whole, whole, half, whole, whole, whole, half
	scale := make([]note.Pitch, 8)
	scale[0] = n
	scale[1] = scale[0].Add(note.Tone)
	scale[2] = scale[1].Add(note.Tone)
	scale[3] = scale[2].Add(note.Semitone)
	scale[4] = scale[3].Add(note.Tone)
	scale[5] = scale[4].Add(note.Tone)
	scale[6] = scale[5].Add(note.Tone)
	scale[7] = scale[6].Add(note.Semitone)

	return scale
}

func majorChord(n note.Pitch) []note.Pitch {
	scale := majorScale(n)
	return []note.Pitch{scale[0], scale[2], scale[4]}
}

func plotChord(wave synth.WaveGenerator, pitches []note.Pitch) {
	names := make([]string, len(pitches))
	for i, p := range pitches {
		names[i] = p.String()
	}

	c := synth.Combine(
		synth.Sustain(wave(sampleRate, pitches[0].Frequency(), 2*time.Second), 0.3),
		synth.Sustain(wave(sampleRate, pitches[1].Frequency(), 2*time.Second), 0.3),
		synth.Sustain(wave(sampleRate, pitches[2].Frequency(), 2*time.Second), 0.3),
	)
	fmt.Println(synth.Plot(c, strings.Join(names, ","), 0.1))
	fmt.Println()
}

func play(wave synth.WaveGenerator, triad []note.Pitch) io.Reader {
	return synth.Combine(
		synth.Sustain(wave(sampleRate, triad[0].Frequency(), 1*time.Second), 0.2),
		io.MultiReader(
			synth.Sustain(wave(sampleRate, triad[1].Frequency(), 200*time.Millisecond), 0),
			synth.Sustain(wave(sampleRate, triad[1].Frequency(), 1*time.Second), 0.2),
		),
		io.MultiReader(
			synth.Sustain(wave(sampleRate, triad[2].Frequency(), 400*time.Millisecond), 0),
			synth.Sustain(wave(sampleRate, triad[2].Frequency(), 1*time.Second), 0.2),
		),
	)
}

func main() {
	p, err := oto.NewPlayer(sampleRate, channelNum, bitDepthInBytes, bufferSizeInBytes)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	pitch := note.C3

	for i := uint8(0); i < 3; i++ {
		triad := majorChord(pitch.Add(i * note.Octave))
		plotChord(synth.NewSineWave, triad)

		sound := play(synth.NewSineWave, triad)
		if _, err := io.Copy(p, sound); err != nil {
			panic(err)
		}
	}
}
