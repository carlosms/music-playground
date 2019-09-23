package main

import (
	"io"

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

func play(wave synth.WaveGenerator, bpm int, staff [][]note.Note) io.Reader {
	readers := []io.Reader{}

	for _, notes := range staff {
		groupReaders := []io.Reader{}

		for _, n := range notes {
			d := n.ToSeconds(note.Quarter, bpm)
			w := wave(sampleRate, n.Frequency(), d)
			r := synth.Sustain(w, 0.2)
			groupReaders = append(groupReaders, r)
		}

		readers = append(readers, synth.Combine(groupReaders...))
	}

	return io.MultiReader(readers...)
}

func main() {
	p, err := oto.NewPlayer(sampleRate, channelNum, bitDepthInBytes, bufferSizeInBytes)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	sound := synth.Combine(
		play(synth.NewSineWave, 180, trebleStaff),
		play(synth.NewSineWave, 180, bassStaff),
	)

	if _, err := io.Copy(p, sound); err != nil {
		panic(err)
	}
}

// Marble Machine
// Composed by Wintergatan, transcribed by Chalmers Huang
// Transcribed to Go painstakingly manually from
// https://musescore.com/user/5631216/scores/1846226
var trebleStaff = [][]note.Note{
	// Bar 1
	[]note.Note{note.NewNote(note.E6, note.Quarter)},
	[]note.Note{note.NewNote(note.E5, note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{note.NewNote(note.E5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	// Bar 2
	[]note.Note{note.NewNote(note.G5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.E5, note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewNote(note.G5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.D6, note.Eighth)},
	// Bar 3
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{note.NewNote(note.E5, note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{note.NewNote(note.E5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	// Bar 4
	[]note.Note{note.NewNote(note.G5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.D5, note.Eighth)},
	[]note.Note{note.NewNote(note.Fsharp5, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewNote(note.G5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.D6, note.Eighth)},
	// Bar 5
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{note.NewNote(note.Fsharp5, note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{note.NewNote(note.Fsharp5, note.Eighth)},
	[]note.Note{note.NewNote(note.D6, note.Eighth)},
	// Bar 6
	[]note.Note{note.NewNote(note.C6, note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewNote(note.Fsharp5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewNote(note.G5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.E5, note.Eighth)},
	// Bar 7
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewNote(note.C5, note.Eighth)},
	[]note.Note{note.NewNote(note.E5, note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewNote(note.B4, note.Eighth)},
	[]note.Note{note.NewNote(note.C5, note.Eighth)},
	[]note.Note{note.NewNote(note.D5, note.Eighth)},
	[]note.Note{note.NewNote(note.D6, note.Eighth)},
	// Bar 8
	[]note.Note{note.NewNote(note.C6, note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewNote(note.Fsharp5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewNote(note.G5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.E6, note.Eighth)},
	// Bar 9
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{note.NewNote(note.E5, note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{note.NewNote(note.E5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	// Bar 10
	[]note.Note{note.NewNote(note.G5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.E5, note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewNote(note.G5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.D6, note.Eighth)},
	// Bar 11
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{note.NewNote(note.Fsharp5, note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{note.NewNote(note.Fsharp5, note.Eighth)},
	[]note.Note{note.NewNote(note.D6, note.Eighth)},
	// Bar 12
	[]note.Note{note.NewNote(note.C6, note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewNote(note.Fsharp5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewNote(note.G5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.E6, note.Eighth)},
	// Bar 13
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{note.NewNote(note.Fsharp5, note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.E6, note.Eighth)},
	// Bar 14
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewNote(note.B5, note.Eighth)},
	[]note.Note{note.NewNote(note.Fsharp5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewNote(note.G5, note.Eighth)},
	[]note.Note{note.NewNote(note.F5, note.Eighth)},
	[]note.Note{note.NewNote(note.E5, note.Eighth)},
	// Bar 15
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewNote(note.B4, note.Eighth)},
	[]note.Note{note.NewNote(note.C5, note.Eighth)},
	[]note.Note{note.NewNote(note.Fsharp5, note.Eighth)},
	[]note.Note{note.NewNote(note.C5, note.Eighth)},
	[]note.Note{note.NewNote(note.E5, note.Eighth)},
	[]note.Note{note.NewNote(note.G5, note.Eighth)},
	[]note.Note{note.NewNote(note.D5, note.Eighth)},
	// Bar 16
	[]note.Note{note.NewNote(note.D5, note.Eighth)},
	[]note.Note{note.NewNote(note.Fsharp5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.B4, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.D5, note.Eighth)},
	[]note.Note{note.NewNote(note.G5, note.Eighth)},
	[]note.Note{note.NewNote(note.A5, note.Eighth)},
	[]note.Note{note.NewNote(note.E6, note.Eighth)},
}

var bassStaff = [][]note.Note{
	// Bar 1
	[]note.Note{note.NewRest(note.Whole)},
	// Bar 2
	[]note.Note{note.NewRest(note.Whole)},
	// Bar 3
	[]note.Note{note.NewRest(note.Whole)},
	// Bar 4
	[]note.Note{note.NewRest(note.Whole)},
	// Bar 5
	[]note.Note{note.NewRest(note.Whole)},
	// Bar 6
	[]note.Note{note.NewRest(note.Whole)},
	// Bar 7
	[]note.Note{note.NewRest(note.Whole)},
	// Bar 8
	[]note.Note{note.NewRest(note.Whole)},
	// Bar 9
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.E4, note.Eighth),
		note.NewNote(note.G4, note.Eighth),
		note.NewNote(note.B4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.E4, note.Eighth),
		note.NewNote(note.G4, note.Eighth),
		note.NewNote(note.B4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	// Bar 10
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.E4, note.Eighth),
		note.NewNote(note.G4, note.Eighth),
		note.NewNote(note.B4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.E4, note.Eighth),
		note.NewNote(note.G4, note.Eighth),
		note.NewNote(note.B4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	// Bar 11
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.D4, note.Eighth),
		note.NewNote(note.Fsharp4, note.Eighth),
		note.NewNote(note.A4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.D4, note.Eighth),
		note.NewNote(note.Fsharp4, note.Eighth),
		note.NewNote(note.A4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	// Bar 12
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.D4, note.Eighth),
		note.NewNote(note.Fsharp4, note.Eighth),
		note.NewNote(note.A4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.D4, note.Eighth),
		note.NewNote(note.Fsharp4, note.Eighth),
		note.NewNote(note.A4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	// Bar 13
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.B3, note.Eighth),
		note.NewNote(note.D4, note.Eighth),
		note.NewNote(note.Fsharp4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.B3, note.Eighth),
		note.NewNote(note.D4, note.Eighth),
		note.NewNote(note.Fsharp4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	// Bar 14
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.B3, note.Eighth),
		note.NewNote(note.D4, note.Eighth),
		note.NewNote(note.Fsharp4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.B3, note.Eighth),
		note.NewNote(note.D4, note.Eighth),
		note.NewNote(note.Fsharp4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	// Bar 15
	[]note.Note{
		note.NewNote(note.C3, note.Eighth),
		note.NewNote(note.E4, note.Eighth),
		note.NewNote(note.G4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	[]note.Note{
		note.NewNote(note.C3, note.Eighth),
		note.NewNote(note.E4, note.Eighth),
		note.NewNote(note.G4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewRest(note.Quarter)},
	// Bar 16
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{
		note.NewNote(note.C3, note.Eighth),
		note.NewNote(note.E4, note.Eighth),
		note.NewNote(note.G4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{
		note.NewNote(note.C3, note.Eighth),
		note.NewNote(note.E4, note.Eighth),
		note.NewNote(note.G4, note.Eighth)},
	[]note.Note{
		note.NewNote(note.C3, note.Eighth),
		note.NewNote(note.E4, note.Eighth),
		note.NewNote(note.G4, note.Eighth)},
	[]note.Note{note.NewRest(note.Eighth)},
	[]note.Note{
		note.NewNote(note.C3, note.Eighth),
		note.NewNote(note.E4, note.Eighth),
		note.NewNote(note.G4, note.Eighth)},
}
