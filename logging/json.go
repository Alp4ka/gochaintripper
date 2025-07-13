package logging

type JSONGetter interface {
	IsJSON(json string) bool
	Exists(json string, fieldPattern string) bool
}

type JSONSetter interface {
	SetValue(json string, fieldPattern, replaceString string) (string, error)
}
