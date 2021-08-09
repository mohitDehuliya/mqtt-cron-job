package gosdk

type Thing struct {
	Key          string
	TemplateName string
	Name         string
	Description  string
}

func NewThing() *Thing {
	return &Thing{}
}
