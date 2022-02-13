package storage

import (
	"errors"
	"time"
)

type ValueType int

const (
	Integer ValueType = iota
	Decimal
	String
	DateTime
	Boolean
	Crypted
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
	v, ok := val.(float64)
	return ok && v == float64(int(v))
}

func checkDecimal(val interface{}) bool {
	_, ok := val.(float64)
	return ok
}

func checkDateTime(val interface{}) bool {
	str, ok := val.(string)
	if ok {
		_, err := time.Parse(time.RFC3339, str)
		ok = err == nil
	}
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
	case DateTime:
		return checkDateTime(v.V)
	case Crypted:
		return checkString(v.V)
	}
	return false
}

type KeyStorage interface {
	Init()
	Get(path string) (Value, error)
	GetEncrypted(path string) (Value, error)
	GetKeys(pathSearch string) (map[string]Value, error)
	GetAllKeys() map[string]Value
	Put(path string, value Value)
	PutEncrypted(path string, maskedValue Value, encryptedValue string)
	Delete(path string)
}
