# Music Playground

This repository is a kind of notebook where I'll experiment programing music theory concepts from scratch, using Go.

I am deliberately avoiding to use any code as reference for the implementation of a synthesizer, or musical notation. There might be better ways to implement music concepts, but the goal is for me to figure out how to achieve it.

## Table of Contents

- [1 Let There Be Noise](#1-let-there-be-noise)
  - [1.1 PCM](#11-pcm)
  - [1.2 Beep Boop](#12-beep-boop)
  - [1.3 Pure Waves](#13-pure-waves)
  - [1.4 World's Smallest Synthesizer](#14-worlds-smallest-synthesizer)

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

In previous examples we used an `io.MultiReader` to concatenate `io.Readers`, playing one sound wave after the other. So we could say the code was a **monophonic** synthesizer, able to play one note at a time. Making **polyphonic** music much more fun, so let's also have some code to mix different samples.

The PCM library used, [Oto](https://github.com/hajimehoshi/oto), already supports sending samples simultaneously. So we could be playing different sounds from a few goroutines. But having having my own function to mix the samples coming from different `io.Readers` will allow me to be more independent of the PCM library. This could be used in the future for example to save the final mixed samples into a sound file.

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
