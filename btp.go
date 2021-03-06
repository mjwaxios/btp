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
		r.data = make([]byte, 0, 1500)
		r.collecting = true
		return nil, nil
	}

	// Ignore the byte if we are not collecting
	if !r.collecting {
		return nil, nil
	}

	// Check if we are esced
	if r.esced {
		switch b {
		case ESTX, EETX, EESC:
			// Add none escapped byte to data
			r.data = append(r.data, (b & 0x7F))
			r.esced = false
		default:
			r.collecting = false
			return nil, errors.New("Bad Escapped Value")
		}
		return nil, nil
	}

	// Check for everything else
	switch b {
	case ETX:
		r.collecting = false
		return r.data, nil
	case ESC:
		r.esced = true
	default:
		r.data = append(r.data, b)
	}

	return nil, nil
}

// Encode takes a byte slice and encodes it into BTP format
func Encode(bin []byte) ([]byte, error) {
	buff := make([]byte, 0, 1500)
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
