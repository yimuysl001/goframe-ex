package einit

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"goframe-ex/encrypt/ucrypt"
	"github.com/shirou/gopsutil/host"
	"os"
	"time"
)

// 设置
func GetKey() string {
	var key = ""
	dir, err := os.Getwd()
	if err != nil {
		key = key + err.Error()
	} else {
		key = key + dir
	}
	info, err := host.Info()
	if err != nil {
		key = key + err.Error()
	} else {
		key = key + info.HostID
	}

	return gmd5.MustEncrypt(key)
}
func GetNowInt64() uint64 {
	loacalInfo, _ := host.Info()
	return loacalInfo.BootTime + loacalInfo.Uptime
}
func GetNowTime() time.Time {
	return time.Unix(int64(GetNowInt64()), 0)
}

func readPri() error {
	getenv := gfile.GetContents(path)

	//getenv := os.Getenv(envkeyPre + key)
	if getenv == "" {
		return errors.New("未设置有效期")
	}

	pbe, err := DecryptPBE(getenv)
	if err != nil {
		return errors.New("有效期获取失败")
	}

	localPri = Pri{}

	err = gjson.New(pbe).Scan(&localPri)
	if err != nil {
		return err
	}

	return nil

}

func check(ctx context.Context, p Pri) {
	cnow := GetNowTime()
	//g.Log().Debug(ctx, "当前时间：", cnow)
	//g.Log().Debug(ctx, "创建时间：", p.BootTime)
	if cnow.Before(p.BootTime) {
		g.Log().Error(ctx, "时间获取异常")
		os.Exit(1)
	}
	//g.Log().Debug(ctx, "有效时间：", p.LastTime)
	if cnow.After(p.LastTime) {
		g.Log().Error(ctx, "程序已过期，过期时间：", p.LastTime)
		g.Log().Error(ctx, "请联系管理员续期，key：", GetKey())
		os.Exit(1)
	}
}
func DecryptPBE(pwdstr string) (string, error) {
	var key = GetKey()
	var str = ""
	var err error
	err = g.Try(gctx.GetInitCtx(), func(ctx context.Context) {
		str, err = ucrypt.Decrypt2(pwdstr, key, 5728)
		//str, err = encryption.Decrypt(pwdstr, getLocalKey()+":"+key)
		if err != nil {
			panic(err)
		}
	})

	return str, err
}

func EncryptPBE(str string) (string, error) {
	var key = GetKey()

	return ucrypt.Encrypt2(str, key, 5728)

	//return encryption.Encrypt(str, getLocalKey()+":"+key)
}
