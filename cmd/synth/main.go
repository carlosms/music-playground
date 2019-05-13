package main

import (
	"fmt"
	"io"
	"time"

	"github.com/carlosms/music-playground/synth"
	"github.com/hajimehoshi/oto"
)

const (
	sampleRate        = 44100
	channelNum        = 1
	bitDepthInBytes   = 2
	bufferSizeInBytes = 5120
)

const (
	c4 = 261.6256
	e4 = 329.6276
	g4 = 391.9954
)

func plotChord(wave synth.WaveGenerator) {
	fmt.Println(synth.Plot(
		synth.Sustain(wave(sampleRate, c4, 2*time.Second), 0.3),
		fmt.Sprintf("C4 f = %.3f kHz, Sustain = %.2f", c4, 0.3), 0.1))
	fmt.Println()
	fmt.Println(synth.Plot(
		synth.Sustain(wave(sampleRate, e4, 2*time.Second), 0.3),
		fmt.Sprintf("E4 f = %.3f kHz, Sustain = %.2f", e4, 0.3), 0.1))
	fmt.Println()
	fmt.Println(synth.Plot(
		synth.Sustain(wave(sampleRate, g4, 2*time.Second), 0.3),
		fmt.Sprintf("G4 f = %.3f kHz, Sustain = %.2f", g4, 0.3), 0.1))
	fmt.Println()

	c := synth.Combine(
		synth.Sustain(wave(sampleRate, c4, 2*time.Second), 0.3),
		synth.Sustain(wave(sampleRate, e4, 2*time.Second), 0.3),
		synth.Sustain(wave(sampleRate, g4, 2*time.Second), 0.3),
	)
	fmt.Println(synth.Plot(c, "C major chord", 0.1))
	fmt.Println()
}

func chord(wave synth.WaveGenerator) io.Reader {
	return synth.Combine(
		synth.Sustain(wave(sampleRate, c4, 2*time.Second), 0.3),
		io.MultiReader(
			synth.Sustain(wave(sampleRate, e4, 400*time.Millisecond), 0),
			synth.Sustain(wave(sampleRate, e4, 2*time.Second), 0.3),
		),
		io.MultiReader(
			synth.Sustain(wave(sampleRate, g4, 800*time.Millisecond), 0),
			synth.Sustain(wave(sampleRate, g4, 2*time.Second), 0.3),
		),
	)
}

func main() {
	p, err := oto.NewPlayer(sampleRate, channelNum, bitDepthInBytes, bufferSizeInBytes)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	fmt.Println("C major, sine wave")
	fmt.Println("--------------------")
	plotChord(synth.NewSineWave)

	sound := chord(synth.NewSineWave)
	if _, err := io.Copy(p, sound); err != nil {
		panic(err)
	}

	time.Sleep(500 * time.Millisecond)

	fmt.Println("C major, square wave")
	fmt.Println("--------------------")
	plotChord(synth.NewSquareWave)

	sound = synth.Sustain(chord(synth.NewSquareWave), 0.6)
	if _, err := io.Copy(p, sound); err != nil {
		panic(err)
	}
}
