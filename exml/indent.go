package exml

import (
	"context"
	"github.com/beevik/etree"
)

func IndentXml(ctx context.Context, str string) string {

	document := etree.NewDocument()
	err := document.ReadFromString(str)
	if err != nil {
		panic(err)
	}

	document.IndentTabs()
	toString, err := document.WriteToString()
	if err != nil {
		panic(err)
	}
	return toString

}
