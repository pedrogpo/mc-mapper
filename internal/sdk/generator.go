package sdk

import (
	"fmt"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
	sdkstruct "github.com/pedrogpo/mc-auto-mapper/internal/sdk/struct"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/generics"
)

var includesFile string = "includes.hpp"

func SDKInit(allMappings constants.Mappings) {
	fmt.Println("SDK Init")

	// create directories
	for clsName, clsMap := range allMappings.Classes {
		if !generics.Contains(constants.ClassesToMap, clsName) {
			continue
		}

		for _, value := range clsMap.SrgMappings {
			sdkstruct.GenerateClassStruct(value.Name, allMappings)
			sdkstruct.GenerateHeaderContent(value.Name, allMappings, includesFile)
			sdkstruct.GenerateCppContent(value.Name, allMappings)
		}
	}

	println("[DEBUG] -> SDK Directories sucessfully created.")

}
