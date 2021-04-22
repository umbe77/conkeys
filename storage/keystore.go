package storage

import "errors"

type ValueType int

const (
	Integer ValueType = iota
	Decimal
	String
	DateTime
	Boolean
)

func (v ValueType) ToString() (string, error) {
	switch v {
	case Integer:
		return "Integer", nil
	case Decimal:
		return "Decimal", nil
	case String:
		return "String", nil
	case DateTime:
		return "DateTime", nil
	case Boolean:
		return "Boolean", nil
	}
		return "", errors.New("Not valid type")
}

func checkString(val interface{}) bool {
	_, ok := val.(string)
	return ok
}

func checkBool(val interface{}) bool {
	_, ok := val.(bool)
	return ok
}

func checkInteger(val interface{}) bool {
	_, ok := val.(int)
	return ok
}

func checkDecimal(val interface{}) bool {
	_, ok := val.(float64)
	return ok
}

type Value struct {
	T ValueType
	V interface{}
}

func (v Value) CheckType() bool {
	switch v.T {
	case String:
		return checkString(v.V)
	case Boolean:
		return checkBool(v.V)
	case Decimal:
		return checkDecimal(v.V)
	case Integer:
		return checkInteger(v.V)
	}
	return false
}

type KeyStorage interface {
	Init()
	Get(path string) (Value, error)
	GetKeys(pathSearch string) (map[string]Value, error)
	GetAllKeys() map[string]Value
	Put(path string, value Value)
}
