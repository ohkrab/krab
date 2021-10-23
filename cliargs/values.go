package cliargs

type Values map[string]interface{}

func (v Values) Get(key string) string {
	val, exists := v[key]
	if exists {
		if s, ok := val.(string); ok {
			return s
		}
	}

	return ""
}
