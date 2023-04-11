package builder

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/generics"
)

func CreateFieldsFile(allMappings constants.Mappings) {
	mappingFields := strings.Builder{}

	for clsName, clsMap := range allMappings.Fields {

		if _, ok := constants.FieldsToMap[clsName]; !ok {
			continue
		}

		mappingFields.WriteString("{\"" + clsName + "\", \n	{ \n")
		for fieldName, fieldMap := range clsMap {
			find := generics.Find(constants.FieldsToMap[clsName], func(e string) bool {
				found := false
				for _, v := range fieldMap.SrgMappings {
					if v.Name == e {
						found = true
					}
				}

				for _, v := range fieldMap.ObfMappings {
					if v.Name == e {
						found = true
					}
				}

				if e == fieldName {
					found = true
				}
				return found
			})

			if find == nil {
				continue
			}

			mappingFields.WriteString("		\"" + fieldName + "\" { \n")

			for _, v := range fieldMap.SrgMappings {
				mappingFields.WriteString("			{\"" + v.Version + "\",\"" + v.Name + "\"}, \n")
			}

			for _, v := range fieldMap.ObfMappings {
				mappingFields.WriteString("			{\"" + v.Version + "\",\"" + v.Name + "\"}, \n")
			}

			mappingFields.WriteString("		}, \n")
		}
		mappingFields.WriteString("	} \n}, \n")
	}

	err := ioutil.WriteFile("out/fields.txt", []byte(mappingFields.String()), 0644)

	if err != nil {
		fmt.Println(err)
		return
	}
}
