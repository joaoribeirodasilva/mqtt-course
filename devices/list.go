package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/joaoribeirodasilva/wait_signals"
)

type Sensors struct {
	Door        SensorDoor    `json:"door"`
	Temperature SensorNumeric `json:"temperature"`
	Humidity    SensorNumeric `json:"humidity"`
}

type DataList struct {
	conf          *Configuration
	mu            sync.Mutex
	list          []*Sensors
	isStarted     bool
	stopRequested bool
	isDirty       bool
	finished      chan bool
}

// NewDataList creates a new thread safe and auto save list pointer
func NewDataList(conf *Configuration) *DataList {

	dl := &DataList{}
	dl.conf = conf
	dl.list = make([]*Sensors, 0)
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

	// make/remake the finished channel
	dl.finished = make(chan bool, 1)

	// tries to read the list data from the file
	if err := dl.Read(); err != nil {
		fmt.Printf("WARNING: [LIST] failed to read data file REASON: %s\n", err.Error())
	}

	// starts it's own thread
	go func() {

		log.Println("INFO: [LIST] data list started")
		// flags it's started
		dl.isStarted = true

		// while not requested to stop then loop
		for !dl.stopRequested {

			if err := dl.Save(); err != nil {
				fmt.Printf("WARNING: [LIST] failed save data file REASON: %s\n", err.Error())
			}

			if sig := wait_signals.SleepWait(time.Duration(dl.conf.Data.SaveInterval)*time.Millisecond, syscall.SIGINT, syscall.SIGTERM); sig != nil {
				break
			}

			log.Printf("INFO: [LIST] buffer has %d messages stored\n", dl.Len())
		}

		log.Println("INFO: [LIST] data list stopping")
		// here the stop request flag was set
		// flag it as not started
		dl.isStarted = false

		// reset the stop request flag
		dl.stopRequested = false

		// save the current DataList array data to the disk
		if err := dl.Save(); err != nil {
			fmt.Printf("ERROR: [LIST] failed save data file REASON: %s", err.Error())
		}

		// sets the is dirty flag to false
		dl.isDirty = false

		// set the channel so the Stop function can stop waiting
		// for loop termination
		dl.finished <- true
	}()

	return nil
}

// Stop requests the autosave functionality to stop
func (dl *DataList) Stop() {

	// if it's started
	if dl.isStarted {

		log.Println("INFO: [LIST] data list stop requested... waiting")

		// create a channel to wait for the thread to terminate

		// sets the stop request to true
		dl.stopRequested = true

		// return when it's finished
		<-dl.finished

		log.Println("INFO: [LIST] data list stopped")
	}
}

// Append ands a new Sensors struct to the DataList array
func (dl *DataList) Append(item Sensors) {

	// lock the list so the thread inserting
	// can have exclusive access
	dl.mu.Lock()

	//if the list is full
	if uint32(len(dl.list)) >= dl.conf.Data.MaxMessages {

		// print a warning to the console
		//log.Printf("WARNING: sensor list reached it's limit of %d messages stored", dl.conf.Data.MaxMessages)

		// remove the oldest item from the list
		dl.list = dl.list[1:]
	}

	// add the item to the list
	dl.list = append(dl.list, &item)

	// set is dirty flag to true so
	// we know there are new items in
	// the list
	dl.isDirty = true

	// unlock the exclusive thread access
	// to the list
	dl.mu.Unlock()
}

// Removes n items from the head of the DataList array
func (dl *DataList) Remove(items int) {

	// lock the list so the thread inserting
	// can have exclusive access
	dl.mu.Lock()

	// if the number of items to remove is greater
	// than the items list
	if items > len(dl.list) {

		// set the items to remove equal to the
		// items list length
		items = len(dl.list)
	}

	// remove the oldest items from the list
	dl.list = dl.list[items:]

	// set the flag is dirty to true
	dl.isDirty = true

	// unlock the exclusive thread access
	// to the list
	dl.mu.Unlock()
}

// GetHead returns a copy of the first list item
func (dl *DataList) GetHead() *Sensors {

	// lock the list so the thread inserting
	// can have exclusive access
	dl.mu.Lock()

	// defer unlock the exclusive thread access
	// to the list
	defer dl.mu.Unlock()

	// if the list has no items return
	if len(dl.list) == 0 {
		return nil
	}

	// copy the sensor data item to a new
	// memory address without deep copy
	// that is processor intensive
	s := &Sensors{
		Door:        dl.list[0].Door,
		Temperature: dl.list[0].Temperature,
		Humidity:    dl.list[0].Humidity,
	}

	// if the door sensor is marked as open
	// we don't want to send the time the door
	// is set to close
	if s.Door.IsOpen {
		s.Door.CloseTime = nil
	}

	// return the copied list item
	return s
}

// Len returns the list size
func (dl *DataList) Len() int {

	// lock the list so the thread inserting
	// can have exclusive access
	dl.mu.Lock()

	// defer unlock the exclusive thread access
	// to the list
	defer dl.mu.Unlock()

	// return the length of the list
	return len(dl.list)
}

// Save saves the DataList array into a JSON file
func (dl *DataList) Save() error {

	// lock the list so the thread inserting
	// can have exclusive access
	dl.mu.Lock()

	// defer unlock the exclusive thread access
	// to the list
	defer dl.mu.Unlock()

	// if there are no changes to the list
	// return
	if !dl.isDirty {
		return nil
	}

	// transform the list data into a JSON array
	// of bytes
	data, err := json.Marshal(dl.list)
	if err != nil {
		return err
	}

	log.Printf("INFO: [LIST] saving data to file")

	// write this JSON byte array to a data file
	err = os.WriteFile(dl.conf.Data.Path, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Reads reads the DataList array from a JSON file
func (dl *DataList) Read() error {

	// lock the list so the thread inserting
	// can have exclusive access
	dl.mu.Lock()

	// defer unlock the exclusive thread access
	// to the list
	defer dl.mu.Unlock()

	log.Printf("INFO: [LIST] reading data from file")

	// read the JSON file into the a byte array
	bytes, err := os.ReadFile(dl.conf.Data.Path)
	if err != nil {
		return err
	}

	// transform the JSON bytes into a memory list
	err = json.Unmarshal([]byte(bytes), &dl.list)
	if err != nil {
		return err
	}

	// sets the flag is dirty equals to false
	dl.isDirty = false

	return nil
}

func (dl *DataList) IsDirty() bool {

	// lock the list so the thread inserting
	// can have exclusive access
	dl.mu.Lock()

	// defer unlock the exclusive thread access
	// to the list
	defer dl.mu.Unlock()

	// return the list is dirty flag status
	return dl.isDirty
}
