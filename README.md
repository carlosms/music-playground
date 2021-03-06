# Music Playground

This repository is a kind of notebook where I'll experiment programing music theory concepts from scratch, using Go.

This is not supposed to be a tutorial, and I don't claim to have any authority in this subject. I am deliberately avoiding to use any code as reference for the implementation of a synthesizer, or musical notation. There might be better ways to implement music concepts, but the goal is for me to figure out how to achieve it.

## Table of Contents

- [1 Let There Be Noise](#1-let-there-be-noise)
  - [1.1 PCM](#11-pcm)
  - [1.2 Beep Boop](#12-beep-boop)
  - [1.3 Pure Waves](#13-pure-waves)
  - [1.4 World's Smallest Synthesizer](#14-worlds-smallest-synthesizer)
- [2 Noteworthy Theory](#2-noteworthy-theory)
  - [2.1 Do, a Deer, a Female Deer](#21-do-a-deer-a-female-deer)
    - [2.1.1 The Theory](#211-the-theory)
    - [2.1.2 The Implementation](#212-the-implementation)
  - [2.2 Pata-Pata-Pata-Pon](#22-pata-pata-pata-pon)
  - [2.3 The Sound... of Silence](#23-the-sound-of-silence)
  - [2.4 Lost My Marbles](#24-lost-my-marbles)
- [3 Harmony](#3-harmony)
  - [3.1 Scaling Up](#31-scaling-up)

## 1 Let There Be Noise

### 1.1 PCM

The starting point of this repository is [Oto (音)](https://github.com/hajimehoshi/oto), a low-level Golang library to play PCM sound.

PCM stands for **Pulse-code modulation**, and [Wikipedia defines it](https://en.wikipedia.org/wiki/Pulse-code_modulation) as:

> Pulse-code modulation (PCM) is a method used to digitally represent sampled analog signals. It is the standard form of digital audio in computers, compact discs, digital telephony and other digital audio applications. In a PCM stream, the amplitude of the analog signal is sampled regularly at uniform intervals, and each sample is quantized to the nearest value within a range of digital steps.
> 
> A PCM stream has two basic properties that determine the stream's fidelity to the original analog signal: the sampling rate, which is the number of times per second that samples are taken; and the bit depth, which determines the number of possible digital values that can be used to represent each sample.

To simplify the first examples I will constraint the data to have 8 **bit depth**, and 1 **channel** (mono audio). The **sample rate** used will be 44100 Hz, a typical value used for example in CD audio.

```go
const (
	sampleRate        = 44100
	channelNum        = 1
	bitDepthInBytes   = 1
)
```

The Oto player can be initialized like this:

```go
p, _ := oto.NewPlayer(sampleRate, channelNum, bitDepthInBytes, bufferSizeInBytes)
```

[`oto.Player`](https://godoc.org/github.com/hajimehoshi/oto#Player) implements [`io.WriteCloser`](https://golang.org/pkg/io/#WriteCloser). To play sound we just need to use the `Write` method to send the samples.

We need to provide sound samples as an stream of bytes. The library expects the data to follow this sequence:

```
[data]      = [sample 1] [sample 2] [sample 3] ...
[sample *]  = [channel 1] ...
[channel *] = [byte 1] [byte 2] ...
```

First we need to allocate some bytes. The number of bytes we need depends on the duration of the sound we want to play, the number of channels, the **bit depth**, and the **sample rate**.

With the values defined above we only need to care about the **sample rate** (44100 Hz) and the duration. For example 2 seconds will require an array of `2 * 44100` bytes.

```go
nBytes := sampleRate * 2
buf := make([]byte, nBytes)
```

Because the **bit depth** is 1 byte, each byte holds the value of one sample. This value can be in the range `[0x00..0xFF]`, so we can fill them with random data.

```go
for i := 0; i < nBytes; i++ {
  buf[i] = byte(rand.Intn(256))
}
```

The last step is to send the data to the `oto.Player` `Write` method.

```go
io.Copy(p, bytes.NewReader(buf))
```

You can find the complete example in [cmd/noise/main.go](./cmd/noise/main.go). If you run it you will be rewarded with the beautiful sound of nostalgia, in the form of TV static noise.

#### [cmd/noise/main.go](./cmd/noise/main.go)

```shell
$ go run cmd/noise/main.go
```

### 1.2 Beep Boop

Sound is produced by the the propagation of a vibration through air.

From the code in this repository we are instructing the computer to move the diaphragm of the speaker back and forth quickly. This diaphragm oscillation causes the air in front of the speaker to move also back and forth, in turn pushing the air in front of it to also move, and so on until this movement reaches the air in contact with our eardrums.

So the vibration that reaches our ear is produced by a wave of pressure. Let's try to oscillate the speaker at a constant rate to see how it sounds.

To do that we will oscillate the PCM sample byte values between 2 fixed values, at a fixed rate. Or, put another way, we will create a continuous sound wave.

A sound wave can be simplified to the following parameters:

- **Amplitude** _A_: the maximum distance that air molecules are displaced from their neutral position. Measured in meters.
- **Frequency** _f_: number of wave oscillation cycles per second. It is measured in [**hertz (Hz)**](https://en.wikipedia.org/wiki/Hertz). 1 Hz means 1 wave cycle per second.
- **Period** _T_: time it takes for one wave cycle to complete.
  
  Period and frequency are related:

  _f_ = 1 / _T_
  
  This means for example that a wave with a period of 500ms has a frequency of 2Hz.
- **Speed** _v_: how fast the wave travels, measured in meters per second. For sound waves the speed depends on the air temperature. For 20°C the speed is ~343m/s.
- **Wavelength** λ: distance, in meters, between 2 points of air in the identical part of an oscillation cycle, e.g. 2 wave crests. The wavelength is related to the frequency:
  
  λ = _v_ / _f_

```
Y: Displacement of air molecules
┤
┤               ◄---------------------------- T ---------------------------►
┤          ╭────────╮                                                  ╭────────╮
┤       ╭──╯    ▲   ╰──╮                                            ╭──╯        ╰──╮
┤    ╭──╯     A |      ╰─╮                                       ╭──╯              ╰─╮
┤  ╭─╯          |        ╰──╮                                  ╭─╯                   ╰──╮
┤╭─╯            ▼           ╰─╮                              ╭─╯                        ╰─╮
┼╯ ------ equilibrium -----   ╰╮                            ╭╯                            ╰╮
┤          position            ╰─╮                        ╭─╯                              ╰─╮
┤                                ╰─╮                   ╭──╯                                  ╰─╮
┤                                  ╰──╮              ╭─╯                                       ╰─
┤                                     ╰──╮        ╭──╯                                           
┤                                        ╰────────╯                                              
 ┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬────┬
 X: Time
```

Knowing this we can create a function to oscillate the values of the PCM samples back and forth at a fixed rate, creating a square wave.
The wave we are creating is a digital wave, so measurement of of the **amplitude** is not expressed in meters. Instead, it is a number of discrete steps.
Because the **bit depth** will be set to 16 bits, the value of each sample can range from `-32768` to `32767`. The middle point of the wave, or **equilibrium**, will be `0`. This means the amplitude must have a value in the range `[0..32767]`.

The function will simply divide the period _T_ in 2, and set `equilibrium + amplitude` or `equilibrium - amplitude`, depending on which half of the period the sample index falls on.

```go
const (
	sampleRate        = 44100
	bitDepthInBytes   = 2

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
```

This `io.Reader` can be consumed by the `oto.Player` like this:

```go
p, _ := oto.NewPlayer(sampleRate, channelNum, bitDepthInBytes, bufferSizeInBytes)

io.Copy(p, sqrWave(600, 25600, time.Second))
```

We can also print an ASCII graph of the sampled values using [github.com/guptarohit/asciigraph](https://github.com/guptarohit/asciigraph). Or, rather, [my fork](https://github.com/carlosms/asciigraph) that adds a couple of new options.

```go
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
```

```
  32767 ┤                                            
  28671 ┤                                            
  24575 ┤                                            
  20479 ┤                                            
  16384 ┤                                            
  12288 ┼───────────╮         ╭───────────╮          
   8192 ┤           │         │           │          
   4096 ┤           │         │           │          
      0 ┼           │         │           │          
  -4096 ┤           │         │           │          
  -8192 ┤           │         │           │          
 -12288 ┤           ╰─────────╯           ╰───────── 
 -16384 ┤                                            
 -20479 ┤                                            
 -24575 ┤                                            
 -28671 ┤                                            
 -32767 ┤                                            
           f = 2.000 kHz, A = 12800

  32767 ┤                                            
  28671 ┤                                            
  24575 ┼───╮  ╭───╮  ╭───╮  ╭───╮  ╭───╮  ╭───╮  ╭─ 
  20479 ┤   │  │   │  │   │  │   │  │   │  │   │  │  
  16384 ┤   │  │   │  │   │  │   │  │   │  │   │  │  
  12288 ┤   │  │   │  │   │  │   │  │   │  │   │  │  
   8192 ┤   │  │   │  │   │  │   │  │   │  │   │  │  
   4096 ┤   │  │   │  │   │  │   │  │   │  │   │  │  
      0 ┼   │  │   │  │   │  │   │  │   │  │   │  │  
  -4096 ┤   │  │   │  │   │  │   │  │   │  │   │  │  
  -8192 ┤   │  │   │  │   │  │   │  │   │  │   │  │  
 -12288 ┤   │  │   │  │   │  │   │  │   │  │   │  │  
 -16384 ┤   │  │   │  │   │  │   │  │   │  │   │  │  
 -20479 ┤   │  │   │  │   │  │   │  │   │  │   │  │  
 -24575 ┤   ╰──╯   ╰──╯   ╰──╯   ╰──╯   ╰──╯   ╰──╯  
 -28671 ┤                                            
 -32767 ┤                                            
           f = 6.000 kHz, A = 25600
```

You can find the complete example in [cmd/sqr/main.go](./cmd/sqr/main.go).
This program plays a few random robot noises that could be sound effects from an old school Atari game.

#### [cmd/sqr/main.go](./cmd/sqr/main.go)

```shell
$ go run cmd/sqr/main.go
```

### 1.3 Pure Waves

Square waves are OK, and you can create amazing music such as chiptune.
But if we want a more _musical_ or _natural_ sound we can create sine waves.

Pure sine waves do not really replicate the feel of any instrument, but we're getting closer. This is the kind of wave that a tuning fork produces.

We can create one with the `math.Sin` function, which takes an argument in radians. To do that first we need to know the offset of the current sample, or number of samples since the first one. If we divide this sample offset by the number of samples in a period, we have the position in the horizontal axis of the wave. Knowing that a full wave **period** is `2 * π` radians, we can then find out the sample position in radians by multiplying.

```go
const (
	sampleRate        = 44100
	bitDepthInBytes   = 2

	equilibrium int16 = 0
)

var max = int16((math.Pow(2, 8*bitDepthInBytes) - 1) / 2) // (2^16 - 1) / 2

func sineWave(freq float64, amplitude int16, duration time.Duration) io.Reader {
	if amplitude > max {
		panic("amplitude max value is 32767")
	}

	samplesPeriod := int(sampleRate / freq)

	nBytes := bitDepthInBytes * int(sampleRate*duration.Seconds())
	buf := make([]byte, nBytes)

	for i := 0; i < nBytes; i += bitDepthInBytes {
		offset := (i / bitDepthInBytes)

		radian := float64(offset) / float64(samplesPeriod) * 2 * math.Pi
		value := equilibrium + int16(math.Round(float64(amplitude)*math.Sin(radian)))

		// little-endian
		buf[i] = byte(value)
		if bitDepthInBytes == 2 {
			buf[i+1] = byte(value >> 8)
		}
	}

	return bytes.NewReader(buf)
}
```

Which, visualized with the previous `plot` function produces this output:

```
  32767 ┤                                            
  28671 ┤                                            
  24575 ┤                                            
  20479 ┤                                            
  16384 ┤                                            
  12288 ┤   ╭───╮                 ╭───╮              
   8192 ┤ ╭─╯   ╰─╮             ╭─╯   ╰─╮            
   4096 ┤╭╯       ╰╮           ╭╯       ╰╮           
      0 ┼╯         ╰╮         ╭╯         ╰╮          
  -4096 ┤           ╰╮       ╭╯           ╰╮       ╭ 
  -8192 ┤            ╰─╮   ╭─╯             ╰─╮   ╭─╯ 
 -12288 ┤              ╰───╯                 ╰───╯   
 -16384 ┤                                            
 -20479 ┤                                            
 -24575 ┤                                            
 -28671 ┤                                            
 -32767 ┤                                            
           f = 2.000 kHz, A = 12800

  32767 ┤                                            
  28671 ┤                                            
  24575 ┤  ╭╮          ╭╮          ╭╮          ╭╮    
  20479 ┤ ╭╯╰╮        ╭╯╰╮        ╭╯╰╮        ╭╯╰╮   
  16384 ┤ │  │        │  │        │  │        │  │   
  12288 ┤╭╯  ╰╮      ╭╯  ╰╮      ╭╯  ╰╮      ╭╯  ╰╮  
   8192 ┤│    │      │    │      │    │      │    │  
   4096 ┤│    │      │    │      │    │      │    │  
      0 ┼╯    ╰╮    ╭╯    ╰╮    ╭╯    ╰╮    ╭╯    ╰╮ 
  -4096 ┤      │    │      │    │      │    │      │ 
  -8192 ┤      │    │      │    │      │    │      │ 
 -12288 ┤      ╰╮  ╭╯      ╰╮  ╭╯      ╰╮  ╭╯      ╰ 
 -16384 ┤       │  │        │  │        │  │         
 -20479 ┤       ╰╮╭╯        ╰╮╭╯        ╰╮╭╯         
 -24575 ┤        ╰╯          ╰╯          ╰╯          
 -28671 ┤                                            
 -32767 ┤                                            
           f = 3.500 kHz, A = 25600

```

You can find the complete example in [cmd/sine/main.go](./cmd/sine/main.go).
This program plays a few random noises that sound a bit more pleasant that the previous square waves.

#### [cmd/sine/main.go](./cmd/sine/main.go)

```shell
$ go run cmd/sine/main.go
```

### 1.4 World's Smallest Synthesizer

Now that we have some code in place, it's time to refactor it and stop being too embarrassed about the contents of this public repository.

Everything lives in each command `main.go` file, and that's not good for obvious reasons. We'll move the code dealing with waves to a new package: `synth`.

The `sineWave` function created all the samples at once in a buffer, and then returned a `bytes.NewReader` as a way to satisfy the `io.Reader` interface. It would be better to refactor it into a real `io.Reader` that creates the bytes on demand. Also, to separate concerns, the sine wave generator will not worry about the amplitude, only about the shape of the wave.

```go
// SineWave is an io.Reader that returns int16 samples. The Reader returns
// int16 samples of a sine wave with maximum amplitude.
// Byte ordering is little endian. The format is:
//     [sample 0 byte 0] [sample 0 byte 1] [sample 1 byte 0] [sample 1 byte 1]...
type SineWave struct {
	freq float64
	// nSamples is the total number of samples that can be read
	nSamples int64
	// offset is measured in number of samples read so far
	offset int64

	sampleRate int
}

const (
	equilibrium int16 = 0
	max               = 32767 // (2^16 - 1) / 2
)

func (s *SineWave) Read(buf []byte) (int, error) {
	if s.offset >= s.nSamples {
		return 0, io.EOF
	}

	samplesPeriod := int64(float64(s.sampleRate) / s.freq)

	var i int
	for i = 0; i < len(buf)/2 && s.offset < s.nSamples; i += 2 {
		radian := float64(pos) / float64(samplesPeriod) * 2 * math.Pi
		value := equilibrium + int16(float64(max)*math.Sin(radian))

		// int16 to 2 bytes, little-endian
		buf[i] = byte(value)
		buf[i+1] = byte(value >> 8)

		s.offset++
	}

	if s.offset >= s.nSamples {
		return i, io.EOF
	}

	return i, nil
}
```

You can find the complete code in [./synth/sine.go](./synth/sine.go). [That directory](./synth) contains other files, for example another `io.Reader` for square wave, [./synth/square.go](./synth/square.go).

But what about the volume? As we saw the wave amplitude affects the sound volume, but our `SineWave` and `SquareWave` Readers don't have a parameter to configure that, the waves are created with the maximum amplitude possible.

Well, that's where having `io.Reader`'s comes in handy. We can create a new function that takes those `int16` samples streams and modifies them in some way. In this case, we are interested in the distance between each sample and the equilibrium position. A simple multiplication will allow us to change the volume of any sound wave, without having to know the actual shape of it:

```go
// SustainedReader takes a Reader that returns int16 samples, and multiplies
// each value by the given percentage. The percentage must be a value between
// 0 and 1
type SustainedReader struct {
	r          io.Reader // underlying reader
	percentage float64
}

func (s *SustainedReader) Read(p []byte) (n int, err error) {
	n, err = s.r.Read(p)

	for i := 0; i < n-1; i += 2 {
		// Convert 2 bytes to int16, little-endian
		v := int16(p[i]) + int16(p[i+1])<<8

		v = int16(float64(v) * s.percentage)

		// int16 back to to 2 bytes, little-endian
		p[i] = byte(v)
		p[i+1] = byte(v >> 8)
	}

	return
}
```

You can find the complete code in [./synth/sustain.go](./synth/sustain.go).

There is a reason why this `Reader` is called `Sustain` and not `Volume` or `Amplitude`. This simple method to change the volume of a sound wave opens a new world full of concepts like **envelope**, **ADSR**, **filter**, **LFO**... but let's put aside this melon and we'll open it some other time*.

- _* In Spanish there is an idiom for "abrir un melón"/"to open a melon" that roughly translates to "open a can of worms"_

In previous examples we used an `io.MultiReader` to concatenate `io.Readers`, playing one sound wave after the other. So we could say the code was a **monophonic** synthesizer, able to play one note at a time. Making **polyphonic** music is much more fun, so let's also have some code to mix different samples.

The PCM library used, [Oto](https://github.com/hajimehoshi/oto), already supports sending samples simultaneously. So we could be playing different sounds from a few goroutines. But having my own function to mix the samples coming from different `io.Readers` will allow me to be more independent of the PCM library. This could be used in the future for example to save the final mixed samples into a sound file.

Combining 2 or more samples is just a matter of adding their value. The only tricky thing to take into account is that in Go `int16` may overflow. For example:

```go
	var n int16 = 30000
	var m int16 = 10000
	fmt.Println(n + m) // -25536
```

Knowing this we can create a simple function that checks for overflows, and clips the results to the maximum or minimum value for `int16`:

```go
func add(a, b int16) int16 {
	v := a + b

	signA := a < 0
	signB := b < 0
	signV := v < 0

	// negative + negative = positive, overflow
	if signA && signB && !signV {
		return math.MinInt16
	}

	// positive + positive = negative, overflow
	if !signA && !signB && signV {
		return math.MaxInt16
	}

	return v
}
```

Now for the combined output of more than one stream of samples, we need to be able to return `io.EOF` when all of the individual streams have been fully read.

```go
// Combine takes an arbitrary number of Readers that return int16 samples, and
// returns a Reader that combines those samples into one single sample value
func Combine(readers ...io.Reader) io.Reader {
	return &CombinedReader{readers}
}

// CombinedReader takes an arbitrary number of Readers that return int16 samples,
// and combines those samples into one single sample value
type CombinedReader struct {
	readers []io.Reader // underlying readers
}

func (m *CombinedReader) Read(p []byte) (int, error) {
	var eofErr error
	var n int
	for n = 0; n < len(p)-1; n += 2 {
		var total int16
		eofErr = io.EOF

		for _, r := range m.readers {
			// read 1 sample (2 bytes)
			buf := make([]byte, 2)
			n, err := io.ReadFull(r, buf)

			if n == 2 {
				eofErr = nil

				// Convert 2 bytes to int16, little-endian
				v := int16(buf[0]) + int16(buf[1])<<8
				total = add(total, v)
			}

			// ErrUnexpectedEOF means the reader had less than 2 bytes, we can't
			// use that as a sample so it is also discarded gracefully
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				return n - 2, err
			}
		}

		// int16 back to to 2 bytes, little-endian
		p[n] = byte(total)
		p[n+1] = byte(total >> 8)
	}

	return n, eofErr
}
```

You can find the complete code in [./synth/readers.go](./synth/readers.go).

The `synth` package can be used to play a C major chord like follows (warning, contains spoilers of the future sections!):

```go
const (
	c4 = 261.6256
	e4 = 329.6276
	g4 = 391.9954
)

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
```

```
C major, sine wave
--------------------
  32767 ┤                                                                                           
  28671 ┤                                                                                           
  24575 ┤                                                                                           
  20479 ┤                                                                                           
  16384 ┤                                                                                           
  12288 ┤                                                                                           
   8192 ┤ ╭────╮           ╭────╮           ╭────╮           ╭───╮            ╭───╮           ╭──── 
   4096 ┤╭╯    ╰╮         ╭╯    ╰╮         ╭╯    ╰╮        ╭─╯   ╰─╮        ╭─╯   ╰─╮        ╭╯     
      0 ┼╯      ╰╮       ╭╯      ╰╮       ╭╯      ╰╮      ╭╯       ╰╮      ╭╯       ╰╮      ╭╯      
  -4096 ┤        ╰─╮   ╭─╯        ╰─╮   ╭─╯        ╰╮    ╭╯         ╰╮    ╭╯         ╰╮    ╭╯       
  -8192 ┤          ╰───╯            ╰───╯           ╰────╯           ╰────╯           ╰────╯        
 -12288 ┤                                                                                           
 -16384 ┤                                                                                           
 -20479 ┤                                                                                           
 -24575 ┤                                                                                           
 -28671 ┤                                                                                           
 -32767 ┤                                                                                           
           C4 f = 261.626 kHz, Sustain = 0.30. Scale 0.10x

  32767 ┤                                                                                           
  28671 ┤                                                                                           
  24575 ┤                                                                                           
  20479 ┤                                                                                           
  16384 ┤                                                                                           
  12288 ┤                                                                                           
   8192 ┤ ╭───╮        ╭───╮         ╭──╮         ╭───╮        ╭───╮         ╭──╮         ╭──╮      
   4096 ┤╭╯   ╰╮      ╭╯   ╰╮       ╭╯  ╰╮       ╭╯   ╰╮      ╭╯   ╰╮      ╭─╯  ╰╮       ╭╯  ╰╮     
      0 ┼╯     ╰╮    ╭╯     ╰╮     ╭╯    ╰╮     ╭╯     ╰╮    ╭╯     ╰╮     │     ╰╮     ╭╯    ╰╮    
  -4096 ┤       ╰╮  ╭╯       ╰╮   ╭╯      ╰╮   ╭╯       ╰╮  ╭╯       ╰╮  ╭─╯      ╰╮   ╭╯      ╰╮   
  -8192 ┤        ╰──╯         ╰───╯        ╰───╯         ╰──╯         ╰──╯         ╰───╯        ╰── 
 -12288 ┤                                                                                           
 -16384 ┤                                                                                           
 -20479 ┤                                                                                           
 -24575 ┤                                                                                           
 -28671 ┤                                                                                           
 -32767 ┤                                                                                           
           E4 f = 329.628 kHz, Sustain = 0.30. Scale 0.10x

  32767 ┤                                                                                           
  28671 ┤                                                                                           
  24575 ┤                                                                                           
  20479 ┤                                                                                           
  16384 ┤                                                                                           
  12288 ┤                                                                                           
   8192 ┤ ╭──╮       ╭──╮       ╭──╮       ╭──╮        ╭──╮       ╭──╮       ╭──╮       ╭──╮        
   4096 ┤╭╯  ╰╮     ╭╯  ╰╮     ╭╯  ╰╮     ╭╯  ╰╮      ╭╯  ╰╮     ╭╯  ╰╮     ╭╯  ╰╮     ╭╯  ╰╮     ╭ 
      0 ┼╯    │    ╭╯    ╰╮    │    ╰╮    │    ╰╮    ╭╯    │    ╭╯    │    ╭╯    ╰╮    │    ╰╮    │ 
  -4096 ┤     ╰╮  ╭╯      ╰╮  ╭╯     ╰╮  ╭╯     ╰╮  ╭╯     ╰╮  ╭╯     ╰╮  ╭╯      ╰╮  ╭╯     ╰╮  ╭╯ 
  -8192 ┤      ╰──╯        ╰──╯       ╰──╯       ╰──╯       ╰──╯       ╰──╯        ╰──╯       ╰──╯  
 -12288 ┤                                                                                           
 -16384 ┤                                                                                           
 -20479 ┤                                                                                           
 -24575 ┤                                                                                           
 -28671 ┤                                                                                           
 -32767 ┤                                                                                           
           G4 f = 391.995 kHz, Sustain = 0.30. Scale 0.10x

  32767 ┤                                                                                           
  28671 ┤                                                                                           
  24575 ┤  ╭─╮                                                                ╭─╮                   
  20479 ┤ ╭╯ │                                                   ╭╮          ╭╯ │                   
  16384 ┤ │  ╰╮                                                 ╭╯╰╮         │  ╰╮                  
  12288 ┤╭╯   │                                                ╭╯  ╰╮       ╭╯   │         ╭╮       
   8192 ┤│    ╰╮       ╭──╮                 ╭──╮               │    │       │    ╰╮       ╭╯╰╮      
   4096 ┤│     │      ╭╯  ╰─╮   ╭──╮       ╭╯  ╰─╮            ╭╯    ╰╮      │     │      ╭╯  ╰╮     
      0 ┼╯     │      │     ╰───╯  ╰╮     ╭╯     ╰─────╮      │      │     ╭╯     │     ╭╯    ╰╮  ╭ 
  -4096 ┤      ╰╮    ╭╯             ╰─╮  ╭╯            ╰─╮   ╭╯      │     │      ╰╮    │      ╰──╯ 
  -8192 ┤       │    │                ╰──╯               ╰───╯       ╰╮    │       │   ╭╯           
 -12288 ┤       ╰╮  ╭╯                                                │   ╭╯       ╰╮  │            
 -16384 ┤        │ ╭╯                                                 │   │         │ ╭╯            
 -20479 ┤        ╰─╯                                                  ╰╮ ╭╯         ╰─╯             
 -24575 ┤                                                              ╰─╯                          
 -28671 ┤                                                                                           
 -32767 ┤                                                                                           
           C major chord. Scale 0.10x

C major, square wave
--------------------
  32767 ┤                                                                                           
  28671 ┤                                                                                           
  24575 ┤                                                                                           
  20479 ┤                                                                                           
  16384 ┤                                                                                           
  12288 ┤                                                                                           
   8192 ┼────────╮       ╭────────╮       ╭────────╮       ╭───────╮        ╭───────╮       ╭────── 
   4096 ┤        │       │        │       │        │       │       │        │       │       │       
      0 ┼        │       │        │       │        │       │       │        │       │       │       
  -4096 ┤        │       │        │       │        │       │       │        │       │       │       
  -8192 ┤        ╰───────╯        ╰───────╯        ╰───────╯       ╰────────╯       ╰───────╯       
 -12288 ┤                                                                                           
 -16384 ┤                                                                                           
 -20479 ┤                                                                                           
 -24575 ┤                                                                                           
 -28671 ┤                                                                                           
 -32767 ┤                                                                                           
           C4 f = 261.626 kHz, Sustain = 0.30. Scale 0.10x

  32767 ┤                                                                                           
  28671 ┤                                                                                           
  24575 ┤                                                                                           
  20479 ┤                                                                                           
  16384 ┤                                                                                           
  12288 ┤                                                                                           
   8192 ┼──────╮      ╭─────╮      ╭──────╮     ╭──────╮      ╭─────╮      ╭──────╮     ╭──────╮    
   4096 ┤      │      │     │      │      │     │      │      │     │      │      │     │      │    
      0 ┼      │      │     │      │      │     │      │      │     │      │      │     │      │    
  -4096 ┤      │      │     │      │      │     │      │      │     │      │      │     │      │    
  -8192 ┤      ╰──────╯     ╰──────╯      ╰─────╯      ╰──────╯     ╰──────╯      ╰─────╯      ╰─── 
 -12288 ┤                                                                                           
 -16384 ┤                                                                                           
 -20479 ┤                                                                                           
 -24575 ┤                                                                                           
 -28671 ┤                                                                                           
 -32767 ┤                                                                                           
           E4 f = 329.628 kHz, Sustain = 0.30. Scale 0.10x

  32767 ┤                                                                                           
  28671 ┤                                                                                           
  24575 ┤                                                                                           
  20479 ┤                                                                                           
  16384 ┤                                                                                           
  12288 ┤                                                                                           
   8192 ┼─────╮     ╭────╮     ╭─────╮    ╭─────╮    ╭─────╮    ╭─────╮     ╭────╮     ╭─────╮    ╭ 
   4096 ┤     │     │    │     │     │    │     │    │     │    │     │     │    │     │     │    │ 
      0 ┼     │     │    │     │     │    │     │    │     │    │     │     │    │     │     │    │ 
  -4096 ┤     │     │    │     │     │    │     │    │     │    │     │     │    │     │     │    │ 
  -8192 ┤     ╰─────╯    ╰─────╯     ╰────╯     ╰────╯     ╰────╯     ╰─────╯    ╰─────╯     ╰────╯ 
 -12288 ┤                                                                                           
 -16384 ┤                                                                                           
 -20479 ┤                                                                                           
 -24575 ┤                                                                                           
 -28671 ┤                                                                                           
 -32767 ┤                                                                                           
           G4 f = 391.995 kHz, Sustain = 0.30. Scale 0.10x

  32767 ┤                                                                                           
  28671 ┼─────╮                                                 ╭──╮        ╭────╮          ╭╮      
  24575 ┤     │                                                 │  │        │    │          ││      
  20479 ┤     │                                                 │  │        │    │          ││      
  16384 ┤     │                                                 │  │        │    │          ││      
  12288 ┤     │                                                 │  │        │    │          ││      
   8192 ┤     ╰╮      ╭─────╮  ╭──╮╭─╮    ╭────────╮ ╭─╮      ╭─╯  ╰╮       │    ╰╮     ╭───╯╰─╮  ╭ 
   4096 ┤      │      │     │  │  ││ │    │        │ │ │      │     │       │     │     │      │  │ 
      0 ┼      │      │     │  │  ││ │    │        │ │ │      │     │       │     │     │      │  │ 
  -4096 ┤      │      │     │  │  ││ │    │        │ │ │      │     │       │     │     │      │  │ 
  -8192 ┤      ╰─╮  ╭─╯     ╰──╯  ╰╯ ╰────╯        ╰─╯ ╰──────╯     ╰─╮    ╭╯     ╰─╮  ╭╯      ╰──╯ 
 -12288 ┤        │  │                                                 │    │        │  │            
 -16384 ┤        │  │                                                 │    │        │  │            
 -20479 ┤        │  │                                                 │    │        │  │            
 -24575 ┤        │  │                                                 │    │        │  │            
 -28671 ┤        ╰──╯                                                 ╰────╯        ╰──╯            
 -32767 ┤                                                                                           
           C major chord. Scale 0.10x
```

You can find the complete example in [cmd/synth/main.go](./cmd/synth/main.go). It's not a symphony, but hey, it's a chord, so I can say we are official graduating from noise to music!

#### [cmd/synth/main.go](./cmd/synth/main.go)

```shell
$ go run cmd/synth/main.go
```

## 2 Noteworthy Theory

Now we have a way to create sounds with a certain **frequency**. To go into proper musical territory, let's start by creating **notes**.

The notes names are **A, B, C, D, E, F, G** in English speaking countries. Or **Do, Re, Mi, Fa, Sol, La, Si** in many other countries.

Let's start with the basics. According to [wikipedia](https://en.wikipedia.org/wiki/Musical_note):

> In music, a **note** is the **pitch** and **duration** of a sound, and also its representation in musical notation (♪, ♩). A note can also represent a pitch class. Notes are the building blocks of much written music: discretizations of musical phenomena that facilitate performance, comprehension, and analysis.

### 2.1 Do, a Deer, a Female Deer

#### 2.1.1 The Theory

The **pitch** of a note is, in practical terms, the same as the **frequency**, and is also measured in **hertz (Hz)**.

The audible range of frequencies is around 20 Hz to 20 kHz. If we plot the frequency in a horizontal line, and pick any point, we could call it a **note**. In relation to this first note, any point in a lower pitch would be called a **lower** or **flatter** note, and any point with a higher pitch can be called a **higher** or **sharper** note.

```
lower pitch                                                         higher pitch
──────────┬────────────────────────────────────┬─────────────────┬──────────────
         note                                 note           also a note

*linear scale
```

If we pick any note pitch **X**, and play it along another note with double its pitch (**2X**), they will feel very similar. In a way they feel like they are the same note, and these two points receive the same name. The distance between a note and another one with double its pitch is called an **octave**.

The octaves below and above a note X can be visualized in a line like this:

```
  octave   octave
      ◄--►◄------►◄--- octave ---►◄----------- octave -----------►
...───┬───┬───────┬───────────────┬───────────────────────────────┬──────────...
    1/8X  1/4X   1/2X             X                               2X            

*linear scale
```

Now let's focus on any octave. Let's divide this interval into some notes. For historical reasons an octave gets divided into 12 notes, with equal distance between them. But, this "equal distance" is only equal on a logarithmic scale. If you take a look at the previous line, the width of each octave grows exponentially, and the same happens with the distance between notes within an octave.

Since a note X and the next X have a ratio of 2:1 (octave), each one of the 12 notes has a ratio with the next one of 2<sup>1/12</sup>:1, or **<sup>12</sup>√2:1**.

```
   ◄------------------------------- octave -------------------------------►
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─...
   1     2     3     4     5     6     7     8     9     10    11    12    1    

*logarithmic scale
```

The distance between any of these points is called a **semitone** or **half step**, and two semitones equal a whole **tone** or **step**. Using the labels above, 4 is one semitone above 3, and one tone above 2.

But how do we match these 1 to 12 semitone labels to the note names, when we only have 7 names? That's because the note names are related to the **major scale**. Without going into details yet, a **scale** is a pattern to select some of the 12 semitones of an octave. For the major scale, the pattern is:

```
whole, whole, half, whole, whole, whole, half
```

Which means that from the starting point we pick the notes (1 to 12) one step higher, two steps higher, 2 and a half steps higher, and so on.

```
   ◄------------------------------- octave -------------------------------►
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─...
   │           │           │     │           │           │           │     │    
   *           *           *     *           *           *           *     *    
       whole       whole    half     whole       whole       whole    half

*logarithmic scale
```

If we apply this pattern starting on the note **C** we get the **C major scale**. The notes of this scale are the ones that get assigned names:

```
   ◄------------------------------- octave -------------------------------►
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─...
   │           │           │     │           │           │           │     │    
   C           D           E     F           G           A           B     C
   Do          Re          Mi    Fa          Sol         La          Si    Do

*logarithmic scale
```

As we can see between E and F there is one semitone, and between D and E a tone. The note one semitone above D is called **D♯** (D **sharp**) or **E♭** (E **flat**). Both names refer to the exact same note.

These names use **accidentals**. An accidental is a symbol that modifies the pitch of a note:

- **♯** (or #): **sharp**, means the note is one semitone higher in pitch.
- **♭** (or b): **flat**, means the note is one semitone lower in pitch.
- **♮**: **natural**, means the note uses its normal note pitch. Used to cancel previous accidentals only, doesn't need to be specified all the time.

```
   ◄------------------------------- octave -------------------------------►
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─...
   │           │           │     │           │           │           │     │    
   C   C#/Db   D   D#/Eb   E     F   F#/Gb   G   G#/Ab   A   A#/Bb   B     C
   Do  Do#/Reb Re  Re#/Mib Mi    Fa Fa#/Solb Sol Sol#/Lab La La#/Sib Si    Do

*logarithmic scale
```

The previous graph with the C major scale corresponds to the piano keys. The white keys are the notes of the scale, and get natural names, and the black keys of the piano are the rest of the 12 notes, the sharps (or flats) of the white keys.

```
    C#  D#      F#  G#  A#      C#  D#
    Db  Eb      Gb  Ab  Bb      Db  Eb
║░░███░███░░║░░███░███░███░░║░░███░███░░║
║░░███░███░░║░░███░███░███░░║░░███░███░░║
║░░███░███░░║░░███░███░███░░║░░███░███░░║
║░░███░███░░║░░███░███░███░░║░░███░███░░║
║░░░║░░░║░░░║░░░║░░░║░░░║░░░║░░░║░░░║░░░║
║░░░║░░░║░░░║░░░║░░░║░░░║░░░║░░░║░░░║░░░║
║░░░║░░░║░░░║░░░║░░░║░░░║░░░║░░░║░░░║░░░║
╚═══╩═══╩═══╩═══╩═══╩═══╩═══╩═══╩═══╩═══╝
  C   D   E   F   G   A   B   C   D   E
```

Now let's see how we can calculate the exact frequency Hz values for each note. First, we need a better name for the note we want to know the exact pitch of.

As we saw if we have a note X of a made up pitch of 100 Hz, the frequencies of 200 Hz, 400 Hz will also get the same note name. To differentiate we use the scientific pitch notation, which combines the note name with its octave, starting with C0.

```
-- octave 2 --►◄------------------ octave 3 ------------------►◄--- octave 4 ---
...┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬─...
   A2  A#2 B2  C3  C#3 D3  D#3 E3  F3  F#3 G3  G#3 A3  A#3 B3  C4  C#4 D4  D#4

*logarithmic scale
```

In this tuning system the reference point from where all the other notes are calculated is **A4**, with a standardized value of **440 Hz**.

From A4 we can obtain the values of the other notes by applying the ratio of **<sup>12</sup>√2:1** for each semitone of distance.

```
n = number of semitones between the note and A4
pitch = 440 * 2^(n/12)
```

For example, B4 is 2 semitones above A4:

```
n = 2
B4 = 440 * 2^(2/12) = 493.88 Hz
```

C4 is 9 semitones below A4:

```
n = -9
C4 = 440 * 2^(-9/12) = 261.62 Hz
```

But why this ratio?

The same way a note C4 and the next C5 have a ratio of 2:1 (octave), there are other ratios that produce notes that sound good together. For example 3:2 (perfect fifth), 4:3 (perfect fourth), 5:4 (major third).

A tuning system with these exact ratios, named **just intonation**, comes with a few practical difficulties. An alternative is an equal temperament system, in which each adjacent note interval has the same ratio. For historical reasons the most common tuning system in western music today is the **12-tone equal temperament**, or 12-TET.

In 12-TET an octave gets divided into 12 notes, with equal (logarithmic) distance between them. While 12-TET creates intervals that are not _exactly_ those nice whole number ratios, in practice the approximations are so close that the difference is hard to notice.

For example the perfect fifth of C is G. In just intonation G would be an interval of 3/2 = 1.5 times the pitch of C.
With 12-TET the interval is 2<sup>​7/12</sup> = 1.498307

| | 12-TET Ratio | 12-TET decimal | Just intonation ratio | Just intonation decimal |
| --- | --- | --- | --- | --- |
| Perfect fifth | 2<sup>​7/12</sup> | 1.498307 | 3/2 | 1.5 |

So equal temperament allows to have notes close enough to the ideal ratios, while solving many practical problems of just intonation tuning.

#### 2.1.2 The Implementation

The easiest way to code all the pitches would be to have a list of all the names with their frequency in Hz. Something like:

```go
const (
	// ...
	C4      float64 = 261.63
	Csharp4 float64 = 277.18
	D4      float64 = 293.66
	Dsharp4 float64 = 311.13
	// ...
)
```

But I'd like to to some operations over the note pitches, things like `C4 + Octave == C5`. So I will assign a sequential int value to each pitch.

This sequence could start at any point, for example it would make sense to make C0 the first element:

```
C0  = 0
C#0 = 1
...
```

Or follow the nomenclature of an 88 key piano, and start with:

```
A0  = 0
A#0 = 1
...
```

But there is already an standard for this: MIDI (Musical Instrument Digital Interface). In MIDI 1.0 instrument keys are defined with 8-bit numbers from 0 to 127. The convention is that A4 is the MIDI key 69. This means the first key is actually below the octave 0, it starts on C-1 and ends in G9:

```
C-1  = 0
C#-1 = 1
D-1  = 2
...
G#4  = 68
A4   = 69
A#4  = 70
...
F9   = 125
F#9  = 126
G9   = 127
```

Taking this into account the code in Go could look like this:

```go
package note

// Note is a musical pitch and duration
type Note struct {
	Pitch
	// Duration
}

// Interval is the distance between pitches, measured in semitones
type Interval = uint8

const (
	// Semitone is the smallest interval between pitches, a 12th of an octave
	Semitone Interval = 1
	// Tone is 2 semitones
	Tone Interval = 2
	// Octave is the distance between a pitch and another with double frequency
	Octave Interval = 12
)

// Pitch represents a musical frequency
type Pitch uint8

func (p Pitch) String() string {
	return pitchValues[p].name
}

func (p Pitch) Frequency() float64 {
	return pitchValues[p].frequency
}

func (p Pitch) Add(i Interval) Pitch {
	return Pitch(uint8(p) + uint8(i))
}

func (p Pitch) Subtract(i Interval) Pitch {
	return Pitch(uint8(p) - uint8(i))
}

const (
	...
	G4       Pitch = 67
	Gsharp4  Pitch = 68
	A4       Pitch = 69
	...
)


var pitchValues = map[Pitch]struct {
	name      string
	frequency float64
}{
	...
	G4:       {"G4", 391.99},
	Gsharp4:  {"G#4", 415.30},
	A4:       {"A4", 440},
	...
}
```

But copying all the values can be tedious. Why would you spend 2 minutes on a tedious task if you can spend 2 hours automating it? Obviously this needs some code generation to make it more interesting.

Using `go generate` is easy. First, let's put the previous type definitions without any of the pitch values in [./theory/note/note.go](./theory/note/note.go).

Then we need to add the `go:generate` keyword:

```go
// Package note contains types to manage musical notes
package note

//go:generate go run ../../gen/main.go
```

Now we need a Go program that will generate all the pitches automatically and place them in a new [./theory/note/frequencies.go](./theory/note/frequencies.go) file.

This text template will become the `frequencies.go` file:

```go
var tpl = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
package note

const (
{{- range .Notes}}
	{{ printf "%-8v Pitch = %v" .VarName .KeyNumber }}
{{- end}}
)

var pitchValues = map[Pitch]struct {
	name      string
	frequency float64
}{
{{- range .Notes}}
	{{ printf "%-9v {%q, %v}," (printf "%v:" .VarName) .Name .Freq }}
{{- end}}
}
`))
```

To fill this template we will need to provide `Notes`, an array of this struct type:

```go
	type note struct {
		VarName   string
		Name      string
		KeyNumber int
		Freq      float64
	}
```

Now comes the interesting part. We need 128 of those `notes`. As we have seen in the previous section, the frequency of any note can be calculated from the reference frequency of A4 = 440 Hz.

```go
	for i := 0; i <= 127; i++ {
		// MIDI key 69 is used for A4, 440 Hz in standard tuning
		// If n = number of semitones between the note and A4
		// then pitch = 440 * 2^(n/12)
		distance := float64(i) - 69
		freq := 440 * math.Pow(2, (distance/12))

		// ...
	}
```

As for the names, they can be built from the list of note names + the octave number. This will require some cleaning to have nice printable names and valid variable names.

```go
// sanitizeVarName replaces forbidden characters for variable names, e.g.
// A#4 => Asharp4; C-1 => C_1
func sanitizeVarName(name string) string {
	return strings.Replace(
		strings.Replace(name, "#", "sharp", 1),
		"-1", "_1", 1)
}

func main() {
	var names = []string{
		"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

	for i := 0; i <= 127; i++ {
		octave := int(i/12) - 1
		name := names[i%12] + strconv.Itoa(octave)

		// ...

		notes = append(notes, note{
			VarName:   sanitizeVarName(name),
			Name:      name,
			KeyNumber: i,
			Freq:      freq,
		})
	}
}
```

You can see all parts put together in

##### [gen/main.go](./gen/main.go)

And after running `go generate`

```shell
$ go generate ./theory/note/
```

The resulting file gets generated in

##### [theory/note/frequencies.go](./theory/note/frequencies.go)

This `Pitch` type can be used in a sample command that takes [cmd/synth/main.go](./cmd/synth/main.go) and extends it to build chords using this new `theory/note` package.

```go
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

func play(wave synth.WaveGenerator, triad []note.Pitch) io.Reader {
	return synth.Combine(
		synth.Sustain(wave(sampleRate, triad[0].Frequency(), 1*time.Second), 0.2),
		// ...
	)
}

func main() {
	// ...

	pitch := note.C3

	for i := uint8(0); i < 3; i++ {
		triad := majorChord(pitch.Add(i * note.Octave))
		plotChord(synth.NewSineWave, triad)

		sound := play(synth.NewSineWave, triad)

		// ...
	}
}
```

You can find the complete example in [cmd/pitch/main.go](./cmd/pitch/main.go). When you run it, it plays the C major chord in 3 different octaves.

##### [cmd/pitch/main.go](./cmd/pitch/main.go)

```shell
$ go run cmd/pitch/main.go
```

```
  32767 ┤                                                                                           
  28671 ┤      ╭╮                                                                                   
  24575 ┤    ╭─╯╰╮                                                                                  
  20479 ┤   ╭╯   ╰─╮                                                                                
  16384 ┤  ╭╯      │                                                                                
  12288 ┤ ╭╯       ╰╮                                                                               
   8192 ┤╭╯         ╰╮               ╭──────╮                                   ╭─────╮             
   4096 ┤│           ╰╮             ╭╯      ╰──╮        ╭─────╮               ╭─╯     ╰──╮          
      0 ┼╯            │            ╭╯          ╰────────╯     ╰──╮          ╭─╯          ╰────╮ ╭── 
  -4096 ┤             ╰╮          ╭╯                             ╰─╮      ╭─╯                 ╰─╯   
  -8192 ┤              ╰╮        ╭╯                                ╰──────╯                         
 -12288 ┤               ╰╮     ╭─╯                                                                  
 -16384 ┤                ╰╮   ╭╯                                                                    
 -20479 ┤                 ╰───╯                                                                     
 -24575 ┤                                                                                           
 -28671 ┤                                                                                           
 -32767 ┤                                                                                           
           C3,E3,G3. Scale 0.10x

  32767 ┤                                                                                           
  28671 ┤                                                                                           
  24575 ┤  ╭─╮                                                                ╭─╮                   
  20479 ┤ ╭╯ │                                                   ╭╮          ╭╯ │                   
  16384 ┤ │  ╰╮                                                 ╭╯╰╮         │  ╰╮                  
  12288 ┤╭╯   │                                                ╭╯  ╰╮       ╭╯   │         ╭╮       
   8192 ┤│    ╰╮       ╭──╮                 ╭──╮               │    │       │    ╰╮       ╭╯╰╮      
   4096 ┤│     │      ╭╯  ╰─╮   ╭──╮       ╭╯  ╰─╮            ╭╯    ╰╮      │     │      ╭╯  ╰╮     
      0 ┼╯     │      │     ╰───╯  ╰╮     ╭╯     ╰─────╮      │      │     ╭╯     │     ╭╯    ╰╮  ╭ 
  -4096 ┤      ╰╮    ╭╯             ╰─╮  ╭╯            ╰─╮   ╭╯      │     │      ╰╮    │      ╰──╯ 
  -8192 ┤       │    │                ╰──╯               ╰───╯       ╰╮    │       │   ╭╯           
 -12288 ┤       ╰╮  ╭╯                                                │   ╭╯       ╰╮  │            
 -16384 ┤        │ ╭╯                                                 │   │         │ ╭╯            
 -20479 ┤        ╰─╯                                                  ╰╮ ╭╯         ╰─╯             
 -24575 ┤                                                              ╰─╯                          
 -28671 ┤                                                                                           
 -32767 ┤                                                                                           
           C4,E4,G4. Scale 0.10x

  32767 ┤                                                                                           
  28671 ┤                                                                                           
  24575 ┤ ╭╮                               ╭╮                                                       
  20479 ┤╭╯│                               ││                               ╭─╮                     
  16384 ┤│ │                        ╭─╮    │╰╮                        ╭╮    │ │               ╭╮    
  12288 ┤│ │                        │ │   ╭╯ │   ╭╮         ╭─╮      ╭╯│    │ │         ╭╮    │╰╮   
   8192 ┤│ ╰╮   ╭╮        ╭─╮       │ │   │  │   │╰╮  ╭─╮   │ │      │ ╰╮   │ │   ╭─╮  ╭╯│   ╭╯ │   
   4096 ┤│  │  ╭╯╰╮ ╭─╮   │ ╰╮     ╭╯ │   │  │   │ │  │ │   │ ╰╮    ╭╯  │  ╭╯ ╰╮  │ │  │ ╰╮  │  │   
      0 ┼╯  │  │  ╰─╯ ╰╮ ╭╯  ╰╮    │  ╰╮  │  │  ╭╯ │ ╭╯ │  ╭╯  │    │   │  │   │  │ │  │  │  │  ╰╮  
  -4096 ┤   │  │       │ │    ╰──╮╭╯   │  │  │  │  ╰─╯  ╰╮ │   ╰╮  ╭╯   │  │   │ ╭╯ ╰╮╭╯  │  │   │  
  -8192 ┤   │ ╭╯       ╰─╯       ╰╯    │ ╭╯  ╰╮ │        │ │    ╰──╯    │  │   │ │   ╰╯   │ ╭╯   ╰╮ 
 -12288 ┤   ╰╮│                        │ │    │╭╯        ╰─╯            ╰╮ │   │ │        ╰╮│     ╰ 
 -16384 ┤    ││                        │ │    ││                         │╭╯   ╰─╯         ││       
 -20479 ┤    ╰╯                        ╰╮│    ╰╯                         ╰╯                ╰╯       
 -24575 ┤                               ╰╯                                                          
 -28671 ┤                                                                                           
 -32767 ┤                                                                                           
           C5,E5,G5. Scale 0.10x
```

### 2.2 Pata-Pata-Pata-Pon

Now that the `Note` contains a `Pitch`, it's time to add the missing half: **duration**.

The **duration** defines how long a note is to be played. But this duration is not measured in seconds, instead the note symbols are relative to the **whole note**.

| Note | Name                        | Relative value |
| ---- | --------------------------- | -------------- |
| 𝅜    | Double note / breve         | 2              |
| 𝅝    | Whole note / semibreve      | 1              |
| 𝅗𝅥    | Half note / minim           | 1/2            |
| ♩    | Quarter note / crotchet     | 1/4            |
| ♪    | Eighth note / quaver        | 1/8            |
| 𝅘𝅥𝅯    | Sixteenth note / semiquaver | 1/16           |

So the duration in seconds of any note can be calculated once we know the value for a whole note. In order to do that, we need to talk about the **staff** and **tempo**

A **staff** (or **stave**) is the group of five lines where the note symbols are placed. The piece of music written in the **staff** gets divided into uniform sections called **bars** (or **measures**).

A **time signature** (or **meter signature**) is placed the start of a **staff** for a piece of music, consisting of 2 numbers, one on top of each other. The top number indicates how many **beats** there are in each **bar**, and the lower number indicates the note equivalent to each **beat**.

For example:
- (⁴₄) means that in each bar there are 4 beats, and each beat corresponds to a quarter note (1/4)
- (³₈) means that in each bar there are 3 beats, and each beat corresponds to an eighth note (1/8)

A **staff** will also contain a **tempo**, measured in **bpm (beats per minute)**. This defines the speed of the piece of music.

The **bpm** gives us how many beats there are in a minute, and the lower number of the **time signature** defines what note is equivalent to a beat. Knowing this we can calculate the length of a note in seconds.

For example, if we have a (⁴₄) signature and 120 bpm, and we want to know the length of a sixteenth note:

![score](./doc/tempo.png)

- 60 seconds per minute / 120 beats per minute = 0.5 **seconds per beat**
- Per the time signature, a **quarter note** (1/4) is a beat, so 1/4 of a note is 0.5 seconds
- 0.5 seconds per beat / (1/4) note per beat = 2 **seconds per whole note**
- 2 seconds per whole note * (1/16) of a note = 0.125 **seconds per sixteenth note**

```go
package note

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
```

```go
func TestToSeconds(t *testing.T) {
	n := note.Note{Duration: note.Quarter}

	assert.Equal(t, 500*time.Millisecond, n.ToSeconds(note.Quarter, 120))
	assert.Equal(t, time.Second, n.ToSeconds(note.Eighth, 120))
	assert.Equal(t, 250*time.Millisecond, n.ToSeconds(note.Half, 120))
}
```

You can find the complete code in [./theory/note/duration.go](./theory/note/duration.go), and the tests to verify that everything works as it should in [./theory/note/duration_test.go](./theory/note/duration_test.go)

In `note.go`, the code relative to intervals and pitch gets moved to [./theory/note/pitch.go](./theory/note/pitch.go). And now `Note` can be a `Pitch` plus a `Duration`:

```go
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
```

And can be used like this:

```go
func TestNoteString(t *testing.T) {
	n := note.NewNote(note.Fsharp5, note.Quarter)

	assert.Equal(t, "♩ F#5", fmt.Sprint(n))
	assert.Equal(t, "♩ F#5", n.String())
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
```

You can find the complete code in [./theory/note/note.go](./theory/note/note.go), and the tests for it in [./theory/note/note_test.go](./theory/note/note_test.go)

### 2.3 The Sound... of Silence

There is something missing in the implementation of musical notes so far. There is no way to define a **rest**.

A **rest** is a symbol that defines an interval of silence, with a length corresponding to a note name.

| Note | Name | Relative value |
| ---- | ---- | -------------- |
| 𝄺 | Double rest / breve rest         | 2    |
| 𝄻 | Whole rest / semibreve rest      | 1    |
| 𝄼 | Half rest / minim rest           | 1/2  |
| 𝄽 | Quarter rest / crotchet rest     | 1/4  |
| 𝄾 | Eighth rest / quaver rest        | 1/8  |
| 𝄿 | Sixteenth rest / semiquaver rest | 1/16 |

Adding support this this in `duration.go` is simple, we only have to account for the different string representation, since the length in seconds is calculated the same way.

```go
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
	...
}

// String returns the musical note symbol
func (d Duration) String() string {
	...
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
```

But how to represent a generic thing that can be a note or a rest? I opted to rename the previous `note.Pitch` type to `note.pitchValue`, create a new `note.restPitch`, and extract the methods to a `note.Pitch` interface. So a `Note` will contain a `Pitch`, that will be initialized to either a **pitch** or a **rest**.

```go
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

// NewNote creates a new musical Note
func NewNote(p Pitch, d Duration) Note {
	return Note{
		Pitch:    p,
		Duration: d,
	}
}

// NewRest creates a new rest Note
func NewRest(d Duration) Note {
	return Note{
		Pitch:    restPitch{},
		Duration: d,
	}
}

// Pitch is an interface for a note musical frequency
type Pitch interface {
	// String returns a human readable name for this pitch
	String() string
	// Frequency returns the hertz value for this pitch
	Frequency() float64
	// Add returns a new Pitch adding an interval to this pitch, making it higher
	Add(i Interval) Pitch
	// Subtract returns a new Pitch subtracting an interval to this pitch, making it lower
	Subtract(i Interval) Pitch
}

// restPitch represents the (absence of) musical frequency for a rest note
type restPitch struct{}

// Frequency returns the 0 hertz value
func (r restPitch) Frequency() float64 {
	return 0
}

// Add is a no operation, returns a RestPitch
func (r restPitch) Add(i Interval) Pitch {
	return r
}

// Subtract is a no operation, returns a RestPitch
func (r restPitch) Subtract(i Interval) Pitch {
	return r
}

// String is a no operation, returns an empty string
func (r restPitch) String() string {
	return ""
}

// pitchValue represents a musical frequency
type pitchValue uint8

// String returns a human readable name for this pitch
func (p pitchValue) String() string {
	return pitchValues[p].name
}

// Frequency returns the hertz value for this pitch
func (p pitchValue) Frequency() float64 {
	return pitchValues[p].frequency
}

// Add returns a new Pitch adding an interval to this pitch, making it higher
func (p pitchValue) Add(i Interval) Pitch {
	return pitchValue(uint8(p) + uint8(i))
}

// Subtract returns a new Pitch subtracting an interval to this pitch, making it lower
func (p pitchValue) Subtract(i Interval) Pitch {
	return pitchValue(uint8(p) - uint8(i))
}
```

### 2.4 Lost My Marbles

Enough theory. The repository still does not have chords, scales, ADSR... but all in good time.

For now with notes and rests we can play **M U S I C**. Since this whole project is a convoluted way to make music, let's play a tribute to another beautifully convoluted music instrument: [Wintergatan's Marble Machine](https://www.youtube.com/watch?v=IvUU8joBb1Q).

A staff is a list of notes played. Sometimes several notes can be played at the same time. So a staff can be represented as a `[][]note.Note`. In the piano score (https://musescore.com/user/5631216/scores/1846226) there are two staffs, the treble and bass:

```go
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

	...

}

var bassStaff = [][]note.Note{
	// Bar 1
	[]note.Note{note.NewRest(note.Whole)},

	...

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

	...

}
```

A staff can produce the `io.Reader` to be be played combining the sound for groups of simultaneous notes `[]note.Note`, and then concatenating them one after the other:

```go
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
```

And then both staffs can be combined with `synth.Combine` too:

```go
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
```

You can find the complete example in [cmd/marble/main.go](./cmd/marble/main.go).

##### [cmd/marble/main.go](./cmd/marble/main.go)

```shell
$ go run cmd/marble/main.go
```

## 3 Harmony

### 3.1 Scaling Up

A **scale** is an ordered group of notes, in ascending or descending order.
It is usually described by the distance (**interval**) between two consecutive notes.

This group of notes is used to guide the composition of melodies and harmonies.

The easiest scale to build is the **chromatic scale**, which defined as intervals between each note looks like:

```
half, half, half, half, half, half, half, half, half, half, half, half
```

Visualized in a line, the chromatic scale is just the list of all the 12 notes in an octave. If we apply the interval patter to the **C** note, we get this:

```
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─...
   │     │     │     │     │     │     │     │     │     │     │     │     │    
   |half |half |half |half |half |half |half |half |half |half |half |half |    
   │     │     │     │     │     │     │     │     │     │     │     │     │    
   *     *     *     *     *     *     *     *     *     *     *     *     *    
   C   C#/Db   D   D#/Eb   E     F   F#/Gb   G   G#/Ab   A   A#/Bb   B     C
```

The other scale already mentioned in a previous section is the **major scale**, which has this sequence of intervals:

```
whole, whole, half, whole, whole, whole, half
```

```
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─...
   │           │           │     │           │           │           │     │    
   *           *           *     *           *           *           *     *    
       whole       whole    half     whole       whole       whole    half
```

If we apply this interval to the note C, we get the **C major scale**. C is the starting point, a **whole** step above it is D, another **whole** step above is E, **half** a step above is F, and so on.

```
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─...
   C   C#/Db   D   D#/Eb   E     F   F#/Gb   G   G#/Ab   A   A#/Bb   B     C
   │           │           │     │           │           │           │     │    
   |   whole   |   whole   |half |   whole   |   whole   |   whole   |half |    
   │           │           │     │           │           │           │     │
   *           *           *     *           *           *           *     *    
   C           D           E     F           G           A           B     C
```

This pattern can be applied to any starting note, not only natural notes. It can be applied to any of the 12 notes, for example, E. This is the **E major scale**:

```
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─...
   E     F   F#/Gb   G   G#/Ab   A   A#/Bb   B     C   C#/Db   D   D#/Eb   E
   │           │           │     │           │           │           │     │    
   |   whole   |   whole   |half |   whole   |   whole   |   whole   |half |    
   │           │           │     │           │           │           │     │
   *           *           *     *           *           *           *     *
   E           F#          G#    A           B           C#          D#    E
```

The **natural minor scale** follows this other pattern of step intervals:

```
whole, half, whole, whole, half, whole, whole
```

```
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─...
   │           │     │           │           │     │           │           │    
   *           *     *           *           *     *           *           *    
       whole    half     whole       whole    half     whole       whole
```

Applying this pattern to A we get the **A natural minor scale**:

```
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─...
   A   A#/Bb   B     C   C#/Db   D   D#/Eb   E     F   F#/Gb   G   G#/Ab   A
   │           │     │           │           │     │           │           │    
   |   whole   |half |   whole   |   whole   |half |   whole   |   whole   |
   │           │     │           │           │     │           │           │
   *           *     *           *           *     *           *           *    
   A           B     C           D           E     F           G           A
```

The reason why a scale is called **major** or **minor** has to do with **intervals**. 

**Step**, **half step**, and **octave** are names for intervals between 2 notes. There are other interval names based on the **major scale**. Because a scale can be moved and applied to any note, the notes of a scale can also be called by their **scale degrees**, their relative position in the **major scale**. Arabic numerals are used for this.

This relative position in the **major scale** gives name to commonly used **intervals**. Note that the last note of the scale is the same as the first one, but one octave above. The same number, 1, is used. The pattern then repeats one octave higher, with the same numbering (1, 2, 3...).

```
   ◄------------------------------- octave -------------------------------►◄----- octave --...
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬───...
   │           │           │     │           │           │           │     │           │
   |   whole   |   whole   |half |   whole   |   whole   |   whole   |half |   whole   |
   │           │           │     │           │           │           │     │           │
   *           *           *     *           *           *           *     *           *
   1           2           3     4           5           6           7     1           2
   ◄► Unison (same note)
   ◄----------► Second (whole step)
   ◄----------------------► Third (2 steps)
   ◄----------------------------► Fourth (2 and a half steps)
   ◄--- Fifth (3 and a half steps) ---------►
   ◄--- Sixth (4 and a half steps) ---------------------►
   ◄--- Seventh (5 and a half steps) -------------------------------►
   ◄--- Octave (12 half steps) -------------------------------------------►
   ◄--- Ninth (7 steps) --------------------------------------------------------------►
```

For the notes that are not in the major scale, we use something similar to the note accidentals. Instead of natural, sharp, and flat, we use **perfect**, **augmented**, and **diminished**. A perfect interval is implied when its name is not augmented or diminished.

```
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬───...
   │           │           │     │           │           │           │     │           │
   *           *           *     *           *           *           *     *           *
   1           2           3     4           5           6           7     1           2
   ◄--- Diminished Sixth -------------------------►
   ◄--- Perfect Sixth ----------------------------------►
   ◄--- Augmented Sixth --------------------------------------►
```

A scale is called **minor** when it contains a diminished third, and **major** when it contains a perfect third. Commonly these 2 intervals are just called **minor third** and **major third**.

Now let's take another look at the major and minor scale patterns:

```
major:         whole, whole, half,  whole, whole, whole, half
natural minor: whole, half,  whole, whole, half,  whole, whole
```

Let's see what happens when we compare the G major scale to the E minor scale:
```
G major scale
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─...
   G   G#/Ab   A   A#/Bb   B     C   C#/Db   D   D#/Eb   E     F   F#/Gb   G
   │           │           │     │           │           │           │     │    
   |   whole   |   whole   |half |   whole   |   whole   |   whole   |half |    
   │           │           │     │           │           │           │     │
   *           *           *     *           *           *           *     *
   G           A           B     C           D           E           F#    G

E minor scale
...┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┬─...
   E     F   F#/Gb   G   G#/Ab   A   A#/Bb   B     C   C#/Db   D   D#/Eb   E
   │           │     │           │           │     │           │           │    
   |   whole   |half |   whole   |   whole   |half |   whole   |   whole   |
   │           │     │           │           │     │           │           │
   *           *     *           *           *     *           *           *    
   E           F#    G           A           B     C           D           E
```

They contain the same pool of notes, but reordered. This is because the patterns are actually the same pattern, but shifted:

```
major:         whole, whole, half,  whole, whole, whole, half
natural minor:                                    whole, half,  whole, whole, half,  whole, whole
```

All the patterns that can be built shifting the major scale pattern are valid scales, but they are usually called **modes**. All of the modes have their own name, this is the complete list:

```
Ionian     whole, whole, half,  whole, whole, whole, half
Dorian            whole, half,  whole, whole, whole, half,  whole
Phrygian                 half,  whole, whole, whole, half,  whole, whole
Lydian                          whole, whole, whole, half,  whole, whole, half
Mixolydian                             whole, whole, half,  whole, whole, half,  whole
Aeolian                                       whole, half,  whole, whole, half,  whole, whole
Locrian                                              half,  whole, whole, half,  whole, whole, whole
```

The Ionian mode is commonly known as _the major scale_, and the Aeolian mode as _the natural minor scale_.

The **major** and **minor** scales are probably the two most common ones. But by definition any pattern, with any number of notes, can be considered a scale. Of course not all of them will sound good.

The most commonly used scales are, in no particular order:

| | | |
| --- | --- | --- |
| Major scale       |                                   | whole, whole, half, whole, whole, whole, half |
| Major scale modes | Ionian (aka Major scale)          | whole, whole, half, whole, whole, whole, half |
|                   | Dorian                            | whole, half, whole, whole, whole, half, whole |
|                   | Phrygian                          | half, whole, whole, whole, half, whole, whole |
|                   | Lydian                            | whole, whole, whole, half, whole, whole, half |
|                   | Mixolydian                        | whole, whole, half, whole, whole, half, whole |
|                   | Aeolian (aka Natural Minor scale) | whole, half, whole, whole, half, whole, whole |
|                   | Locrian                           | half, whole, whole, half, whole, whole, whole |
| Minor scales      | Natural Minor       | whole, half, whole, whole, half, whole, whole |
|                   | Harmonic Minor      | whole, half, whole, whole, half, whole and a half, half |
|                   | Melodic Minor       | whole, half, whole, whole, whole, whole, half |
| Pentatonic scales | Minor Pentatonic    | whole and a half, whole, whole, whole and a half, whole |
|                   | Major Pentatonic    | whole, whole, whole and a half, whole, whole and a half |
| Blues scale       |                     | whole and a half, whole, half, half, whole and a half, whole |
| Whole-tone scale  |                     | whole, whole, whole, whole, whole, whole |

The implementation of scales in code could look like this:

```go
func Major(tonic note.Pitch) []note.Pitch {
	var pattern = []note.Interval{
		note.Tone, note.Tone, note.Semitone, note.Tone, note.Tone, note.Tone, note.Semitone}

	scale := []note.Pitch{tonic}
	for _, interval := range pattern {
		scale = append(scale, scale[len(scale)-1].Add(interval))
	}

	return scale
}
```

Just take a pattern, and build a slice adding intervals to the last element. Refactoring this to be more generic, the final code looks like this:

```go
func scale(tonic note.Pitch, pattern []note.Interval) []note.Pitch {
	scale := []note.Pitch{tonic}
	for _, interval := range pattern {
		scale = append(scale, scale[len(scale)-1].Add(interval))
	}

	return scale
}

var major = []note.Interval{
	note.Tone, note.Tone, note.Semitone, note.Tone, note.Tone, note.Tone, note.Semitone}

func Major(tonic note.Pitch) []note.Pitch {
	return scale(
		tonic,
		major,
	)
}
```

For the modes of the major scale we can shift the pattern we already have, like this:

```go

// shift a pattern n positions to the left
func shift(pattern []note.Interval, n int) []note.Interval {
	return append(pattern[n:], pattern[0:n]...)
}

func Ionian(tonic note.Pitch) []note.Pitch {
	return Major(tonic)
}

func Dorian(tonic note.Pitch) []note.Pitch {
	return scale(tonic, shift(major, 1))
}

func Phrygian(tonic note.Pitch) []note.Pitch {
	return scale(tonic, shift(major, 2))
}

...
```