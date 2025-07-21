package epongo

import (
	"context"
	"github.com/beevik/etree"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	timeconv "github.com/Andrew-M-C/go.timeconv"
	carbon "github.com/dromara/carbon/v2"
	"goframe-ex/ecache"
	"goframe-ex/exml"
	"net/url"
	"regexp"
	"time"
)

func MapFunc() map[string]interface{} {

	return map[string]interface{}{
		"func_ln": func() string {
			return "\n"
		},
		"func_lr": func() string {
			return "\r"
		},
		"func_char": func(char int) string {
			return string(rune(char))
		},
		"func_for": func(count int) []int {
			is := make([]int, count)
			for i := 0; i < count; i++ {
				is[i] = i
			}
			return is
		},

		"func_isMobile": IsMobile,
		"func_isEmail":  IsEmail,
		"func_isURL":    IsURL,
		"func_isIDCard": IsIDCard,
		"func_strToMap": func(data interface{}) map[string]interface{} {
			return gjson.New(data).Map()
		},
		"func_mapToJson": func(data interface{}) string {
			return gjson.New(data).MustToJsonString()
		},
		"func_mapToXml": func(data interface{}, rootTag ...string) string {
			return gjson.New(data).MustToXmlString(rootTag...)
		},
		"func_cfg": func(path string, name ...string) interface{} {
			return g.Cfg(name...).MustGet(context.Background(), path, nil).Val()
		},
		"func_newCarbon": carbon.NewCarbon,
		"func_addDate":   timeconv.AddDate,
		"func_getLocalData": func() *gmap.StrAnyMap {
			return localData
		},
		"func_setLocalCache": func(key string, value any) any {
			localData.Set(key, value)
			return value
		},
		"func_getLocalCache": func(key string) any {
			return localData.Get(key)
		},
		"func_setCommonCache": func(ctx context.Context, key string, value interface{}, duration int) error {
			return ecache.Cache("pongo").Set(ctx, "localpongo:"+key, value, time.Duration(duration)*time.Second)
		},

		"func_getCommonCache": func(ctx context.Context, key string) any {
			get, err := ecache.Cache("pongo").Get(ctx, "localpongo:"+key)
			if err != nil {
				g.Log().Error(ctx, "获取值失败")
				return nil
			}
			return get.Val()
		},
		"func_removeCommonCache": func(ctx context.Context, key ...interface{}) any {
			get, err := ecache.Cache("pongo").Remove(ctx, key...)
			if err != nil {
				g.Log().Error(ctx, "删除值失败")
				return nil
			}
			if get != nil {
				return get.Val()
			}
			return nil
		},
		"func_getJsonData": gjson.New,
		"func_getXmlRoot": func(body string) *etree.Element {
			return exml.GetElementRoot(body)
		},
		"func_getXpathData": exml.GetElementValue,
		"func_getElement":   exml.GetElement,
		"func_getElements":  exml.GetElements,
		"func_getXmlDoc":    exml.GetMustDoc,
	}
}

// IsMobile 是否为手机号码
func IsMobile(mobile string) bool {
	pattern := `^(1[2|3|4|5|6|7|8|9][0-9]\d{8})$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(mobile)
}

// IsEmail 是否为邮箱地址
func IsEmail(email string) bool {
	// pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z].){1,4}[a-z]{2,4}$` //匹配电子邮箱
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// IsURL 是否是url地址
func IsURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}
	URL, err := url.Parse(u)
	if err != nil || URL.Scheme == "" || URL.Host == "" {
		return false
	}
	return true
}

// IsIDCard 是否为身份证
func IsIDCard(idCard string) bool {
	sz := len(idCard)
	if sz != 18 {
		return false
	}
	weight := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	validate := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
	sum := 0
	for i := 0; i < len(weight); i++ {
		sum += weight[i] * int(byte(idCard[i])-'0')
	}
	m := sum % 11
	return validate[m] == idCard[sz-1]
}
