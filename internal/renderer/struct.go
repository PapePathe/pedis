package renderer

import (
	"fmt"
	"reflect"
	"strings"
)

type RespRenderer struct {
	Tagname string
}

func (rr RespRenderer) Render(i interface{}) []byte {
	data := []byte{}
	value := reflect.ValueOf(i)

	fields := fmt.Sprintf("%%%d", value.NumField())
	data = append(data, []byte(fields)...)
	data = append(data, []byte{13, 10}...)

	for i := 0; i < value.NumField(); i++ {
		tag := value.Type().Field(i).Tag.Get(rr.Tagname)

		if tag == "" || tag == "-" {
			continue
		}

		args := strings.Split(tag, ",")
		e := value.Field(i).Interface()
		name := strings.ToLower(value.Type().Field(i).Name)

		switch args[0] {
		case "string":
			r := StringRenderer{}
			data = append(data, r.Render(name, e.(string))...)
		case "int":
			r := IntRenderer{}
			data = append(data, r.Render(name, e.(int))...)
		default:
			panic("data type not yet handled")
		}
	}

	return data
}
