package kwiscale

// Generates fields from a struct
// Use Tags to set type
//
// Example:
//    type P struct {
//        Name string `form:"text,The label,required"
//    }

import (
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

var TPLTEXTINPUT string = `
<label> %s
<input type="text" name="%s" value="%s" %s/>
</label>
`

var TPLPWDINPUT string = `
<label> %s
<input type="password" name="%s" value="%s" %s/>
</label>
`

var TPLTEXTAREA string = `
<label> %s
<textarea name="%s" %s>%s</textarea>
</label>
`

var TPLSELECT string = `
<label> %s
<select name="%s" %s>%s</select>
</label>
`

var TPLOPTION string = `
    <option value="%s">%s</option>`

func renderForm(s interface{}) []string {

	ret := []string{}
	v := reflect.TypeOf(s)
	value := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {

		f := v.Field(i)
		if len(f.Tag) == 0 {
			continue
		}

		label := f.Name
		var val string
		parts := strings.Split(f.Tag.Get("form"), ",")

		switch f.Type.Kind() {
		case reflect.Slice, reflect.Array:
			valuefield := parts[3]
			optionfield := parts[4]

			sl := value.Field(i).Interface()
			intface := reflect.ValueOf(sl)

			for j := 0; j < intface.Len(); j++ {
				opt := intface.Index(j)
				val += fmt.Sprintf(TPLOPTION, opt.FieldByName(valuefield).String(), opt.FieldByName(optionfield).String())
			}

		default:
			val = value.Field(i).String()
		}

		req := ""
		if len(parts) > 1 {
			if parts[1] != "" {
				label = parts[1]
			}
			if len(parts) > 2 {
				req = parts[2]
			}
		}

		switch parts[0] {
		case "text":
			ret = append(ret, fmt.Sprintf(TPLTEXTINPUT, label, f.Name, val, req))
		case "password":
			ret = append(ret, fmt.Sprintf(TPLPWDINPUT, label, f.Name, val, req))
		case "textarea":
			ret = append(ret, fmt.Sprintf(TPLTEXTAREA, label, f.Name, req, val))
		case "select":
			ret = append(ret, fmt.Sprintf(TPLSELECT, label, f.Name, req, val))
		default:
			panic("error " + parts[0] + " not known")
		}
	}
	return ret
}

func CreateForm(T interface{}) template.HTML {

	s := renderForm(T)
	r := ""
	for _, i := range s {
		r += i
	}

	return template.HTML(r)

}
