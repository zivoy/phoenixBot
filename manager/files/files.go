package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/gommon/bytes"
)

const (
	Cache  = "./.cache"
	Images = "./images"

	purgeUsageUnder = 5
	purgeUsageOver  = 10 * time.Minute
	purgeMaxSize    = 100 * bytes.MiB
)

var (
	stopPurge = make(chan interface{})
	images    = &imageList{Images: make(map[string]*Image)}
)

func init() {
	_ = os.Mkdir(Cache, os.ModePerm)
	_ = os.Mkdir(Images, os.ModePerm)

	images.LoadImagesFromDisc()

	go func() {
		t := time.NewTicker(3 * time.Minute)
		for {
			select {
			case <-t.C:
				go images.PurgeList()
			case <-stopPurge:
				t.Stop()
				return
			}
		}
	}()
}

type imageList struct {
	Images map[string]*Image
	mu     sync.RWMutex
}

type Image struct {
	Id       string    `json:"id"`
	LastCall time.Time `json:"last_requested"`

	mu          sync.Mutex
	DefaultSize *ImageSize `json:"size"`
	sizes       map[string]*ImageSize

	Deleted bool `json:"-"`
}

type ImageSize struct {
	lastCall    time.Time
	timesCalled int

	Size   int64 `json:"filesize"`
	Width  int   `json:"width"`
	Height int   `json:"height"`

	Name string `json:"-"`
	Path string `json:"-"`
}

func StopCleaning() {
	stopPurge <- true
}

func (l *imageList) PurgeList() {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for _, im := range l.Images {
		im.mu.Lock()
		var size int64
		now := time.Now()
		for id, s := range im.sizes {
			if s.lastCall.Add(purgeUsageOver).Before(now) && s.timesCalled < purgeUsageUnder {
				go im.removeVariant(id)
				continue
			}
			size += s.Size
		}
		if size > purgeMaxSize {
			go im.purgeVariants()
		}
		im.mu.Unlock()
	}
}

func (l *imageList) LoadImagesFromDisc() {
	l.mu.Lock()
	defer l.mu.Unlock()
	files, err := os.ReadDir(Images)
	if err != nil {
		log.Fatal("cant load images", err)
	}

	for _, f := range files {
		name := f.Name()
		if name[len(name)-3:] != "dat" {
			continue
		}

		f, err := os.Open(fmt.Sprintf("%s/%s", Images, name))
		if err != nil {
			log.Println("error opening ", name, err)
			continue
		}

		dec := gob.NewDecoder(f)
		i := new(Image)
		err = dec.Decode(i)
		if err != nil {
			log.Println("error decoding file ", name, err)
			_ = f.Close()
			continue
		}
		i.sizes = make(map[string]*ImageSize)
		l.Images[i.Id] = i

		_ = f.Close()
	}
}

func GetImages() []*Image {
	images.mu.RLock()
	defer images.mu.RUnlock()
	l := make([]*Image, 0, len(images.Images))
	for _, v := range images.Images {
		if !v.Deleted {
			l = append(l, v)
		}
	}
	return l
}

func GetImage(id string, width, height int, imType string) *ImageSize {
	images.mu.RLock()
	defer images.mu.RUnlock()
	image, ok := images.Images[id]
	if !ok || image.Deleted {
		return nil
	}

	if width > image.DefaultSize.Width {
		width = image.DefaultSize.Width
	}
	if height > image.DefaultSize.Height {
		height = image.DefaultSize.Height
	}

	if width == 0 && height == 0 { // todo keep track of type as well
		image.inc(image.DefaultSize)
		return image.DefaultSize
	}

	// variation exists
	if size, ok := image.sizes[getImageSize(imType, width, height)]; ok {
		image.inc(size)
		return size
	}

	// size doesn't exist, resize
	return image.Resize(imType, width, height)
}

func (i *Image) inc(s *ImageSize) {
	i.mu.Lock()
	defer i.mu.Unlock()
	s.lastCall = time.Now()
	i.LastCall = s.lastCall
	s.timesCalled++
}

func getImageSize(imType string, width, height int) string {
	return fmt.Sprintf("%s-%dx%d", imType, width, height)
}

func AddImage(image io.Reader, name string) (*Image, error) {
	images.mu.Lock()
	defer images.mu.Unlock()
	_ = name
	defaultImage, _, err := decodeImage(image)
	if err != nil {
		return nil, err
	}

	// todo resize to make fit in max size to save space

	im := &Image{
		Id:       strings.ToUpper(strings.Replace(uuid.New().String(), "-", "", -1)),
		LastCall: time.Now(),
		sizes:    make(map[string]*ImageSize),
	}
	im.DefaultSize = im.saveImage(defaultImage, "webp", Images)
	if err = im.saveInfo(); err != nil {
		return nil, err
	}
	images.Images[im.Id] = im
	log.Println("uploaded image ", im.Id)
	return im, nil
}

func (i *Image) saveInfo() error { // todo save to db
	f, err := os.Create(fmt.Sprintf("%s/%s.dat", Images, i.Id))
	if err != nil {
		log.Println("can't create info file ", err)
		return err
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Println("error closing file ", f.Name())
		}
	}(f)
	enc := gob.NewEncoder(f)
	err = enc.Encode(i)
	if err != nil {
		log.Println("can't write info file ", err)
		return err
	}
	return nil
}

func RemoveImage(imageid string) error {
	images.mu.Lock()
	defer images.mu.Unlock()
	image, ok := images.Images[imageid]
	if !ok {
		return errors.New("image not found")
	}

	go image.purgeVariants()
	image.mu.Lock()
	image.Deleted = true
	image.mu.Unlock()

	log.Println("disabled image ", imageid)
	return nil
}

func (i *Image) removeVariant(size string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	im, ok := i.sizes[size]
	if !ok {
		return
	}
	delete(i.sizes, size)
	if err := os.Remove(im.Path); err != nil {
		log.Println("could not delete ", im.Name)
	}
}

func (i *Image) purgeVariants() {
	i.mu.Lock()
	defer i.mu.Unlock()
	for _, s := range i.sizes {
		err := os.Remove(s.Path)
		if err != nil {
			log.Println("cannot remove ", s.Path, err)
		}
	}
	i.sizes = make(map[string]*ImageSize)
}
