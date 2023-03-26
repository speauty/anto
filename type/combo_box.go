package _type

type StdComboBoxModel struct {
	Key  string
	Name string
}

func (customSTD *StdComboBoxModel) BindKey() string {
	return "Key"
}

func (customSTD *StdComboBoxModel) DisplayKey() string {
	return "Name"
}
