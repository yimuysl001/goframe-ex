package econf

import (
	"github.com/flosch/pongo2/v6"
	"github.com/gogf/gf/v2/frame/g"
	"goframe-ex/encrypt/ucrypt"
)

const (
	key = "myyixinfixkdsxdo"
	iv  = "rfgtecpoxserfghu"
)

var (
	localKey = ""
	localIv  = ""
)

func SetLocalKey(ckey string) {
	localKey = ckey
}
func SetLocalVi(civ string) {
	localIv = civ
}
func getLocalKey() string {
	if localKey == "" {
		return key
	}
	return localKey
}

func getLocalIv() string {
	if localIv == "" {
		return iv
	}

	return localIv
}

func init() {
	pongo2.DefaultSet.Globals.Update(map[string]any{
		"dec": func(data string) string {
			decrypt2, err := ucrypt.Decrypt2(data, getLocalKey(), 5728)
			if err != nil {
				panic(err)
			}
			return decrypt2
		},
	})

}

// InitConf 注意需最优先加载
func InitConf(file ...string) {
	g.Cfg().SetAdapter(NewAdapterFile(file...))
}
