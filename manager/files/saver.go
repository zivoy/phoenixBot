package files

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

var ValidImage = regexp.MustCompile(`^(webp|png|jpe?g)$`)

func (i *Image) Resize(imType string, width, height int) *ImageSize {
	r := i.DefaultSize
	i.inc(r)

	data, err := os.ReadFile(r.Path)
	if err != nil {
		log.Println("error reading default image size ", err)
		return nil
	}

	im, t, err := decodeImage(bytes.NewReader(data))
	if err != nil {
		log.Println("can't decode image of type ", t)
		return nil
	}

	var res *image.NRGBA
	if width == 0 || height == 0 {
		res = imaging.Resize(im, width, height, imaging.Lanczos /*imaging.CatmullRom*/)
	} else {
		res = imaging.Fit(im, width, height, imaging.Lanczos /*imaging.CatmullRom*/)

	}
	s := i.saveImage(res, imType, Cache)
	size := getImageSize(imType, width, height)
	log.Printf("Creating %s at %s -- %s\n", i.Id, size, s.Name)
	i.sizes[size] = s
	i.inc(s)
	return s
}

func nameImage(image []byte) string {
	h := sha256.New()
	h.Write(image)
	return fmt.Sprintf("%X", h.Sum(nil))
}

func (i *Image) saveImage(image image.Image, imType, location string) *ImageSize {
	i.mu.Lock()
	defer i.mu.Unlock()

	var buf bytes.Buffer
	if err := encodeImage(&buf, image, imType); err != nil {
		log.Println("error encoding image", err)
		return nil
	}
	b := buf.Bytes()
	name := nameImage(b)

	if imType == "jpg" {
		imType = "jpeg"
	}

	path := fmt.Sprintf("%s/%s-%s.%s", location, i.Id, name, imType)
	if err := os.WriteFile(path, b, 0664); err != nil {
		log.Println("error writing image do disc", err)
		return nil
	}

	box := image.Bounds()
	im := &ImageSize{
		Size:   int64(len(b)),
		Width:  box.Dx(),
		Height: box.Dy(),
		Name:   name,
		Path:   path,
	}

	return im
}

func decodeImage(reader io.Reader) (image.Image, string, error) {
	return image.Decode(reader)
}

func encodeImage(buf *bytes.Buffer, image image.Image, imType string) error {
	switch imType {
	case "webp":
		return webp.Encode(buf, image, &webp.Options{Lossless: true})
	case "png":
		return png.Encode(buf, image)
	case "jpg", "jpeg":
		return jpeg.Encode(buf, image, &jpeg.Options{Quality: 80})
	default:
		return errors.New("invalid image type")
	}
}
