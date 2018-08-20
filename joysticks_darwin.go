// +build darwin

package joysticks

import "net"

// see if Device exists.
func DeviceExists(index uint8) bool {
	return true
}

// Connect sets up a go routine that puts a joysticks events onto registered channels.
// to register channels use the returned HID object's On<xxx>(index) methods.
// Note: only one event, of each type '<xxx>', for each 'index', so re-registering, or deleting, an event stops events going on the old channel.
// It Needs the HID objects ParcelOutEvents() method to be running to perform routing.(so usually in a go routine.)
func Connect(index int) (d *HID) {
	r, _ := net.Dial("tcp", "rpi-cnc.local:5001")
	d = &HID{make(chan osEventRecord), make(map[uint8]button), make(map[uint8]hatAxis), make(map[eventSignature]chan Event)}
	// start thread to read joystick events to the joystick.state osEvent channel
	go eventPipe(r, d.OSEvents)
	d.populate()
	return d
}
