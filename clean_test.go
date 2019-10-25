package ansii

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanConverter(t *testing.T) {
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
