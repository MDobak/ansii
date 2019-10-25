package ansii

import (
	"strings"
)

const CleanAll = ""
const AllowColors = "m"

type EscapeCodesCleaner struct {
	LastCsiCode        CsiCode
	LastAllowedCsiCode CsiCode

	allowedCsiModes string
}

func NewEscapeCodesCleaner(allowedCsiModes string) *EscapeCodesCleaner {
	return &EscapeCodesCleaner{
		allowedCsiModes: allowedCsiModes,
	}
}

func (p *EscapeCodesCleaner) HandleByte(b byte) []byte {
	return nil
}

func (p *EscapeCodesCleaner) HandleCsi(csi CsiCode) []byte {
	p.LastCsiCode = csi

	if strings.Contains(p.allowedCsiModes, string(csi.Code)) {
		p.LastAllowedCsiCode = csi
		return csi.Bytes()
	}

	return nil
}
