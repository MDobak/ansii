package ansii

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnsii_SimpleReadWrite(t *testing.T) {
	a := NewAnsii(NewEscapeCodesCleaner(""))

	rb := make([]byte, 4)
	wl, wErr := a.Write([]byte("test"))
	rl, rErr := a.Read(rb)

	assert.Nil(t, wErr)
	assert.Nil(t, rErr)
	assert.Equal(t, 4, wl)
	assert.Equal(t, 4, rl)
	assert.Equal(t, []byte("test"), rb)
}

func TestAnsii_SimpleReadWriteWithResetCsi(t *testing.T) {
	a := NewAnsii(NewEscapeCodesCleaner(""))

	wl, wErr := a.Write([]byte("a\033[mb"))
	rb := make([]byte, a.Len())
	rl, rErr := a.Read(rb)

	assert.Nil(t, wErr)
	assert.Nil(t, rErr)
	assert.Equal(t, 5, wl)
	assert.Equal(t, 2, rl)
	assert.Equal(t, []byte("ab"), rb)
}

func TestAnsii_SimpleReadWriteWithColorCsi(t *testing.T) {
	a := NewAnsii(NewEscapeCodesCleaner(""))

	wl, wErr := a.Write([]byte("a\033[1;2mb"))
	rb := make([]byte, a.Len())
	rl, rErr := a.Read(rb)

	assert.Nil(t, wErr)
	assert.Nil(t, rErr)
	assert.Equal(t, 8, wl)
	assert.Equal(t, 2, rl)
	assert.Equal(t, []byte("ab"), rb)
}

func TestAnsii_WriteByte(t *testing.T) {
	a := NewAnsii(NewEscapeCodesCleaner(""))

	err1 := a.WriteByte('a')
	err2 := a.WriteByte('b')

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, []byte("ab"), a.Bytes())
}

func TestAnsii_Read(t *testing.T) {
	a := NewAnsii(NewEscapeCodesCleaner(""))

	_, _ = a.Write([]byte("ab"))
	rb1 := make([]byte, 2)
	rl1, err1 := a.Read(rb1)
	rb2 := make([]byte, 2)
	rl2, err2 := a.Read(rb2)

	assert.Nil(t, err1)
	assert.Equal(t, 2, rl1)
	assert.Equal(t, []byte("ab"), rb1)

	assert.Equal(t, io.EOF, err2)
	assert.Equal(t, 0, rl2)
	assert.Equal(t, []byte{0, 0}, rb2)
}

func TestAnsii_ReadBytes(t *testing.T) {
	a := NewAnsii(NewEscapeCodesCleaner(""))

	_, _ = a.Write([]byte("ab cd"))
	rb1, err1 := a.ReadBytes(' ')
	_, err2 := a.ReadBytes(' ')

	assert.Nil(t, err1)
	assert.Equal(t, []byte("ab "), rb1)
	assert.Equal(t, io.EOF, err2)
}

func TestAnsii_ReadRune(t *testing.T) {
	a := NewAnsii(NewEscapeCodesCleaner(""))

	_, _ = a.Write([]byte("\U0001F600"))
	rr1, rl1, err1 := a.ReadRune()
	_, _, err2 := a.ReadRune()

	assert.Nil(t, err1)
	assert.Equal(t, 4, rl1)
	assert.Equal(t, "\U0001F600", string(rr1))
	assert.Equal(t, io.EOF, err2)
}

func TestAnsii_ReadString(t *testing.T) {
	a := NewAnsii(NewEscapeCodesCleaner(""))

	_, _ = a.Write([]byte("ab cd"))
	rb1, err1 := a.ReadString(' ')
	_, err2 := a.ReadString(' ')

	assert.Nil(t, err1)
	assert.Equal(t, "ab ", rb1)
	assert.Equal(t, io.EOF, err2)
}

func TestAnsii_ReadLine(t *testing.T) {
	a := NewAnsii(NewEscapeCodesCleaner(""))

	_, _ = a.Write([]byte("ab\ncd"))
	rb1, err1 := a.ReadLine()
	_, err2 := a.ReadLine()

	assert.Nil(t, err1)
	assert.Equal(t, "ab\n", rb1)
	assert.Equal(t, io.EOF, err2)
}

func TestAnsii_Bytes(t *testing.T) {
	a := NewAnsii(NewEscapeCodesCleaner(""))

	_, _ = a.Write([]byte("ab cd"))
	rb1 := a.Bytes()
	rb2 := a.Bytes()
	_, _ = a.Read(make([]byte, 3))
	rb3 := a.Bytes()

	assert.Equal(t, []byte("ab cd"), rb1)
	assert.Equal(t, []byte("ab cd"), rb2)
	assert.Equal(t, []byte("cd"), rb3)
}

func TestAnsii_String(t *testing.T) {
	a := NewAnsii(NewEscapeCodesCleaner(""))

	_, _ = a.Write([]byte("ab cd"))
	rs1 := a.String()
	rs2 := a.String()
	_, _ = a.Read(make([]byte, 3))
	rs3 := a.String()

	assert.Equal(t, "ab cd", rs1)
	assert.Equal(t, "ab cd", rs2)
	assert.Equal(t, "cd", rs3)
}

func TestAnsii_Len(t *testing.T) {
	a := NewAnsii(NewEscapeCodesCleaner(""))

	_, _ = a.Write([]byte("ab cd"))
	rl1 := a.Len()
	rl2 := a.Len()
	_, _ = a.Read(make([]byte, 3))
	rl3 := a.Len()

	assert.Equal(t, 5, rl1)
	assert.Equal(t, 5, rl2)
	assert.Equal(t, 2, rl3)
}
