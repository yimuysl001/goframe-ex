package epongo

import (
	"github.com/clbanning/mxj/v2"
	"github.com/flosch/pongo2/v6"
	"github.com/Masterminds/sprig/v3"
)

func init() {
	pongo2.DefaultSet.Debug = false
	mxj.XMLEscapeChars(true)
	BuildFunction(sprig.FuncMap())
	BuildFunction(MapFunc())
}
