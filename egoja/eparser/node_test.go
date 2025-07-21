package eparser

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/util/gmeta"
	"testing"
)

func TestSql(t *testing.T) {

	sql, parameters, err := ParseSql(`
select *  from yxhis${aaa.substring(0,4)}..tbmzfymx${aaa.substring(0,6)} where 1=1 
      ?{cbrh , and cmzh=#{cbrh} }
       ?{ typeof cbrid !== 'undefined' && cbrid!=null  && cbrid!='' ,   and cbrid=#{cbrid} }
      ?{ ighzl!=null&&ighzl.length>0 , and ighzl in( #{ighzl} )}
        `, map[string]interface{}{"aaa": "20250506", "cbrh": "1234", "cbrid": "", "ighzl": []int{1, 2, 3}})

	fmt.Println(sql, parameters, err)

}

func TestNameb(t *testing.T) {
	sql, parameters, err := ParseSql(`
select *   from yxhis${aaa.substring(0,4)}..tbmzfymx${aaa.substring(0,6)} where 1=1 
      ?{cbrh , and cmzh=#{cbrh} }
       ?{  !cbrid    || cbrid =='' ,   and cbrid="" }
      ?{ ighzl!=null&&ighzl.length>0 , and ighzl in( #{ighzl} )}
        `, map[string]interface{}{"aaa": "20250506", "cbrh": "1234", "cbrid": "", "ighzl": []int{1, 2, 3}})

	fmt.Println(sql, parameters, err)
}

type MetaTag struct {
	g.Meta `tableName:"meta_tag" comment:"meta_tag" `
	Id     string `json:"id,omitempty" orm:"id" comment:"ID"`
	Name   string `json:"name,omitempty" orm:"name" comment:"说明"`
}

func TestTagAndMeta(t *testing.T) {
	var m MetaTag
	fmt.Println(gmeta.Data(m))
	fmt.Println(gstructs.TagFields(m, []string{"orm", "gorm"}))

}

type Ctx struct {
	Name string
	ctx  context.Context
}

func (c *Ctx) GetCtx() context.Context {
	return c.ctx
}

func GetCtxs(ctx context.Context) *Ctx {
	value := ctx.Value("ctxc")
	if value == nil {
		return nil
	}
	return value.(*Ctx)
}

func TestDpS(t *testing.T) {
	var (
		ctxc = &Ctx{Name: "ctxa"}
		ctx  = context.WithValue(context.Background(), "ctxc", ctxc)
	)
	ctxc.ctx = ctx

	fmt.Println(GetCtxs(GetCtxs(ctx).ctx).Name)

}
