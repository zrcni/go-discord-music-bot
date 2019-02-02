package audiorepository

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/zrcni/go-discord-music-bot/videoaudio"
)

const packageName = "AudioRepository"

var packagePrefix = fmt.Sprintf("[%s]", packageName)

var audioRepository = &repository{
	items: make(map[string]Item),
}

type AudioType int

const (
	RAW AudioType = iota
	FILE
)

type Item interface {
	GetID() string
}

type Raw struct {
	ID  string
	Raw []byte
}

// GetID returns the item ID
func (r Raw) GetID() string {
	return r.ID
}

type File struct {
	ID   string
	Name string
}

// GetID returns the item ID
func (f File) GetID() string {
	return f.ID
}

func (f File) GetRawData() ([]byte, error) {
	path, err := videoaudio.GetTrackFilepath(f.Name)
	if err != nil {
		return nil, errors.Wrap(err, packageName)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, packageName)
	}

	return data, nil
}

type repository struct {
	items map[string]Item
}

// Add adds an item to the audiorepository, return false if item already exists
func Add(item Item) (ok bool) {
	itemID := item.GetID()
	if audioRepository.items[itemID] != nil {
		return false
	}

	audioRepository.items[itemID] = item
	return true
}

// Get gets audio data and removes it, return false if item was not found
func Get(ID string) ([]byte, error) {
	log.Debugf("%s getting item %s", packagePrefix, ID)

	if audioRepository.items[ID] == nil {
		return nil, errors.Wrap(fmt.Errorf("audio by id %s doesn't exist in audio repository", ID), packageName)
	}

	switch item := audioRepository.items[ID].(type) {
	case *Raw:
		defer deleteItem(item)
		return item.Raw, nil
	case *File:
		defer deleteItem(item)
		return item.GetRawData()
	}

	// This code should never be reached
	panic(errors.Wrap(fmt.Errorf("audio repository item is not of valid type: %v", audioRepository.items[ID]), packageName))
}

func deleteItem(item Item) error {
	itemID := item.GetID()

	log.Debugf("%s deleting item %s", packagePrefix, itemID)

	audioRepository.items[itemID] = nil

	f, ok := item.(*File)
	if !ok {
		return errors.Wrap(fmt.Errorf("%s could not delete item %s", packagePrefix, itemID), packageName)
	}

	path, err := videoaudio.GetTrackFilepath(f.Name)
	if err != nil {
		return errors.Wrap(err, packageName)
	}

	err = os.Remove(path)
	if err != nil {
		return errors.Wrap(err, packageName)
	}

	return nil
}
