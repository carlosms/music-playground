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

`oto.Player` implements `io.WriteCloser`. To play sound we just need to use the `Write` method to send the samples.

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

```shell
$ go run cmd/noise/main.go
```
