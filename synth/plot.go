package synth

import (
	"fmt"
	"io"
	"math"

	"github.com/carlosms/asciigraph"
)

// width is the number of data points to print in the X axis.
// Magic number that makes the plot fit in the GitHub README.md width
const width = 92

// Plot reads the int16 samples from the io.Reader and returns a XY plot.
// The number of points in the X axis is fixed, use scale to control zooming
// out, e.g. for a scale of 0.25 each data point represents 4 samples
func Plot(r io.Reader, caption string, scale float64) string {
	samplesPoint := int(math.Round(1 / scale))

	// Each data point is made of 2 bytes per sample, times the number of
	// samples for each point (scale)
	buf := make([]byte, width*2*samplesPoint)
	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		panic(err)
	}

	// number of bytes for each data point
	bytesPoint := 2 * samplesPoint

	var plotData []float64
	for i := 0; i < n-bytesPoint; i += bytesPoint {
		// Do the average of the number of samples per data point (samplesPoint)
		var v float64
		for j := 0; j < samplesPoint; j++ {
			v += float64(int16(buf[i]) + int16(buf[i+1])<<8)
		}
		v = v / float64(samplesPoint)

		// little-endian
		plotData = append(plotData, v)
	}

	return asciigraph.Plot(plotData,
		asciigraph.Caption(fmt.Sprintf("%s. Scale %.2fx", caption, scale)),
		asciigraph.Height(15), asciigraph.Min(float64(-max)), asciigraph.Max(float64(max)))
}
