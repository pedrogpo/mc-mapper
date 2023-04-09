package constants

type VersionMap struct {
	Version string
	Name    string
}

type Map struct {
	ObfMappings []VersionMap
	SrgMappings []VersionMap
}

type MethodsSig struct {
	Version    string
	Params     []string
	ReturnType string
}

type MethodMap struct {
	MethodsSig []MethodsSig
	Map
}
