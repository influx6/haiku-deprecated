package reactive

import (
	"fmt"
	"reflect"
)

//AcceptableKind matches the Kind of a type against a interface supplied
func AcceptableKind(ktype reflect.Kind, m interface{}) bool {
	if GetKind(m) == ktype {
		return true
	}
	return false
}

//MakeType validates accepted types and returns the (Immutable, error)
func MakeType(val interface{}, chain bool, m Timer) (Immutable, error) {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Struct:
		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "struct")
	case reflect.Map:
		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "map")
	case reflect.Array:
		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "array")
	case reflect.Slice:
		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "slice")
	}

	return StrictAtom(val, chain, m), nil
}
