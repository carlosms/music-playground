package main

import (
	"bytes"
	"io"
	"math/rand"

	"github.com/hajimehoshi/oto"
)

const (
	sampleRate        = 44100
	channelNum        = 1
	bitDepthInBytes   = 1
	bufferSizeInBytes = 4096
)

func main() {
	p, err := oto.NewPlayer(sampleRate, channelNum, bitDepthInBytes, bufferSizeInBytes)
	if err != nil {
		panic(err)
	}
	defer p.Close()

	// sampleRate is the number of samples played in 1 second. A number of bytes
	// of sampleRate * 2 makes for a duration of 2 seconds
	nBytes := sampleRate * 2
	buf := make([]byte, nBytes)
	for i := 0; i < nBytes; i++ {
		buf[i] = byte(rand.Intn(256))
	}

	if _, err := io.Copy(p, bytes.NewReader(buf)); err != nil {
		panic(err)
	}
}
