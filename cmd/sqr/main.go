package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"time"

	"github.com/carlosms/asciigraph"
	"github.com/hajimehoshi/oto"
)

const (
	sampleRate        = 44100
	channelNum        = 1
	bitDepthInBytes   = 2
	bufferSizeInBytes = 512

	equilibrium int16 = 0
)

var max = int16((math.Pow(2, 8*bitDepthInBytes) - 1) / 2) // (2^16 - 1) / 2

func sqrWave(freq float64, amplitude int16, duration time.Duration) io.Reader {
	if amplitude > max {
		panic(fmt.Sprintf("wrong value %v for amplitude, max value is %v", amplitude, max))
	}

	samplesPeriod := int(sampleRate / freq)
	samplesHalfPeriod := samplesPeriod / 2

	nBytes := bitDepthInBytes * int(sampleRate*duration.Seconds())
	buf := make([]byte, nBytes)

	for i := 0; i < nBytes; i += bitDepthInBytes {
		pos := (i / bitDepthInBytes) % samplesPeriod

		var value int16
		if pos <= samplesHalfPeriod {
			value = equilibrium + amplitude
		} else {
			value = equilibrium - amplitude
		}

		// little-endian
		buf[i] = byte(value)
		if bitDepthInBytes == 2 {
			buf[i+1] = byte(value >> 8)
		}
	}

	return bytes.NewReader(buf)
}

func plot(freq float64, amplitude int16, duration time.Duration) string {
	r := sqrWave(freq, amplitude, duration)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	var plotData []float64
	for i := 0; i < len(data); i += bitDepthInBytes {
		switch bitDepthInBytes {
		case 1:
			plotData = append(plotData, float64(data[i]))
		case 2:
			// little-endian
			plotData = append(plotData, float64(int16(data[i])+int16(data[i+1])<<8))
		}
	}

	caption := fmt.Sprintf("f = %.3f kHz, A = %d", freq/1000, amplitude)

	return asciigraph.Plot(plotData,
		asciigraph.Caption(caption),
		asciigraph.Height(15), asciigraph.Min(float64(-max)), asciigraph.Max(float64(max)))
}

func main() {
	fmt.Println(plot(2000, 12800, time.Millisecond))
	fmt.Println()
	fmt.Println(plot(6000, 25600, time.Millisecond))

	r := io.MultiReader(
		sqrWave(100, 16383, time.Second/2),
		sqrWave(300, 6553, time.Second/2),
		sqrWave(600, 6553, time.Second),
		sqrWave(1, 0, time.Second/3),
		sqrWave(400, 9830, time.Second/2),
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
