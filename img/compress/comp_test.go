package compress

import (
	"bytes"
	"github.com/fogleman/gg"
	"gotit/parallel"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

func TestCompress(t *testing.T) {
	jobs, _ := newJobs("/Users/zms/Downloads/workspace")
	parallel.Do(jobs, 1)
}

func TestCompressJPEG(t *testing.T) {
	f := "/Users/zms/Downloads/workspace/compress/2.jpg"
	img, _ := gg.LoadImage(f)
	dest, _ := os.Create(f + "-com.jpg")
	jpeg.Encode(dest, img, &jpeg.Options{
		Quality: 70,
	})
}

func TestCompressPNG(t *testing.T) {

	f := "/Users/zms/Downloads/workspace/compress/2.png"
	img, _ := gg.LoadImage(f)
	dest, _ := os.Create(f + "-com.png")

	tmp := &bytes.Buffer{}
	jpeg.Encode(tmp, img, &jpeg.Options{
		Quality: 60,
	})
	img2, _, _ := image.Decode(tmp)

	en := png.Encoder{CompressionLevel: png.BestCompression}
	en.Encode(dest, img2)

}
