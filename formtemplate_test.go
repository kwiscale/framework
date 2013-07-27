package kwiscale

import "testing"

type GenderType struct {
	Type string
	Ref  string
}

type T struct {
	Name    string       `form:"text,The name"`
	Passwd  string       `form:"password,The password"`
	Address string       `form:"textarea,The address"`
	Gender  []GenderType `form:"select,Gender,,Ref,Type"`
}

func TestFormTemplateInput(t *testing.T) {

	genders := []GenderType{}
	genders = append(genders, GenderType{"Male", "M"})
	genders = append(genders, GenderType{"Female", "F"})

	data := T{"Patrice", "mypass", "My Address", genders}

	res := renderForm(data)
    if len(res.Fields) != 4 {
        t.Error( "not good fields number: ", len(res.Fields))
    }
}
