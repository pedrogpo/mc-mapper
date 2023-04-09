package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pedrogpo/mc-auto-mapper/internal/constants"
	csvutil "github.com/pedrogpo/mc-auto-mapper/internal/utils/csv"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/generics"
	_ "github.com/pedrogpo/mc-auto-mapper/internal/utils/java"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/joined"
	"github.com/pedrogpo/mc-auto-mapper/internal/utils/threadpool"
)

var rootPath string = "data/mappings/"

type Mappings struct {
	classes map[string]constants.Map
	fields  map[string]constants.Map
	methods map[string]constants.MethodMap
}

var allMappings Mappings = Mappings{
	classes: make(map[string]constants.Map),
	fields:  make(map[string]constants.Map),
	methods: make(map[string]constants.MethodMap),
}

var classesMutex sync.Mutex

func (m *Mappings) addClass(clsName string, obfName string, clsPath string, minecraftVersion string) {
	classesMutex.Lock()
	if _, ok := m.classes[clsName]; !ok {
		m.classes[clsName] = constants.Map{
			ObfMappings: []constants.VersionMap{
				{
					Version: minecraftVersion,
					Name:    obfName,
				},
			},
			SrgMappings: []constants.VersionMap{
				{
					Version: minecraftVersion,
					Name:    clsPath,
				},
			},
		}
	} else {
		obfMappings := m.classes[clsName].ObfMappings
		obfMappings = append(obfMappings, constants.VersionMap{
			Version: minecraftVersion,
			Name:    obfName,
		})

		srgMappings := m.classes[clsName].SrgMappings
		srgMappings = append(srgMappings, constants.VersionMap{
			Version: minecraftVersion,
			Name:    clsPath,
		})

		m.classes[clsName] = constants.Map{
			ObfMappings: obfMappings,
			SrgMappings: srgMappings,
		}
	}
	classesMutex.Unlock()
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

		allMappings.addClass(clsName, obfName, clsPath, minecraftVersion)
	case joined.FIELD:
		// Get the second row from the CSV file
		// obfSplitted := strings.Split(parts[1], "/")
		// obfName := obfSplitted[1]

		srgSplitted := strings.Split(parts[2], "/")
		// clsFromName := srgSplitted[len(srgSplitted)-2]
		srgName := srgSplitted[len(srgSplitted)-1]

		fieldRow := fieldsCsv.GetRowIdx(srgName)

		if generics.IsOutOfBound(fieldRow, 1) {
			return
		}

		// fieldName := fieldRow[1]

	case joined.METHOD:
		// fmt.Println("METHOD")
	}
}

func mapForVersion(minecraftVersion string) {
	joinedFile, err1 := os.Open(rootPath + minecraftVersion + "/joined.srg")
	fieldsCsv, err2 := csvutil.NewCSV(rootPath + minecraftVersion + "/fields.csv")
	methodsCsv, err3 := csvutil.NewCSV(rootPath + minecraftVersion + "/methods.csv")

	defer joinedFile.Close()
	defer fieldsCsv.Close()
	defer methodsCsv.Close()

	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("[ERR] It was not possible to read the current version: " + minecraftVersion)
		return
	}

	scanner := bufio.NewScanner(joinedFile)

	pool := threadpool.NewThreadPool(5)
	pool.Start()

	for scanner.Scan() {
		line := scanner.Text()
		// pool.AddTask(func() {
		handleJoinedType(line, fieldsCsv, methodsCsv, minecraftVersion)
		// })
	}

	pool.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func createClassesFile() {
	mappingsClasses := strings.Builder{}
	mappingsClasses.WriteString("const std::map<const char*, s_class_info> mappings_classes = {")

	for clsName, clsMap := range allMappings.classes {
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

func main() {
	// cria o timer
	inicio := time.Now()

	fmt.Println("Starting")

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

		versionName := entry.Name()
		fmt.Printf("Mapping to current directory: %s\n", versionName)

		wg.Add(1)
		go func() {
			defer wg.Done()
			mapForVersion(versionName)
		}()
	}

	wg.Wait()

	fmt.Println("Finished")

	createClassesFile()

	// para o timer
	fim := time.Now()

	// calcula o tempo decorrido em segundos
	tempoDecorrido := fim.Sub(inicio).Seconds()

	fmt.Printf("Tempo decorrido: %.2f segundos\n", tempoDecorrido)

	fmt.Scanln()
}
