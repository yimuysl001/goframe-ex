package epongo

import (
	"context"
	"fmt"
	"github.com/flosch/pongo2/v6"
)

//var localTemp = gmap.NewStrStrMap(true)

func BuildFunction(f map[string]any) {
	pongo2.DefaultSet.Globals.Update(f)
}

func ParseContent(ctx context.Context, str string, data pongo2.Context) (out string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	fromString, err := pongo2.FromString(str)
	if data == nil {
		data = pongo2.Context{}
	}
	data["ctx"] = ctx
	if err != nil {
		return "", err
	}
	out, err = fromString.Execute(data)

	return

}

func ParseContentFile(filename string, data pongo2.Context) (out string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	fromString, err := pongo2.FromCache(filename)

	if err != nil {
		return "", err
	}
	out, err = fromString.Execute(data)

	return

}

func RegisterFunction(name string, fun pongo2.FilterFunction) error {
	if pongo2.FilterExists(name) {
		return pongo2.ReplaceFilter(name, fun)
	}
	return pongo2.RegisterFilter(name, fun)

}

func RegisterFunctionMap(mapf map[string]pongo2.FilterFunction) error {
	for s, function := range mapf {
		err := RegisterFunction(s, function)
		if err != nil {
			return err
		}

	}
	return nil
}

//func SetLocalMap(key, temp string) {
//	localTemp.Set(key, temp)
//}
//
//func GetLocalMap(key string) (string, error) {
//	value, found := localTemp.Search(key)
//	if !found {
//		return "", errors.New("key not found")
//	}
//	return value, nil
//}
