package sdkutils

import (
	"strings"

	"github.com/pedrogpo/mc-auto-mapper/internal/utils/java"
)

func GetReturnTypeForSDK(returnType string) (string, bool) {
	isSDKType := false

	if strings.Contains(returnType, "/") && strings.Contains(returnType, "minecraft") {
		returnType = strings.ReplaceAll(returnType, "/", "::")
		returnType = returnType[1 : len(returnType)-1]
		returnType = "sdk::" + returnType
		isSDKType = true
	} else if len(returnType) > 1 && returnType[0] == '[' {
		returnType = java.HandleJNIArrayType(returnType[1:])
	} else {
		returnType = java.HandleJNINonArrayType(returnType)
	}
	return returnType, isSDKType
}
