package sdk

import (
	"fmt"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
	sdkstruct "github.com/pedrogpo/mc-auto-mapper/internal/sdk/struct"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/generics"
)

var includesFile string = "includes.hpp"

func SDKInit(allMappings constants.Mappings) {
	fmt.Println("[DEBUG] -> SDK Builder Initialized")

	// create directories
	for clsName, clsMap := range allMappings.Classes {
		if !generics.Contains(constants.ClassesToMap, clsMap.Name) {
			continue
		}

		sdkstruct.GenerateClassStruct(clsName, allMappings)
		sdkstruct.GenerateHeaderContent(clsName, allMappings, includesFile)
		sdkstruct.GenerateCppContent(clsName, allMappings)
	}

	println("[DEBUG] -> SDK Directories sucessfully created.")
}
