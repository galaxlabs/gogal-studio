package slug

import (
	"regexp"
	"strings"
)

var multiUnderscore = regexp.MustCompile(`_+`)

func FromAppName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, " ", "_")
	name = multiUnderscore.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")

	return name
}

func FromModuleName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = multiUnderscore.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")

	return name
}

func FromDocTypeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = multiUnderscore.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")

	return name
}

func DocTypeFolderPath(moduleName string, doctypeName string) string {
	moduleSlug := FromModuleName(moduleName)
	doctypeSlug := FromDocTypeName(doctypeName)

	return "modules/" + moduleSlug + "/doctype/" + doctypeSlug
}

func DocTypeJSONPath(moduleName string, doctypeName string) string {
	doctypeSlug := FromDocTypeName(doctypeName)

	return DocTypeFolderPath(moduleName, doctypeName) + "/" + doctypeSlug + ".json"
}
