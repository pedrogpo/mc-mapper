package builder

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
)

type VersionInfo struct {
	Version    string
	Params     []string
	ReturnType string
}

func CreateMethodsFile(allMappings constants.Mappings) {
	mappingMethods := strings.Builder{}
	mappingMethods.WriteString("const std::map<std::string, std::map<std::string, s_method>> mappings_methods = { \n")

	for clsName := range allMappings.Methods {

		lastNameSplitted := strings.Split(clsName, "/")
		lastName := lastNameSplitted[len(lastNameSplitted)-1]

		if _, hasInMethodsToMap := constants.MethodsToMap[lastName]; !hasInMethodsToMap {
			continue
		}
		mappingMethods.WriteString(fmt.Sprintf("\t{\"%s\", { \n", clsName))

		for methodName, methodMap := range constants.GetMethodsToMapInClass(allMappings, clsName) {

			mappingMethods.WriteString(fmt.Sprintf("\t\t{\"%s\", { \n", methodName))

			// parameters
			mappingMethods.WriteString("\t\t\t{ \n")

			groupedSigs := make(map[string][]string)

			for _, sig := range methodMap.MethodsSig {
				sigStr := "("
				for _, Param := range sig.Params {
					sigStr += Param
				}
				sigStr += ")"
				sigStr += sig.ReturnType

				groupedSigs[sigStr] = append(groupedSigs[sigStr], sig.Version)
			}

			for sig, versions := range groupedSigs {
				var versionStr string
				for _, version := range versions {
					versionStr += "\"" + version + "\","
				}
				versionStr = versionStr[:len(versionStr)-1]

				mappingMethods.WriteString(fmt.Sprintf("\t\t\t\ts_try_method{\"%s\", {%s}}, \n", sig, versionStr))
			}

			mappingMethods.WriteString("\t\t\t}, \n")
			mappingMethods.WriteString("\t\t\t{ \n")

			groupedSrgs := make(map[string][]string)

			for _, srg := range methodMap.SrgMappings {
				groupedSrgs[srg.Name] = append(groupedSrgs[srg.Name], srg.Version)
			}

			for srg, versions := range groupedSrgs {
				var versionStr string
				for _, version := range versions {
					versionStr += "\"" + version + "\","
				}
				versionStr = versionStr[:len(versionStr)-1]
				mappingMethods.WriteString(fmt.Sprintf("\t\t\t\ts_try_method{\"%s\", {%s}}, \n", srg, versionStr))
			}

			mappingMethods.WriteString("\t\t\t}, \n")

			mappingMethods.WriteString("\t\t\t{ \n")

			for _, v := range methodMap.ObfMappings {
				mappingMethods.WriteString(fmt.Sprintf("\t\t\t\t{\"%s\",\"%s\"}, \n", v.Version, v.Name))
			}

			mappingMethods.WriteString("\t\t\t}, \n")

			mappingMethods.WriteString("\t\t}}, \n")
		}
		mappingMethods.WriteString("\t}}, \n")
	}

	mappingMethods.WriteString("};")

	err := ioutil.WriteFile("out/methods.txt", []byte(mappingMethods.String()), 0644)

	if err != nil {
		fmt.Println(err)
		return
	}
}
