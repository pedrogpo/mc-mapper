package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/pedrogpo/mc-auto-mapper/internal/builder"
	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
	"github.com/pedrogpo/mc-auto-mapper/internal/sdk"
	csvutil "github.com/pedrogpo/mc-auto-mapper/internal/utils/csv"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/generics"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/java"
	_ "github.com/pedrogpo/mc-auto-mapper/internal/utils/java"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/joined"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/threadpool"
)

var rootPath string = "data/mappings/"

var allMappings constants.Mappings = constants.Mappings{
	Classes: make(map[string]constants.Map),
	Fields:  make(map[string](map[string]constants.Map)),
	Methods: make(map[string](map[string]constants.MethodMap)),
}

func handleJoinedType(line string, fieldsCsv *csvutil.CSV, methodsCsv *csvutil.CSV, minecraftVersion string) {
	joinedType := joined.GetJoinedType(line)

	parts := strings.Split(line, " ")

	switch joinedType {
	case joined.CLASS:
		obfName := parts[1]
		clsPath := parts[2]
		clsSplitted := strings.Split(clsPath, "/")
		clsName := clsSplitted[len(clsSplitted)-1]

		allMappings.AddClass(clsName, obfName, clsPath, minecraftVersion)
	case joined.FIELD:
		// Get the second row from the CSV file
		obfSplitted := strings.Split(parts[1], "/")
		obfName := obfSplitted[1]

		srgSplitted := strings.Split(parts[2], "/")
		clsFromName := srgSplitted[len(srgSplitted)-2]
		srgName := srgSplitted[len(srgSplitted)-1]

		fieldRow := fieldsCsv.GetRowIdx(srgName)

		if generics.IsOutOfBound(fieldRow, 1) {
			return
		}

		fieldName := fieldRow[1]

		allMappings.AddField(clsFromName, fieldName, obfName, srgName, minecraftVersion)
	case joined.METHOD:
		obfSplitted := strings.Split(parts[1], "/")
		obfName := obfSplitted[1]

		srgSplitted := strings.Split(parts[3], "/")
		clsFromName := srgSplitted[len(srgSplitted)-2]
		srgName := srgSplitted[len(srgSplitted)-1]

		signature := parts[4]

		methodRow := methodsCsv.GetRowIdx(srgName)

		if generics.IsOutOfBound(methodRow, 1) {
			return
		}

		clsFromPathSplitted := strings.Split(parts[3], "/")
		clsFromPath := strings.Join(clsFromPathSplitted[:len(clsFromPathSplitted)-1], "/")

		methodName := methodRow[1]

		params, returnType := java.ExtractParamsAndReturn(signature)

		allMappings.AddMethod(clsFromName, methodName, obfName, srgName, params, returnType, minecraftVersion, clsFromPath)
	}
}

func mapForVersion(minecraftVersion string) {
	fmt.Printf("[DEBUG] -> Mapping to current directory: %s\n", minecraftVersion)

	joinedFile, err1 := os.Open(rootPath + minecraftVersion + "/joined.srg")
	fieldsCsv, err2 := csvutil.NewCSV(rootPath + minecraftVersion + "/fields.csv")
	methodsCsv, err3 := csvutil.NewCSV(rootPath + minecraftVersion + "/methods.csv")

	defer joinedFile.Close()
	defer fieldsCsv.Close()
	defer methodsCsv.Close()

	// that's not the best way to handle it, but we don't care so much about it
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("[ERR] It was not possible to read the current version: " + minecraftVersion)
		return
	}

	scanner := bufio.NewScanner(joinedFile)

	pool := threadpool.NewThreadPool(5)
	pool.Start()

	for scanner.Scan() {
		line := scanner.Text()
		pool.AddTask(func() {
			handleJoinedType(line, fieldsCsv, methodsCsv, minecraftVersion)
		})
	}

	pool.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("[DEBUG] -> Starting")

	file, err := os.Open(rootPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	entries, err := file.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		minecraftVersion := entry.Name()

		wg.Add(1)
		go func() {
			defer wg.Done()
			mapForVersion(minecraftVersion)
		}()
	}

	wg.Wait()

	// create out/classes.txt
	builder.CreateClassesFile(allMappings)

	// create out/fields.txt
	builder.CreateFieldsFile(allMappings)

	// create out/methods.txt
	builder.CreateMethodsFile(allMappings)

	// init sdk generator
	sdk.SDKInit(allMappings)
}
