package eid

import (
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/sony/sonyflake"
	"github.com/yitter/idgenerator-go/idgen"
	"strconv"
)

const (
	machineID uint16 = 1
)

var (
	localSnowID *sonyflake.Sonyflake
)

func init() {
	localSnowID = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: gtime.NewFromStr("2024-01-01").Time,
		MachineID: GetMachineId,
	})
}

func GetMachineId() (uint16, error) {
	return machineID, nil
}

func GenSnowID() (uint64, error) {
	return localSnowID.NextID()
}

func GenSnowIDStr(base ...int) (string, error) {
	id := idgen.NextId()
	//if err != nil {
	//	return "", err
	//}
	b := 10
	if len(base) > 0 && 2 < base[0] && base[0] < 36 {
		b = base[0]
	}

	return strconv.FormatInt(id, b), nil
}
