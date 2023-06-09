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

	result := make([]constants.MethodsSig, len(uniqueMethods))
	i := 0
	for _, sig := range uniqueMethods {
		result[i] = sig
		i++
	}

	return result
}

func GenerateMethodFunctionForCPPFile(returnType string, methodName string, isSDKType bool, paramList []string, clsName string) string {
	var builder strings.Builder

	for i := range paramList {
		builder.WriteString(", param")
		builder.WriteString(strconv.Itoa(i))
	}

	params := builder.String()

	if isSDKType {
		return `const auto obj = this->env->CallObjectMethod(this->instance, sdk::g_mapper->classes["` + clsName + `"]->methods["` + methodName + `"]` + params + `);
	return std::make_shared<` + returnType + `>(this->env, obj);`
	}

	functionType := java.GetJNIFunctionType(returnType)

	if len(returnType) > 1 && returnType[0] == '[' {
		functionType = "CallObjectMethod"
	}

	return `return this->env->` + functionType + `(this->instance, sdk::g_mapper->classes["` + clsName + `"]->methods["` + methodName + `"]` + params + `);`
}

func GenerateMethodDefinition(methodName string, methodMap constants.MethodMap) string {
	method := ``

	withoutDuplicatedSigs := RemoveDuplicateSigs(methodMap)

	mapDuplicated := make(map[string]int)
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

		suffix := ""
		key := methodName + strings.Join(generics.Map(sig.Params, java.GetJniTypeFromSignature), "")

		if _, ok := mapDuplicated[key]; !ok {
			mapDuplicated[key] = 0
		} else {
			mapDuplicated[key]++
			suffix = strconv.Itoa(mapDuplicated[key])
		}

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

		method += " " + methodName + suffix + `(`

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

	mapDuplicated := make(map[string]int)
	for _, sig := range withoutDuplicatedSigs {
		returnType, isSDKType := GetReturnTypeForSDK(sig.ReturnType)

		objectNameWithC := returnType

		suffix := ""
		key := methodName + strings.Join(generics.Map(sig.Params, java.GetJniTypeFromSignature), "")

		if _, ok := mapDuplicated[key]; !ok {
			mapDuplicated[key] = 0
		} else {
			mapDuplicated[key]++
			suffix = strconv.Itoa(mapDuplicated[key])
		}

		// generate doc
		if len(sig.Params) > 0 {
			comment := "/** \n"
			paramNameForComment := ""
			for i, param := range sig.Params {
				paramNameForComment = "param" + strconv.Itoa(i)
				comment += "* @param " + paramNameForComment + " " + param
			}
			comment += "\n*/\n"

			method += comment
		}

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

		method += namespace + "::C" + clsName + "::" + methodName + suffix + `(`

		paramName := "param0"
		for i, param := range sig.Params {
			paramName = "param" + strconv.Itoa(i)
			method += java.GetJniTypeFromSignature(param) + " " + paramName + ", "
		}

		method = strings.TrimSuffix(method, ", ")

		content := GenerateMethodFunctionForCPPFile(objectNameWithC, methodName, isSDKType, sig.Params, clsName)

		method += `) {
	` + content + `
}
	
`
	}

	return method
}
