package ansii

import (
	"bytes"
	"sync"
)

const (
	stateChar = iota + 1
	stateEscape
	stateCsiEscape
)

const (
	escapeChar      = 0x1b
	csiFinalByteMin = 0x40
	csiByteMax      = 0x7e
)

type CsiCode struct {
	Arg  string
	Code byte
}

func (c CsiCode) Bytes() []byte {
	return []byte("\033[" + c.Arg + string(c.Code))
}

type EscapeCodeHandler interface {
	HandleByte(b byte) []byte
	HandleCsi(csi CsiCode) []byte
}

type Ansii struct {
	mu sync.Mutex

	buf       *bytes.Buffer
	tmp       []byte
	state     int
	converter EscapeCodeHandler
}

func NewAnsii(c EscapeCodeHandler) *Ansii {
	return &Ansii{
		buf:       bytes.NewBuffer(make([]byte, 0)),
		tmp:       make([]byte, 0),
		state:     stateChar,
		converter: c,
	}
}

func (a *Ansii) Write(p []byte) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, b := range p {
		a.writeByte(b)
	}

	return len(p), nil
}

func (a *Ansii) WriteByte(c byte) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.writeByte(c)
	return nil
}

func (a *Ansii) Read(p []byte) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.buf.Read(p)
}

func (a *Ansii) ReadBytes(delim byte) ([]byte, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.buf.ReadBytes(delim)
}

func (a *Ansii) ReadByte() (byte, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.buf.ReadByte()
}

func (a *Ansii) ReadRune() (r rune, size int, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.buf.ReadRune()
}

func (a *Ansii) ReadString(delim byte) (line string, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.buf.ReadString(delim)
}

func (a *Ansii) ReadLine() (line string, err error) {
	return a.ReadString('\n')
}

func (a *Ansii) Bytes() []byte {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.buf.Bytes()
}

func (a *Ansii) String() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.buf.String()
}

func (a *Ansii) Len() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.buf.Len()
}

func (a *Ansii) writeByte(b byte) {
	var processed bool

	for {
		switch a.state {
		case stateChar:
			processed = a.processChar(b)
		case stateEscape:
			processed = a.processEscape(b)
		case stateCsiEscape:
			processed = a.processCsi(b)
		}

		if processed {
			break
		}
	}
}

func (a *Ansii) processChar(b byte) bool {
	// clean buffer
	if len(a.tmp) > 0 {
		a.tmp = make([]byte, 0)
	}

	if b == escapeChar {
		// escape char
		return a.change(stateEscape)
	}

	if a.handleByte(b) {
		return a.next(stateChar)
	}

	a.buf.WriteByte(b)
	return a.next(stateChar)
}


func (a *Ansii) processEscape(b byte) bool {
	a.tmp = append(a.tmp, b)

	// second byte is used to determine escape sequence type
	if 2 == len(a.tmp) {
		if b == '[' {
			// CSI
			return a.next(stateCsiEscape)
		} else if b == escapeChar{
			// double escape code are converted to single one
			a.buf.WriteByte(escapeChar)
			return a.next(stateChar)
		} else {
			// unknown or unsupported escape code
			return a.next(stateChar)
		}
	}

	return a.next(stateEscape)
}

func (a *Ansii) processCsi(b byte) bool {
	if b >= csiFinalByteMin && b <= csiByteMax {
		// if last byte was final byte
		a.handleCsi(CsiCode{string(a.tmp[2:len(a.tmp)]), b})

		return a.next(stateChar)
	}

	a.tmp = append(a.tmp, b)
	return a.next(stateCsiEscape)
}

func (a *Ansii) handleByte(b byte) bool {
	if a.converter != nil {
		data := a.converter.HandleByte(b)
		if data != nil {
			a.buf.Write(data)
			return true
		}
	}

	return false
}

func (a *Ansii) handleCsi(csi CsiCode) {
	if a.converter != nil {
		data := a.converter.HandleCsi(csi)
		if data != nil {
			a.buf.Write(data)
		}
	}
}

func (a *Ansii) next(state int) bool {
	a.state = state

	return true
}

func (a *Ansii) change(state int) bool {
	a.state = state

	return false
}
