// +build linux

package joysticks

import (
	"os"
	"strconv"
)

// common path root, so Connect and DeviceExists are not thread safe.
var inputPathSlice = []byte("/dev/input/js ")[0:13]

// see if Device exists.
func DeviceExists(index uint8) bool {
	_, err := os.Stat(string(strconv.AppendUint(inputPathSlice, uint64(index-1), 10)))
	return err == nil
}

// Connect sets up a go routine that puts a joysticks events onto registered channels.
// to register channels use the returned HID object's On<xxx>(index) methods.
// Note: only one event, of each type '<xxx>', for each 'index', so re-registering, or deleting, an event stops events going on the old channel.
// It Needs the HID objects ParcelOutEvents() method to be running to perform routing.(so usually in a go routine.)
func Connect(index int) (d *HID) {
	r, e := os.OpenFile(string(strconv.AppendUint(inputPathSlice, uint64(index-1), 10)), os.O_RDWR, 0)
	if e != nil {
		return nil
	}
	d = &HID{make(chan osEventRecord), make(map[uint8]button), make(map[uint8]hatAxis), make(map[eventSignature]chan Event)}
	// start thread to read joystick events to the joystick.state osEvent channel
	go eventPipe(r, d.OSEvents)
	d.populate()
	return d
}
