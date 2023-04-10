package java

import (
	"regexp"
)

func ExtractParamsAndReturn(signature string) (params []string, returnType string) {
	paramRegex := regexp.MustCompile(`\((.*?)\)(.+)`)
	matches := paramRegex.FindStringSubmatch(signature)

	if len(matches) == 3 {
		paramsString := matches[1]
		returnType = matches[2]

		paramTypeRegex := regexp.MustCompile(`(?:\[(?:\[)*|)([ZBCSIJFD]|L[^;]+;)`)
		paramMatches := paramTypeRegex.FindAllStringSubmatch(paramsString, -1)

		for _, match := range paramMatches {
			params = append(params, match[1])
		}
	}

	return params, returnType
}

// Function to handle function types
func GetJNIFunctionType(returnType string) string {
	functionTypeMap := map[string]string{
		"jboolean": "CallBooleanMethod",
		"jbyte":    "CallByteMethod",
		"jchar":    "CallCharMethod",
		"jshort":   "CallShortMethod",
		"jint":     "CallIntMethod",
		"jlong":    "CallLongMethod",
		"jfloat":   "CallFloatMethod",
		"jdouble":  "CallDoubleMethod",
		"void":     "CallVoidMethod",
	}

	functionType, ok := functionTypeMap[returnType]
	if !ok {
		functionType = "CallObjectMethod"
	}

	return functionType
}

// Function to handle array types
func HandleJNIArrayType(returnType string) string {
	switch returnType {
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
	return returnType
}

// Function to handle non-array types
func HandleJNINonArrayType(returnType string) string {
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
	return returnType
}

func GetJniTypeFromSignature(str string) string {
	if len(str) > 1 && str[0] == '[' {
		return HandleJNIArrayType(str)
	}
	return HandleJNINonArrayType(str)
}
