package common_controllers

type By struct {
	Operator       string
	EqualSign      []string
	AttrsName      []string
	AttrsAliasName []string
	AttrsValue     []interface{}
	AttrsOldValue  []interface{}
}
