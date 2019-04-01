package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/carlosms/asciigraph"
	"github.com/hajimehoshi/oto"
)

const (
	sampleRate        = 44100
	channelNum        = 1
	bitDepthInBytes   = 1
	bufferSizeInBytes = 4096
)

func sqrWave(freq, amplitude int, duration time.Duration) io.Reader {
	const equilibrium = 127

	samplesPeriod := sampleRate / freq
	samplesHalfPeriod := samplesPeriod / 2

	nBytes := int(sampleRate * duration.Seconds())
	buf := make([]byte, nBytes)

	for i := 0; i < nBytes; i++ {
		even := (i/samplesHalfPeriod)%2 == 0

		if even {
			buf[i] = byte(equilibrium + amplitude)
		} else {
			buf[i] = byte(equilibrium - amplitude)
		}
	}

	return bytes.NewReader(buf)
}

func plot(freq, amplitude int, duration time.Duration) string {
	r := sqrWave(freq, amplitude, duration)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	var plotData []float64
	for _, b := range data {
		plotData = append(plotData, float64(b))
	}

	caption := fmt.Sprintf("f = %d kHz, A = %d", freq/1000, amplitude)

	return asciigraph.Plot(plotData,
		asciigraph.Caption(caption),
		asciigraph.Height(15), asciigraph.Min(0), asciigraph.Max(255))
}

func main() {
	fmt.Println(plot(2000, 50, time.Millisecond))
	fmt.Println()
	fmt.Println(plot(6000, 100, time.Millisecond))

	r := io.MultiReader(
		sqrWave(100, 50, time.Second/2),
		sqrWave(300, 20, time.Second/2),
		sqrWave(600, 20, time.Second),
		sqrWave(1, 0, time.Second/3),
		sqrWave(400, 30, time.Second/2),
	)

	p, err := oto.NewPlayer(sampleRate, channelNum, bitDepthInBytes, bufferSizeInBytes)
	if err != nil {
		panic(err)
	}

	defer p.Close()

	if _, err := io.Copy(p, r); err != nil {
		panic(err)
	}
}
