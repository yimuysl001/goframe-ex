package ejson

import (
	"bytes"
	"context"
	"encoding/json"
	"regexp"
	"strings"
)

func IndentJson(ctx context.Context, str string) string {
	var bb bytes.Buffer
	err := json.Indent(&bb, []byte(str), "", "    ")
	if err != nil {
		panic(err)
	}
	return bb.String()

}

func ReplaceJson(ctx context.Context, template string) string {
	var jsoncomma = `,\s*\}`
	compile1 := regexp.MustCompile(jsoncomma)
	template = compile1.ReplaceAllString(template, "}")

	jsoncomma = `,\s*\]`
	compile1 = regexp.MustCompile(jsoncomma)
	template = compile1.ReplaceAllString(template, "]")

	var yhstr = `"(?s:(.*?))"`
	compile := regexp.MustCompile(yhstr)
	allString := compile.FindAllString(template, -1)
	//
	splits := compile.Split(template, -1)

	var sb strings.Builder
	sb.WriteString(splits[0])
	for i, split := range allString {
		split = strings.ReplaceAll(split, "\n", "\\n")
		split = strings.ReplaceAll(split, "\r", "\\r")
		sb.WriteString(split)
		sb.WriteString(splits[i+1])
	}

	return strings.ReplaceAll(sb.String(), "&quot;", `\"`)

}
