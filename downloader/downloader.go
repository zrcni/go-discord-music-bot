package downloader

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/zrcni/go-discord-music-bot/queue"
)

const packageName string = "Downloader"

// Downloader handles downloading tracks
type Downloader struct {
	queue queue.Queue
}

// Downloadable has the
type Downloadable struct {
	ID  string
	Get func()
}

// New returns a new downloader
func New(len int) Downloader {
	return Downloader{
		queue: queue.New(len),
	}
}

// Queue for download
func (d *Downloader) Queue(item Downloadable) {
	err := d.queue.Add(item)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("[%s] added to queue: %s", packageName, item.ID)

	d.processQueue()
}

func (d *Downloader) processQueue() {
	log.Debugf("[%s] processing queue", packageName)

	if d.queue.Length() == 0 {
		log.Infof("[%s] queue is empty", packageName)
		return
	}

	item := d.queue.Shift()

	downloadable, ok := item.(Downloadable)
	if !ok {
		panic(fmt.Sprintf("[%s]item is not of type Downloadable: %+v", packageName, downloadable))
	}

	d.download(downloadable)
	log.Debugf("[%s] queue processed", packageName)

	d.processQueue()
}

func (d *Downloader) download(dl Downloadable) {
	log.Infof("[%s] downloading: %s", packageName, dl.ID)
	dl.Get()
}
