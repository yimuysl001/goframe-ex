package eid

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/yitter/idgenerator-go/idgen"
)

func init() {
	var options = idgen.NewIdGeneratorOptions(machineID)
	options.BaseTime = gtime.NewFromStr("2025-01-01").Time.UnixMilli()
	idgen.SetIdGenerator(options)
}
