// BTP provides framing to binary streaming data
package btp

import "errors"

const (
	STX  = 0x02
	ETX  = 0x03
	ESC  = 0x1B
	ESTX = 0x82
	EETX = 0x83
	EESC = 0x9B
)

type Receiver struct {
	collecting bool
	esced      bool
	data       []byte
}

// Decode takes a byte and decodes and adds it to the buffer,  if it marks the end of a packet
//   we return the decoded packet
func (r *Receiver) Decode(b byte) ([]byte, error) {
	// If start of Packet
	if b == STX {
		r.data = make([]byte, 1500)
		r.collecting = true
		return nil, nil
	}

	// Ignore the byte if we are not collecting
	if !r.collecting {
		return nil, nil
	}

	switch b {
	case ETX:
		r.collecting = false
		return r.data, nil
	case ESC:
		r.esced = true
	case ESTX, EETX, EESC:
		if r.esced {
			// Add none escapped byte to data
			r.data = append(r.data, b&0x7F)
		}
	default:
		if r.esced {
			return nil, errors.New("Bad Escapped Value")
		}
		r.data = append(r.data, b)
	}

	return nil, nil
}

// Encode takes a byte slice and encodes it into BTP format
func Encode(bin []byte) ([]byte, error) {
	var buff []byte

	buff = append(buff, STX)

	for _, b := range bin {
		switch b {
		case STX, ETX, ESC:
			buff = append(buff, ESC)
			buff = append(buff, (0x80 | b))
		default:
			buff = append(buff, b)
		}
	}

	buff = append(buff, ETX)

	return buff, nil
}
