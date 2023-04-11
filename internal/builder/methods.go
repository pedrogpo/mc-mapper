package builder

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/generics"
)

type VersionInfo struct {
	Version    string
	Params     []string
	ReturnType string
}

func CreateMethodsFile(allMappings constants.Mappings) {
	mappingMethods := strings.Builder{}

	for clsName, clsMap := range allMappings.Methods {

		if _, ok := constants.MethodsToMap[clsName]; !ok {
			continue
		}

		mappingMethods.WriteString("{\"" + clsName + "\", \n	{ \n")
		for methodName, methodMap := range clsMap {
			find := generics.Find(constants.MethodsToMap[clsName], func(e string) bool {
				found := false
				for _, v := range methodMap.SrgMappings {
					if v.Name == e {
						found = true
					}
				}

				for _, v := range methodMap.ObfMappings {
					if v.Name == e {
						found = true
					}
				}

				if e == methodName {
					found = true
				}
				return found
			})

			if find == nil {
				continue
			}

			mappingMethods.WriteString("		\"" + methodName + "\", { \n")

			// parameters
			mappingMethods.WriteString("		{")

			grouped := make(map[string][]string)

			versionStr := ""
			for _, sig := range methodMap.MethodsSig {
				versionStr += "\"" + sig.Version + "\", "
				sigStr := "("
				for _, Param := range sig.Params {
					sigStr += Param
				}
				sigStr += ")"
				sigStr += sig.ReturnType

				grouped[sigStr] = append(grouped[sigStr], sig.Version)
			}

			for sig, versions := range grouped {
				fmt.Println(methodName, sig, versions)
			}

			mappingMethods.WriteString("		}")

			for _, v := range methodMap.SrgMappings {
				mappingMethods.WriteString("			{\"" + v.Version + "\",\"" + v.Name + "\"}, \n")
			}

			for _, v := range methodMap.ObfMappings {
				mappingMethods.WriteString("			{\"" + v.Version + "\",\"" + v.Name + "\"}, \n")
			}

			mappingMethods.WriteString("		}, \n")
		}
		mappingMethods.WriteString("	} \n}, \n")
	}

	err := ioutil.WriteFile("out/methods.txt", []byte(mappingMethods.String()), 0644)

	if err != nil {
		fmt.Println(err)
		return
	}
}
