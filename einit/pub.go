package einit

import (
	"context"
	"errors"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
	"time"
)

func GetByKeyByCmd(parser *gcmd.Parser) error {
	opt := parser.GetOpt("p", "")
	return Init(gconv.String(opt))

}
func Init(key ...string) error {
	// 从环境中获取有效期
	if key == nil || len(key) == 0 || key[0] == "" {
		return GetByEnv()
	}

	return GetByKey(key[0])
}

func GetByEnv() error {
	var key = GetKey()

	g.Log().Info(gctx.GetInitCtx(), "业务key:", key)

	err := g.Try(gctx.GetInitCtx(), func(ctx context.Context) {
		err := readPri()
		if err != nil {
			panic(err)
		}
		StartCron(localPri)
	})

	return err

}

func GetByKey(key string) error {
	if gfile.Exists(path) {
		err := readPri()
		if err != nil {
			return err
		}
		if localPri.Key == key {
			g.Log().Info(gctx.GetInitCtx(), "已注册过服务,使用本地配置")
			return GetByEnv()
		}

	}

	decrypt, err := DecryptPBE(key)
	if err != nil {
		return errors.New("解密失败,key：" + GetKey())
	}
	if decrypt == "" {
		return errors.New("解密失败,key：" + GetKey())
	}
	decrypt = strings.ReplaceAll(decrypt, "-", "")
	decrypt = strings.ReplaceAll(decrypt, ":", "")
	decrypt = strings.ReplaceAll(decrypt, " ", "")
	decrypt = strings.ReplaceAll(decrypt, ".", "")

	if len(decrypt) < 8 {
		return errors.New("解密参数不正确,key：" + GetKey())
	}
	parse := time.Now()
	if len(decrypt) < 14 {
		parse, err = time.Parse("20060102", decrypt[:8])
		if err != nil {
			return errors.New("解密参数不正确,key：" + GetKey())
		}
	} else {
		parse, err = time.ParseInLocation("20060102150405", decrypt[:14], time.Local)
		if err != nil {
			return errors.New("解密参数不正确,key：" + GetKey())
		}
	}

	var now = GetNowTime()

	if parse.Before(time.Now()) || parse.Before(now) {
		g.Log().Info(gctx.GetInitCtx(), "设置有效期：", parse)
		return errors.New("有效期小于当前时间,key：" + GetKey())
	}
	var p = Pri{
		BootTime: GetNowTime(),
		Key:      key,
		LastTime: parse,
	}
	StartCron(p)
	pbe, err := EncryptPBE(gjson.New(p).MustToJsonString())
	if err != nil {
		return err
	}

	return gfile.PutContents(path, pbe)

}
