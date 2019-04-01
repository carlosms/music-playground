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

One byte holds the value of one sample. This value can be in the range [0..255], so we can fill them with random data.

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
Because the bit depth is set to 8 bits, the value of each sample can range from 0 to 255. The middle point of the wave, or **equilibrium**, will be 127. This means the amplitude must have a value in the range [0..127].

The function will simply divide the period _T_ in 2, and set `equilibrium + amplitude` or `equilibrium - amplitude`, depending on which half of the period the sample index falls on.

```go
const (
	sampleRate        = 44100
	bitDepthInBytes   = 1
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
```

This `io.Reader` can be consumed by the `oto.Player` like this:

```go
p, _ := oto.NewPlayer(sampleRate, channelNum, bitDepthInBytes, bufferSizeInBytes)

io.Copy(p, sqrWave(600, 20, time.Second))
```

We can also print an ASCII graph of the sampled values using [github.com/guptarohit/asciigraph](https://github.com/guptarohit/asciigraph). Or, rather, [my fork](https://github.com/carlosms/asciigraph) that adds a couple of new options.

```go
func plot(freq, amplitude int, duration time.Duration) string {
	r := sqrWave(freq, amplitude, duration)
	data, _ := ioutil.ReadAll(r)

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
}
```

```
 255 ┼                                            
 238 ┤                                            
 221 ┤                                            
 204 ┤                                            
 187 ┤                                            
 170 ┼──────────╮          ╭──────────╮           
 153 ┤          │          │          │           
 136 ┤          │          │          │           
 119 ┤          │          │          │           
 102 ┤          │          │          │           
  85 ┤          ╰──────────╯          ╰────────── 
  68 ┤                                            
  51 ┤                                            
  34 ┤                                            
  17 ┤                                            
   0 ┤                                            
        f = 2 kHz, A = 50

 255 ┼                                            
 238 ┤                                            
 221 ┼──╮  ╭──╮  ╭──╮  ╭──╮  ╭──╮  ╭──╮  ╭──╮  ╭─ 
 204 ┤  │  │  │  │  │  │  │  │  │  │  │  │  │  │  
 187 ┤  │  │  │  │  │  │  │  │  │  │  │  │  │  │  
 170 ┤  │  │  │  │  │  │  │  │  │  │  │  │  │  │  
 153 ┤  │  │  │  │  │  │  │  │  │  │  │  │  │  │  
 136 ┤  │  │  │  │  │  │  │  │  │  │  │  │  │  │  
 119 ┤  │  │  │  │  │  │  │  │  │  │  │  │  │  │  
 102 ┤  │  │  │  │  │  │  │  │  │  │  │  │  │  │  
  85 ┤  │  │  │  │  │  │  │  │  │  │  │  │  │  │  
  68 ┤  │  │  │  │  │  │  │  │  │  │  │  │  │  │  
  51 ┤  │  │  │  │  │  │  │  │  │  │  │  │  │  │  
  34 ┤  ╰──╯  ╰──╯  ╰──╯  ╰──╯  ╰──╯  ╰──╯  ╰──╯  
  17 ┤                                            
   0 ┤                                            
        f = 6 kHz, A = 100
```

You can find the complete example in [cmd/wave/main.go](./cmd/wave/main.go).
This program plays a few random robot noises that could be sound effects from an old school Atari game.

#### [cmd/wave/main.go](./cmd/wave/main.go)

```shell
$ go run cmd/wave/main.go
```
