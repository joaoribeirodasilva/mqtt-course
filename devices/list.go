package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Sensors struct {
	door        SensorDoor
	temperature SensorNumeric
	humidity    SensorNumeric
}

type DataList struct {
	conf          *Configuration
	mu            sync.Mutex
	list          []*Sensors
	isStarted     bool
	stopRequested bool
	isDirty       bool
}

// NewDataList creates a new thread safe and auto save list pointer
func NewDataList(conf *Configuration) *DataList {

	dl := &DataList{}
	dl.conf = conf
	dl.list = make([]*Sensors, 1)
	dl.isStarted = false
	dl.stopRequested = false
	dl.isDirty = false

	return dl
}

// Start starts the list autosave functionality
func (dl *DataList) Start() error {

	// if this thread is already started
	// then don't restart it again
	if dl.isStarted {

		return nil
	}

	// tries to read the list data from the file
	if err := dl.Read(); err != nil {
		fmt.Errorf("WARNING: failed to read data file REASON: %s", err.Error())
	}

	// starts it's own thread
	go func() {

		// flags it's started
		dl.isStarted = true

		// while not requested to stop then loop
		for !dl.stopRequested {

			if err := dl.Save(); err != nil {
				fmt.Errorf("WARNING: failed save data file REASON: %s", err.Error())
			}

			// sleep until next time to save
			time.Sleep(time.Duration(dl.conf.Data.SaveInterval) * time.Millisecond)

		}

		// here the stop request flag was set
		// flag it as not started
		dl.isStarted = false

		// reset the stop request flag
		dl.stopRequested = false

		// save the current DataList array data to the disk
		if err := dl.Save(); err != nil {
			fmt.Errorf("ERROR: failed save data file REASON: %s", err.Error())
		}

		// sets the is dirty flag to false
		dl.isDirty = false
	}()

	return nil
}

// Stop requests the autosave functionality to stop
func (dl *DataList) Stop() {

	// if it's started
	if dl.isStarted {

		// create a channel to wait for the thread to terminate
		done := make(chan bool)

		// sets the stop request to true
		dl.stopRequested = true

		done <- !dl.isStarted

		<-done
	}
}

// Append ands a new Sensors struct to the DataList array
func (dl *DataList) Append(item Sensors) {

	dl.mu.Lock()

	if uint32(len(dl.list)) >= dl.conf.Data.MaxMessages {
		log.Printf("WARNING: sensor list reached it's limit of %d messages stored", dl.conf.Data.MaxMessages)
		dl.list = dl.list[1:]
	}

	dl.list = append(dl.list, &item)
	dl.isDirty = true

	dl.mu.Unlock()
}

// Removes n items from the head of the DataList array
func (dl *DataList) Remove(items int) {

	dl.mu.Lock()

	if items > len(dl.list) {
		items = len(dl.list)
	}

	dl.list = dl.list[items:]
	dl.isDirty = true

	dl.mu.Unlock()
}

func (dl *DataList) GetHead() *Sensors {

	dl.mu.Lock()
	defer dl.mu.Unlock()
	if len(dl.list) == 0 {
		return nil
	}

	s := &Sensors{
		door:        dl.list[0].door,
		temperature: dl.list[0].temperature,
		humidity:    dl.list[0].humidity,
	}

	if s.door.isOpen {
		s.door.closeTime = time.Unix(0, 0)
	}

	return s
}

func (dl *DataList) Len() int {

	dl.mu.Lock()
	defer dl.mu.Unlock()

	return len(dl.list)
}

// Save saves the DataList array into a JSON file
func (dl *DataList) Save() error {

	if !dl.isDirty {
		return nil
	}

	dl.mu.Lock()
	defer dl.mu.Unlock()

	data, err := json.Marshal(dl.list)
	if err != nil {
		return err
	}

	err = os.WriteFile(dl.conf.Data.Path, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Reads reads the DataList array from a JSON file
func (dl *DataList) Read() error {

	dl.mu.Lock()
	defer dl.mu.Unlock()

	bytes, err := os.ReadFile(dl.conf.Data.Path)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(bytes), dl.list)
	if err != nil {
		return err
	}

	dl.isDirty = false

	return nil
}

func (dl *DataList) IsDirty() bool {

	dl.mu.Lock()
	defer dl.mu.Unlock()
	return dl.isDirty
}
