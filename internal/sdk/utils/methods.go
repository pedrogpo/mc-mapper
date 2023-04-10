package sdkutils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/generics"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/java"
)

func RemoveDuplicateSigs(methodMap constants.MethodMap) []constants.MethodsSig {
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
	return std::make_shared<` + returnType + `>(obj);`
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

func GenerateMethodDefinition(methodName string, methodMap constants.MethodMap) string {
	method := ``

	withoutDuplicatedSigs := RemoveDuplicateSigs(methodMap)

	for _, sig := range withoutDuplicatedSigs {

		// SDK Problem - TODO: it should not be there btw
		returnTypeSplitted := strings.Split(sig.ReturnType, "/")
		returnTypeCls := strings.ReplaceAll(returnTypeSplitted[len(returnTypeSplitted)-1], ";", "")
		if len(returnTypeCls) > 2 {
			if !generics.Contains(constants.ClassesToMap, returnTypeCls) {
				fmt.Printf("[ALERT] -> Class %s doesn't exists in ClassesToMap - Used in method: %s! \n", sig.ReturnType, methodName)
			}
		}

		returnType, isSDKType := GetReturnTypeForSDK(sig.ReturnType)

		method += "		"

		if isSDKType {
			method += `std::shared_ptr<`

			parts := strings.Split(returnType, "::")
			parts[len(parts)-1] = "C" + parts[len(parts)-1]

			method += strings.Join(parts, "::")
		} else {
			method += returnType
		}

		if isSDKType {
			method += `>`
		}

		method += " " + methodName + `(`

		for _, param := range sig.Params {
			method += java.GetJniTypeFromSignature(param) + ", "
		}

		method = strings.TrimSuffix(method, ", ")

		method += `);
`
	}
	return method
}

func GenerateMethodContent(clsName string, methodName string, methodMap constants.MethodMap, namespace string) string {
	method := ``

	withoutDuplicatedSigs := RemoveDuplicateSigs(methodMap)

	for _, sig := range withoutDuplicatedSigs {
		returnType, isSDKType := GetReturnTypeForSDK(sig.ReturnType)

		objectNameWithC := returnType

		if isSDKType {
			method += `std::shared_ptr<`

			parts := strings.Split(returnType, "::")
			parts[len(parts)-1] = "C" + parts[len(parts)-1]
			objectNameWithC = strings.Join(parts, "::")

			method += objectNameWithC
		} else {
			method += returnType + " "
		}

		if isSDKType {
			method += `> `
		}

		method += namespace + "::C" + clsName + "::" + methodName + `(`

		paramName := "param0"
		for i, param := range sig.Params {
			paramName = "param" + strconv.Itoa(i)
			method += java.GetJniTypeFromSignature(param) + " " + paramName + ", "
		}

		method = strings.TrimSuffix(method, ", ")

		content := GenerateMethodFunctionForCPPFile(objectNameWithC, methodName, isSDKType, sig.Params)

		method += `) {
	` + content + `
}
	
`
	}

	return method
}
