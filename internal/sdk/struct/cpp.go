package sdkstruct

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
	sdkutils "github.com/pedrogpo/mc-auto-mapper/internal/sdk/utils"
)

func GenerateCppContent(clsPath string, allMappings constants.Mappings) {
	// Remove the file name from the path
	// path := "out/sdk/" + strings.Replace(clsPath, "net/", "", 1)

	path := "out/sdk/" + clsPath

	clsPathParts := strings.Split(clsPath, "/")
	clsName := clsPathParts[len(clsPathParts)-1]

	// write content in .cpp file
	cpp := `#include "` + clsName + `.hpp"
#include <sdk/mapper.hpp>

`

	namespace := "sdk::"
	for i := 0; i < len(clsPathParts)-1; i++ {
		namespace += clsPathParts[i] + "::"
	}

	namespace = strings.TrimSuffix(namespace, "::")

	cpp += namespace + `::C` + clsName + `::C` + clsName + `(JNIEnv* env) {
	this->env = env;
}

` + namespace + `::C` + clsName + `::C` + clsName + `(JNIEnv* env, jobject instance) : instance(instance) {
	this->env = env;
}

` + namespace + `::C` + clsName + `::~C` + clsName + `() {
	this->env->DeleteLocalRef(this->instance);
}

`

	methods := ``

	for methodName, methodMap := range constants.GetMethodsToMapInClass(allMappings, clsPath) {
		methods += sdkutils.GenerateMethodContent(clsName, methodName, methodMap, namespace)
	}

	cpp += methods

	err := ioutil.WriteFile(path+".cpp", []byte(cpp), 0644)

	if err != nil {
		fmt.Println(err)
		return
	}
}
