package sdkutils

import "strings"

func GetReturnTypeForSDK(returnType string) (string, bool) {
	isSDKType := false

	if strings.Contains(returnType, "/") && strings.Contains(returnType, "minecraft") {
		returnType = strings.ReplaceAll(returnType, "/", "::")
		returnType = returnType[1 : len(returnType)-1]
		returnType = "sdk::" + returnType
		isSDKType = true
	} else if len(returnType) > 1 && returnType[0] == '[' {
		// Handle array types
		switch returnType[1:] {
		case "Z":
			returnType = "jbooleanArray"
		case "B":
			returnType = "jbyteArray"
		case "C":
			returnType = "jcharArray"
		case "S":
			returnType = "jshortArray"
		case "I":
			returnType = "jintArray"
		case "J":
			returnType = "jlongArray"
		case "F":
			returnType = "jfloatArray"
		case "D":
			returnType = "jdoubleArray"
		case "V":
			returnType = "voidArray"
		default:
			returnType = "jobjectArray"
		}
	} else {
		// Handle non-array types
		switch returnType {
		case "Z":
			returnType = "jboolean"
		case "B":
			returnType = "jbyte"
		case "C":
			returnType = "jchar"
		case "S":
			returnType = "jshort"
		case "I":
			returnType = "jint"
		case "J":
			returnType = "jlong"
		case "F":
			returnType = "jfloat"
		case "D":
			returnType = "jdouble"
		case "V":
			returnType = "void"
		default:
			returnType = "jobject"
		}
	}
	return returnType, isSDKType
}
