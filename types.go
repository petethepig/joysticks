package joysticks

import (
	"encoding/binary"
	"io"
	"time"
)

// see; https://www.kernel.org/doc/Documentation/input/joystick-api.txt
type osEventRecord struct {
	Time  uint32 // event timestamp, unknown base, in milliseconds 32bit, so about a month
	Value int16  // value
	Type  uint8  // event type
	Index uint8  // axis/button
}

// fill in the joysticks available events from the synthetic events burst produced initially by the driver.
func (d HID) populate() {
	for buttonNumber, hatNumber, axisNumber := 1, 1, 1; ; {
		evt := <-d.OSEvents
		switch evt.Type {
		case 0x81:
			d.Buttons[evt.Index] = button{uint8(buttonNumber), toDuration(evt.Time), evt.Value != 0}
			buttonNumber += 1
		case 0x82:
			d.HatAxes[evt.Index] = hatAxis{uint8(hatNumber), uint8(axisNumber), false, toDuration(evt.Time), float32(evt.Value) / maxValue}
			axisNumber += 1
			if axisNumber > 2 {
				axisNumber = 1
				hatNumber += 1
			}
		default:
			go func() { d.OSEvents <- evt }() // have to consume a real event to know we reached the end of the synthetic burst, so refire it.
			return
		}
	}
	return
}

// pipe any readable events onto channel.
func eventPipe(r io.Reader, c chan osEventRecord) {
	var evt osEventRecord
	for {
		if binary.Read(r, binary.LittleEndian, &evt) != nil {
			close(c)
			return
		}
		c <- evt
	}
}

func toDuration(m uint32) time.Duration {
	return time.Duration(m) * 1000000
}

const maxValue = 1<<15 - 1
