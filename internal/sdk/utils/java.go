package sdkutils

import "strings"

func GetReturnTypeForSDK(returnType string) string {
	if strings.Contains(returnType, "/") && strings.Contains(returnType, "minecraft") {
		returnType = strings.ReplaceAll(returnType, "/", "::")
		returnType = returnType[1 : len(returnType)-1]
		return "sdk::" + returnType
	}

	if len(returnType) > 1 && returnType[0] == '[' {
		// Handle array types
		switch returnType[1:] {
		case "Z":
			return "jbooleanArray"
		case "B":
			return "jbyteArray"
		case "C":
			return "jcharArray"
		case "S":
			return "jshortArray"
		case "I":
			return "jintArray"
		case "J":
			return "jlongArray"
		case "F":
			return "jfloatArray"
		case "D":
			return "jdoubleArray"
		case "V":
			return "voidArray"
		default:
			return "jobjectArray"
		}
	} else {
		// Handle non-array types
		switch returnType {
		case "Z":
			return "jboolean"
		case "B":
			return "jbyte"
		case "C":
			return "jchar"
		case "S":
			return "jshort"
		case "I":
			return "jint"
		case "J":
			return "jlong"
		case "F":
			return "jfloat"
		case "D":
			return "jdouble"
		case "V":
			return "void"
		default:
			return "jobject"
		}
	}
}
