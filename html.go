package ansii

import (
	"strconv"
)

const (
	controlCharMax  = 0x1f
)

type HtmlConverter struct {
	tagOpened  bool
	foreground string
	background string
	bold       bool
	underline  bool
}

var colors = [16]string{
	// Regular
	"#282A2E",
	"#A54242",
	"#8C9440",
	"#DE935F",
	"#5F819D",
	"#85678F",
	"#85678F",
	"#707880",

	// Bright
	"#373B41",
	"#CC6666",
	"#B5BD68",
	"#F0C674",
	"#81A2BE",
	"#B294BB",
	"#8ABEB7",
	"#C5C8C6",
}

func NewHtmlConverter() *HtmlConverter {
	return &HtmlConverter{
		tagOpened:  false,
		foreground: "",
		background: "",
		bold:       false,
		underline:  false,
	}
}

func (p *HtmlConverter) HandleByte(b byte) []byte {
	if ' ' == b {
		return []byte("&nbsp;")
	}

	if '\n' == b {
		return []byte("<br>")
	}

	if '<' == b {
		return []byte("&lt;")
	}

	if '>' == b {
		return []byte("&gt;")
	}

	if b <= controlCharMax {
		return []byte{}
	}

	return nil
}

func (p *HtmlConverter) HandleCsi(csi CsiCode) []byte {
	if 'm' == csi.Code {
		if "0" == csi.Arg {
			p.bold = false
			p.underline = false
			p.background = ""
			p.foreground = ""
			
			return p.printCloseTag()
		}

		// Bold
		if "1" == csi.Arg {
			p.bold = true
			
			return p.printOpenTag()
		}

		// Underline
		if "4" == csi.Arg {
			p.underline = true
			
			return p.printOpenTag()
		}

		// Regular foreground
		for i := 0; i < 8; i++ {
			if "3"+strconv.Itoa(i) == csi.Arg {
				p.foreground = colors[i]
				
				return p.printOpenTag()
			}
		}

		// Bright foreground
		for i := 8; i < 16; i++ {
			if i > 8 && "3"+strconv.Itoa(i-8)+";1" == csi.Arg {
				p.foreground = colors[i]
				
				return p.printOpenTag()
			}
		}

		// Regular background
		for i := 0; i < 8; i++ {
			if "4"+strconv.Itoa(i) == csi.Arg {
				p.background = colors[i]
				
				return p.printOpenTag()
			}
		}

		// Bright background
		for i := 8; i < 16; i++ {
			if i > 8 && "4"+strconv.Itoa(i-8)+";1" == csi.Arg {
				p.background = colors[i]
				
				return p.printOpenTag()
			}
		}
	}

	return nil
}

func (p *HtmlConverter) printOpenTag() []byte {
	tag := ""

	if p.tagOpened {
		tag += string(p.printCloseTag())
	}

	p.tagOpened = true

	tag += `<span style="`

	if "" != p.foreground {
		tag += `color:`+p.foreground+`;`
	}

	if "" != p.background {
		tag += `background-color:`+p.background+`;`
	}

	if p.bold {
		tag += `font-weight:bold;`
	}

	if p.underline {
		tag += `text-decoration:underline;`
	}

	tag += `">`

	return []byte(tag)
}

func (p *HtmlConverter) printCloseTag() []byte {
	if p.tagOpened {
		p.tagOpened = false
		return []byte(`</span>`)
	}

	return nil
}
