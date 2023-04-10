package sdkutils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/generics"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/java"
)

func removeDuplicateMethods(methodMap constants.MethodMap) []constants.MethodsSig {
	uniqueMethods := make(map[string]constants.MethodsSig)

	for _, sig := range methodMap.MethodsSig {
		key := strings.Join(sig.Params, "") + sig.ReturnType
		_, exists := uniqueMethods[key]
		if !exists {
			uniqueMethods[key] = sig
		}
	}

	result := make([]constants.MethodsSig, 0, len(uniqueMethods))
	for _, sig := range uniqueMethods {
		result = append(result, sig)
	}

	return result
}

func GenerateMethodFunctionForCPPFile(returnType string, methodName string, isSDKType bool, paramList []string) string {
	params := ""

	for i := range paramList {
		params += ", param" + strconv.Itoa(i)
	}

	if isSDKType {
		return `const auto obj = this->env->CallObjectMethod(this->instance, g_mapper->methods["` + methodName + `"]` + params + `);
	return std::make_shared<` + returnType + `>(obj)`
	}

	functionType := ""

	if len(returnType) > 1 && returnType[0] == '[' {
		functionType = "CallObjectMethod"
	} else {
		// Handle non-array types
		switch returnType {
		case "jboolean":
			functionType = "CallBooleanMethod"
		case "jbyte":
			functionType = "CallByteMethod"
		case "jchar":
			functionType = "CallCharMethod"
		case "jshort":
			functionType = "CallShortMethod"
		case "jint":
			functionType = "CallIntMethod"
		case "jlong":
			functionType = "CallLongMethod"
		case "jfloat":
			functionType = "CallFloatMethod"
		case "jdouble":
			functionType = "CallDoubleMethod"
		case "void":
			functionType = "CallVoidMethod"
		default:
			functionType = "CallObjectMethod"
		}
	}

	return `return this->env->` + functionType + `(this->instance, g_mapper->methods["` + methodName + `"]` + params + `);`
}

func GenerateMethodFunction(methodName string, methodMap constants.MethodMap) string {
	method := ``

	withoutDuplicated := removeDuplicateMethods(methodMap)

	for _, sig := range withoutDuplicated {

		// SDK Problem - TODO: it should not be there btw
		returnTypeSplitted := strings.Split(sig.ReturnType, "/")
		returnTypeCls := strings.ReplaceAll(returnTypeSplitted[len(returnTypeSplitted)-1], ";", "")
		if len(returnTypeCls) > 2 {
			if !generics.Contains(constants.ClassesToMap, returnTypeCls) {
				fmt.Printf("[ALERT] -> Class %s doesn't exists in ClassesToMap - Used in method: %s! \n", sig.ReturnType, methodName)
			}
		}

		returnType, _ := GetReturnTypeForSDK(sig.ReturnType)

		method += returnType + ` `

		method += methodName + `(`

		for _, param := range sig.Params {
			method += java.GetJniTypeFromSignature(param) + ", "
		}

		method = strings.TrimSuffix(method, ", ")

		method += `);
	
`
	}
	return method
}

func GenerateMethodContent(clsName string, methodName string, methodMap constants.MethodMap) string {
	method := ``

	withoutDuplicated := removeDuplicateMethods(methodMap)

	for _, sig := range withoutDuplicated {
		returnType, isSDKType := GetReturnTypeForSDK(sig.ReturnType)

		method += returnType + ` `

		method += clsName + "::" + methodName + `(`

		paramName := "param0"
		for i, param := range sig.Params {
			paramName = "param" + strconv.Itoa(i)
			method += java.GetJniTypeFromSignature(param) + " " + paramName + ", "
		}

		method = strings.TrimSuffix(method, ", ")

		content := GenerateMethodFunctionForCPPFile(returnType, methodName, isSDKType, sig.Params)

		method += `) {
	` + content + `
}
	
`
	}

	return method
}
