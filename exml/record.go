package exml

import (
	"strings"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
)

func RecordToXml(vs gdb.Record, rootTag ...string) string {
	var strmap = make(map[string]string)
	for k, v := range vs {
		strmap[k] = v.String()
	}
	//Mylog.MyInfo().Println("strmap：", strmap)
	xmlString := gjson.New(strmap).MustToXmlIndentString(rootTag...)
	//Mylog.MyInfo().Println("xmlString：", xmlString)
	return xmlString
}

func ResultToXml(vs gdb.Result, rootTag ...string) string {
	var sb strings.Builder

	for _, v := range vs {
		if len(rootTag) > 0 {
			sb.WriteString(RecordToXml(v, rootTag[0]))
		} else {
			sb.WriteString(RecordToXml(v))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
