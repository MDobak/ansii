package ansii

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHtmlConverter(t *testing.T) {
	a := NewAnsii(NewHtmlConverter())

	_, wErr := a.Write([]byte("a\033[1m\n\033[31mb\033[0m"))
	rb := make([]byte, a.Len())
	_, rErr := a.Read(rb)

	assert.Nil(t, wErr)
	assert.Nil(t, rErr)
	assert.Equal(
		t,
		string("a<span style=\"font-weight:bold;\"><br></span><span style=\"color:#A54242;font-weight:bold;\">b</span>"),
		string(rb),
	)
}

func TestHtmlConverter_HandleByte(t *testing.T) {
	c := NewHtmlConverter()

	assert.Equal(t, []byte("<br>"), c.HandleByte('\n'))
	assert.Equal(t, []byte("&lt;"), c.HandleByte('<'))
	assert.Equal(t, []byte("&gt;"), c.HandleByte('>'))
	assert.Equal(t, []byte("&nbsp;"), c.HandleByte(' '))
	assert.Equal(t, []byte{}, c.HandleByte('\r'))
	assert.Nil(t, c.HandleByte('a'))
}

func TestHtmlConverter_HandleCsi(t *testing.T) {
	c := NewHtmlConverter()

	// bold
	assert.Equal(t, "<span style=\"font-weight:bold;\">", string(c.HandleCsi(CsiCode{"1", 'm'})))

	// underline (should close previous tag)
	assert.Equal(t, "</span><span style=\"font-weight:bold;text-decoration:underline;\">", string(c.HandleCsi(CsiCode{"4", 'm'})))

	// invalid (ignore)
	assert.Equal(t, "", string(c.HandleCsi(CsiCode{"2", 'm'})))

	// reset (close previous tag)
	assert.Equal(t, "</span>", string(c.HandleCsi(CsiCode{"0", 'm'})))

	// reset #2 (tag is already closed so this time sequence should be ignored)
	assert.Equal(t, "", string(c.HandleCsi(CsiCode{"0", 'm'})))

	// background
	assert.Equal(t, "<span style=\"color:#A54242;\">", string(c.HandleCsi(CsiCode{"31", 'm'})))

	// foreground
	assert.Equal(t, "</span><span style=\"color:#A54242;background-color:#8C9440;\">", string(c.HandleCsi(CsiCode{"42", 'm'})))
}
