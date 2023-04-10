package sdkstruct

import (
	"fmt"
	"os"
	"strings"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
)

func GenerateClassStruct(clsPath string, allMappings constants.Mappings) {
	// Remove the file name from the path
	// path := "out/sdk/" + strings.Replace(clsPath, "net/", "", 1)

	path := "out/sdk/" + clsPath

	dir := path[:strings.LastIndex(path, "/")]

	// check if path exists
	_, errExists := os.Stat(dir)

	if !os.IsNotExist(errExists) {
		return
	}

	// create directories
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	// create cpp file inside the folder
	file, err := os.Create(path + ".cpp")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// create hpp file inside the folder
	file, err = os.Create(path + ".hpp")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
}
