package krab

type Templates struct {
}

func (t Templates) ProcessArguments(tpl string, args map[string]any) (string, error) {
	return tpl, nil
}

func EmptyArgs() map[string]any {
	return map[string]any{}
}
