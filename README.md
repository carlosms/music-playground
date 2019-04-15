# Music Playground

This repository is a kind of notebook where I'll experiment programing music theory concepts from scratch, using Go.

## 1. Let there be noise

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

We can create one with this function:

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
		pos := (i / bitDepthInBytes) % samplesPeriod

		radian := float64(pos) / float64(samplesPeriod) * 2 * math.Pi
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