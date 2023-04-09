package joined

import "strings"

type JoinedType int

const (
	CLASS JoinedType = iota
	FIELD
	METHOD
	PK
	UNDEFINED
)

func GetJoinedType(str string) JoinedType {
	if strings.HasPrefix(str, "CL: ") {
		return CLASS
	} else if strings.HasPrefix(str, "FD: ") {
		return FIELD
	} else if strings.HasPrefix(str, "MD: ") {
		return METHOD
	} else if strings.HasPrefix(str, "PK: ") {
		return PK
	} else {
		return UNDEFINED
	}
}
