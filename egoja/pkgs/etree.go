package pkgs

import (
	"github.com/beevik/etree"
	"goframe-ex/egoja/gojaapi"
	"goframe-ex/exml"
)

func init() {
	gojaapi.RegisterImport("etree", map[string]any{
		"ErrXML":        etree.ErrXML,
		"ReadSettings":  etree.ReadSettings{},
		"WriteSettings": etree.WriteSettings{},
		//"Token":         etree.Token{},
		"Document":  etree.Document{},
		"Element":   etree.Element{},
		"Attr":      etree.Attr{},
		"CharData":  etree.CharData{},
		"Comment":   etree.Comment{},
		"Directive": etree.Directive{},
		"ProcInst":  etree.ProcInst{},
		"Path":      etree.Path{},

		"NewDocument":     etree.NewDocument,
		"NewText":         etree.NewText,
		"NewCData":        etree.NewCData,
		"NewCharData":     etree.NewCharData,
		"NewComment":      etree.NewComment,
		"NewDirective":    etree.NewDirective,
		"NewProcInst":     etree.NewProcInst,
		"CompilePath":     etree.CompilePath,
		"MustCompilePath": etree.MustCompilePath,
		"TraverseXmlRoot": exml.TraverseXmlRoot,
		"SetElementValue": exml.SetElementValue,
		"DelAttr":         exml.DelAttr,
		"GetElement":      exml.GetElement,
		"GetElements":     exml.GetElements,
		"GetROOT":         exml.GetROOT,
		"GetElementRoot":  exml.GetElementRoot,
		"ResultToXml":     exml.ResultToXml,
		"RecordToXml":     exml.RecordToXml,
		"GetMustDoc":      exml.GetMustDoc,
		"GetDoc":          exml.GetDoc,
	})
}
