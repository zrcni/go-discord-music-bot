package downloader

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/zrcni/go-discord-music-bot/chanqueue"
)

const packageName string = "Downloader"

// Downloader handles downloading tracks
type Downloader struct {
	queue chanqueue.Queue
}

// Downloadable has the
type Downloadable struct {
	ID  string
	Get func()
}

// // New returns a new downloader
// func New(len int) Downloader {
// 	downloader := Downloader{
// 		queue: listqueue.New(len),
// 	}
// 	return downloader
// }

// // Queue for download
// func (d *Downloader) Queue(item Downloadable) {
// 	log.Infof("----- QUEUEUE LENGTH: %v -----", d.queue.Length())
// 	err := d.queue.Add(item)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	}
// 	log.Infof("[%s] added to queue: %s", packageName, item.ID)

// 	d.processQueue()
// }

// func (d *Downloader) download(dl Downloadable) {
// 	log.Infof("[%s] downloading: %s", packageName, dl.ID)
// 	dl.Get()
// }

// func (d *Downloader) processQueue() error {
// 	item, err := d.queue.Process()
// 	if err != nil {
// 		return err
// 	}

// 	downloadable, ok := item.(Downloadable)
// 	if !ok {
// 		panic(fmt.Sprintf("[%s]item is not of type Downloadable: %+v", packageName, downloadable))
// 	}

// 	d.download(downloadable)

// 	return nil
// }

// CHANNEL QUEUE !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// New returns a new downloader
func New(len int) Downloader {
	queue := chanqueue.New(len)
	downloader := Downloader{}

	queue.Start()

	queue.ProcessItem(func(item interface{}) {
		downloadable, ok := item.(Downloadable)
		if !ok {
			panic(fmt.Sprintf("[%s]item is not of type Downloadable: %+v", packageName, downloadable))
		}

		downloader.download(downloadable)
	})

	downloader.queue = queue

	return downloader
}

// Queue for download
func (d *Downloader) Queue(item Downloadable) {
	err := d.queue.Add(item)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("[%s] added to queue: %s", packageName, item.ID)
}

func (d *Downloader) download(dl Downloadable) {
	log.Infof("[%s] downloading: %s", packageName, dl.ID)
	go dl.Get()
}
