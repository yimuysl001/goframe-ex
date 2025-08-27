package pkgs

import (
	"github.com/Masterminds/sprig/v3"
	"goframe-ex/egoja/gojaapi"
)

func init() {

	for s, a := range sprig.FuncMap() {
		gojaapi.RegisterCommonParameter("f_"+s, a)
	}

}
