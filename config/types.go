package config

import "slices"

// JSON、TOML、YAML、HCL、envfile和Java properties

type FileType int

const UNSET = -1
const (
	JSON FileType = iota
	TOML
	YAML
	YML
	PROPERTIES
	PROPS
	PROP
	HCL
	TFVARS
	DOTENV
	ENV
	INI
)

var supportedFileTypes = []string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "tfvars", "dotenv", "env", "ini"}

func ParseFileType(value string) FileType {
	index := slices.Index(supportedFileTypes, value)
	if index >= 0 {
		return FileType(index)

	}
	panic("Unsupported file types")
}
func (t FileType) String() string {
	return supportedFileTypes[t]
}
