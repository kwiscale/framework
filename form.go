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

type Field struct {
	Tag        string
	Type       string
	Attributes string
	Value      string
	Label      string
	Name       string
}

func (f Field) String() string {

	var ret string

	switch f.Type {
	case "text":
		ret = fmt.Sprintf(TPLTEXTINPUT, f.Label, f.Name, f.Value, f.Attributes)
	case "password":
		ret = fmt.Sprintf(TPLPWDINPUT, f.Label, f.Name, f.Value, f.Attributes)
	case "textarea":
		ret = fmt.Sprintf(TPLTEXTAREA, f.Label, f.Name, f.Attributes, f.Value)
	case "select":
		ret = fmt.Sprintf(TPLSELECT, f.Label, f.Name, f.Attributes, f.Value)
	default:
		panic("error " + f.Type + " not known")
	}

	return string(template.HTML(ret))
}

func (f Field) Html() template.HTML {
	return template.HTML(f.String())
}

type Form struct {
	Fields []Field
}

func (f *Form) Append(field Field) {
	f.Fields = append(f.Fields, field)
}

func renderForm(s interface{}) *Form {

	v := reflect.TypeOf(s)
	value := reflect.ValueOf(s)
	form := new(Form)
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

		field := Field{
			Label:      label,
			Value:      val,
			Attributes: req,
			Name:       f.Name,
			Type:       parts[0],
		}

		form.Append(field)

	}
	return form
}

func CreateForm(T interface{}) *Form {
	return renderForm(T)
}
