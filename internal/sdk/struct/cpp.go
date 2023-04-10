package sdkstruct

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
)

func GenerateCppContent(clsPath string, allMappings constants.Mappings) {
	// Remove the file name from the path
	path := "out/sdk/" + strings.Replace(clsPath, "net/", "", 1)

	clsPathParts := strings.Split(clsPath, "/")
	clsName := clsPathParts[len(clsPathParts)-1]

	// write content in .cpp file
	cpp := `#include "` + clsName + `.hpp"

`

	namespace := "sdk::"
	for i := 0; i < len(clsPathParts)-1; i++ {
		namespace += clsPathParts[i] + "::"
	}

	namespace = strings.TrimSuffix(namespace, "::")

	cpp += namespace + `::` + clsName + `(JNIEnv* env) {
	this->env = env;
}

` + namespace + `::` + clsName + `(JNIEnv* env, jobject instance) : instance(instance) {
	this->env = env;
}

` + namespace + `::~` + clsName + `() {
	this->env->DeleteLocalRef(this->instance);
}

`

	err := ioutil.WriteFile(path+".cpp", []byte(cpp), 0644)

	if err != nil {
		fmt.Println(err)
		return
	}
}
