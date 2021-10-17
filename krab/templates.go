package krab

type Templates struct {
}

func (t Templates) ProcessArguments(tpl string, args map[string]interface{}) (string, error) {
	return tpl, nil
}

func EmptyArgs() map[string]interface{} {
	return map[string]interface{}{}
}
