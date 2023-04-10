package sdkstruct

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
	sdkutils "github.com/pedrogpo/mc-auto-mapper/internal/sdk/utils"
)

func GenerateHeaderContent(clsPath string, allMappings constants.Mappings, includesFile string) {
	// Remove the file name from the path
	// path := "out/sdk/" + strings.Replace(clsPath, "net/", "", 1)

	path := "out/sdk/" + clsPath

	clsPathParts := strings.Split(clsPath, "/")
	clsName := clsPathParts[len(clsPathParts)-1]

	methods := ``
	imports := ``

	returnTypesToImport := []string{}

	for methodName, methodMap := range constants.GetMethodsToMapInClass(allMappings, clsName) {
		withoutDuplicatedSigs := sdkutils.RemoveDuplicateSigs(methodMap)

		// imports
		for _, sig := range withoutDuplicatedSigs {
			returnType, isSDKType := sdkutils.GetReturnTypeForSDK(sig.ReturnType)
			pathToImport := strings.ReplaceAll(returnType, "::", "/")

			if strings.ReplaceAll(pathToImport, "sdk/", "") == strings.ReplaceAll(clsPath, "sdk/", "") {
				continue
			}

			if isSDKType {
				valueAlreadyInList := false
				for _, existingValue := range returnTypesToImport {
					if existingValue == pathToImport {
						valueAlreadyInList = true
						break
					}
				}
				if !valueAlreadyInList {
					returnTypesToImport = append(returnTypesToImport, pathToImport)
				}
			}
		}

		methods += sdkutils.GenerateMethodDefinition(methodName, methodMap)
	}

	for _, returnType := range returnTypesToImport {
		imports += `#include "` + returnType + `.hpp"
`
	}

	// write content in .hpp file
	hpp := `#pragma once
#include "` + includesFile + `"
`

	hpp += imports

	namespace := "sdk::"
	for i := 0; i < len(clsPathParts)-1; i++ {
		namespace += clsPathParts[i] + "::"
	}

	namespace = strings.TrimSuffix(namespace, "::")

	hpp += `
namespace ` + namespace + ` {
	class C` + clsName + ` {
	private:
		JNIEnv* env;
		jobject instance;
	public:
		C` + clsName + `(JNIEnv* env);
		C` + clsName + `(JNIEnv* env, jobject instance);
		jobject getInstance() { return this->instance; }

		~C` + clsName + `();

`
	hpp += methods
	hpp += `	};
}
`

	err := ioutil.WriteFile(path+".hpp", []byte(hpp), 0644)

	if err != nil {
		fmt.Println(err)
		return
	}
}
