package sdkutils

import (
	"fmt"
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

		method += GetReturnTypeForSDK(sig.ReturnType) + ` `

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
