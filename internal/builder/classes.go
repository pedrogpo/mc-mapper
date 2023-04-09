package builder

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/generics"
)

func CreateClassesFile(allMappings constants.Mappings) {
	mappingsClasses := strings.Builder{}
	mappingsClasses.WriteString("const std::map<const char*, s_class_info> mappings_classes = {")

	for clsName, clsMap := range allMappings.Classes {
		if !generics.Contains(constants.ClassesToMap, clsName) {
			continue
		}

		mappingsClasses.WriteString("{\"" + clsName + "\",")
		mappingsClasses.WriteString("{{")

		tryList := strings.Builder{}
		grouped := make(map[string][]string)
		for _, t := range clsMap.SrgMappings {
			grouped[t.Name] = append(grouped[t.Name], t.Version)
		}

		for key, value := range grouped {
			ss := strings.Builder{}
			ss.WriteString(fmt.Sprintf("{\"%s\", {", key))
			for _, v := range value {
				ss.WriteString(fmt.Sprintf("\"%s\",", v))
			}
			ss.WriteString("}},")
			tryList.WriteString(ss.String())
		}

		tryListStr := strings.TrimSuffix(tryList.String(), ",")

		mappingsClasses.WriteString(tryListStr)

		mappingsClasses.WriteString("},{")

		for _, obfMapping := range clsMap.ObfMappings {
			mappingsClasses.WriteString("{\"" + obfMapping.Version + "\", \"" + obfMapping.Name + "\"},")
		}

		mappingsClasses.WriteString("}}")
		mappingsClasses.WriteString("},")
	}

	mappingsClasses.WriteString("};")

	err := ioutil.WriteFile("out/classes.txt", []byte(mappingsClasses.String()), 0644)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Classes file created")
}

func CreateFieldsFile(allMappings constants.Mappings) {
	mappingFields := strings.Builder{}

	for clsName, clsMap := range allMappings.Fields {
		for fieldName, fieldMap := range clsMap {
			if !generics.Contains(
				constants.FieldsToMap[clsName],
				fieldName,
			) {
				continue
			}

			mappingFields.WriteString("{" + clsName + ", { \n")

			mappingFields.WriteString(fieldName + " { \n")

			for _, v := range fieldMap.SrgMappings {
				mappingFields.WriteString("{\"" + v.Version + "\",\"" + v.Name + "\"}, \n")
			}

			for _, v := range fieldMap.ObfMappings {
				mappingFields.WriteString("{\"" + v.Version + "\",\"" + v.Name + "\"}, \n")
			}

			mappingFields.WriteString("}, \n")

			mappingFields.WriteString("}, \n")
			mappingFields.WriteString("}, \n")
		}

	}

	err := ioutil.WriteFile("out/fields.txt", []byte(mappingFields.String()), 0644)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Fields file created")
}
