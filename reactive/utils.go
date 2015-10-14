package reactive

// const (
// 	// ErrUnacceptedTypeMessage defines the message for types that are not part of the basic units/types in go
// 	ErrUnacceptedTypeMessage = "Type %s is not acceptable"
// )

// GetKind returns the kind of the value
// func GetKind(m interface{}) reflect.Kind {
// 	return reflect.TypeOf(m).Kind()
// }

// // AcceptableKind matches the Kind of a type against a interface supplied
// func AcceptableKind(ktype reflect.Kind, m interface{}) bool {
// 	if GetKind(m) == ktype {
// 		return true
// 	}
// 	return false
// }

// MakeType validates accepted types and returns the (Immutable, error)
// func MakeType(val interface{}, chain bool) (Immutable, error) {
// 	switch reflect.TypeOf(val).Kind() {
// 	case reflect.Struct:
// 		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "struct")
// 	case reflect.Map:
// 		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "map")
// 	case reflect.Array:
// 		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "array")
// 	case reflect.Slice:
// 		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "slice")
// 	}
//
// 	return StrictAtom(val, chain), nil
// }
