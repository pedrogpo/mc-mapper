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
	mappingMethods.WriteString("const std::map<std::string, std::map<std::string, s_method>> mappings_methods = { \n")

	for clsName, clsMap := range allMappings.Methods {
		if _, hasInMethodsToMap := constants.MethodsToMap[clsName]; !hasInMethodsToMap {
			continue
		}

		mappingMethods.WriteString(fmt.Sprintf("\t{\"%s\", { \n", clsName))

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
