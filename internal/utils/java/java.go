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

func GetJniTypeFromSignature(str string) string {
	if len(str) > 1 && str[0] == '[' {
		// Handle array types
		switch str[1:] {
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
		switch str {
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
