package types

type Variables interface {
	GetVar(name string) interface{}
	LookupVar(name string) (interface{}, bool)
	SetVar(name string, value interface{})
	ConsumeVars(config interface{}, consumeMap map[string]string) error
}
