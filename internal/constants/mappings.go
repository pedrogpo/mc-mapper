package constants

import "sync"

type VersionMap struct {
	Version string
	Name    string
}

type MethodsSig struct {
	Version    string
	Params     []string
	ReturnType string
}

type Map struct {
	Name        string
	ObfMappings []VersionMap
	SrgMappings []VersionMap
}

// TODO extends Map..
type MethodMap struct {
	Name        string
	clsFromName string
	MethodsSig  []MethodsSig
	ObfMappings []VersionMap
	SrgMappings []VersionMap
}

type Mappings struct {
	Classes map[string]Map
	Fields  map[string](map[string]Map)
	Methods map[string](map[string]MethodMap)
}

var classesMutex sync.Mutex
var fieldsMutex sync.Mutex
var methodsMutex sync.Mutex

// TODO: create a addMapping function to repetitive functions

func (m *Mappings) AddClass(clsName string, obfName string, clsPath string, minecraftVersion string) {
	classesMutex.Lock()
	if _, ok := m.Classes[clsPath]; !ok {
		m.Classes[clsPath] = Map{
			Name: clsName,
			ObfMappings: []VersionMap{
				{
					Version: minecraftVersion,
					Name:    obfName,
				},
			},
			SrgMappings: []VersionMap{
				{
					Version: minecraftVersion,
					Name:    clsPath,
				},
			},
		}
	} else {
		obfMappings := m.Classes[clsPath].ObfMappings
		obfMappings = append(obfMappings, VersionMap{
			Version: minecraftVersion,
			Name:    obfName,
		})

		srgMappings := m.Classes[clsPath].SrgMappings
		srgMappings = append(srgMappings, VersionMap{
			Version: minecraftVersion,
			Name:    clsPath,
		})

		m.Classes[clsPath] = Map{
			Name:        clsName,
			ObfMappings: obfMappings,
			SrgMappings: srgMappings,
		}
	}
	classesMutex.Unlock()
}

func (m *Mappings) AddField(clsName string, fieldName string, obfName string, srgName string, minecraftVersion string) {
	fieldsMutex.Lock()
	if _, ok := m.Fields[clsName]; !ok {
		m.Fields[clsName] = make(map[string]Map)
	}

	if _, ok := m.Fields[clsName][fieldName]; !ok {
		m.Fields[clsName][fieldName] = Map{
			Name: fieldName,
			ObfMappings: []VersionMap{
				{
					Version: minecraftVersion,
					Name:    obfName,
				},
			},
			// that's the try list
			SrgMappings: []VersionMap{
				{
					Version: minecraftVersion,
					Name:    srgName,
				},
			},
		}
	} else {
		obfMappings := m.Fields[clsName][fieldName].ObfMappings
		obfMappings = append(obfMappings, VersionMap{
			Version: minecraftVersion,
			Name:    obfName,
		})

		srgMappings := m.Fields[clsName][fieldName].SrgMappings
		srgMappings = append(srgMappings, VersionMap{
			Version: minecraftVersion,
			Name:    srgName,
		})

		m.Fields[clsName][fieldName] = Map{
			ObfMappings: obfMappings,
			SrgMappings: srgMappings,
		}
	}
	fieldsMutex.Unlock()
}

func (m *Mappings) AddMethod(clsName string, methodName string, obfName string, srgName string, params []string, returnType string, minecraftVersion string, clsFromPath string) {
	methodsMutex.Lock()
	if _, ok := m.Methods[clsFromPath]; !ok {
		m.Methods[clsFromPath] = make(map[string]MethodMap)
	}

	if _, ok := m.Methods[clsFromPath][methodName]; !ok {
		m.Methods[clsFromPath][methodName] = MethodMap{
			Name:        methodName,
			clsFromName: clsName,
			ObfMappings: []VersionMap{
				{
					Version: minecraftVersion,
					Name:    obfName,
				},
			},
			SrgMappings: []VersionMap{
				{
					Version: minecraftVersion,
					Name:    srgName,
				},
			},
			MethodsSig: []MethodsSig{
				{
					Version:    minecraftVersion,
					Params:     params,
					ReturnType: returnType,
				},
			},
		}
	} else {
		obfMappings := m.Methods[clsFromPath][methodName].ObfMappings
		obfMappings = append(obfMappings, VersionMap{
			Version: minecraftVersion,
			Name:    obfName,
		})

		srgMappings := m.Methods[clsFromPath][methodName].SrgMappings
		srgMappings = append(srgMappings, VersionMap{
			Version: minecraftVersion,
			Name:    srgName,
		})

		methodsSig := m.Methods[clsFromPath][methodName].MethodsSig
		methodsSig = append(methodsSig, MethodsSig{
			Version:    minecraftVersion,
			Params:     params,
			ReturnType: returnType,
		})

		m.Methods[clsFromPath][methodName] = MethodMap{
			Name:        methodName,
			clsFromName: clsName,
			ObfMappings: obfMappings,
			SrgMappings: srgMappings,
			MethodsSig:  methodsSig,
		}
	}
	methodsMutex.Unlock()
}
