package config_utils

import (
	"fmt"
	"github.com/siper92/core-utils"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ConfContent []byte

func GetYamlConfigPath(file string) string {
	if !strings.Contains(file, ".yaml") || file == "" {
		return ""
	}

	cwd, err := os.Getwd()
	if err != nil {
		core_utils.Debug("Error getting current working directory: " + err.Error())
		cwd = "./"
	}

	local := filepath.Join(cwd, strings.Replace(file, ".yaml", ".local.yaml", 1))
	if core_utils.FileExists(local) {
		return local
	}

	dev := filepath.Join(cwd, strings.Replace(file, ".yaml", ".dev.yaml", 1))
	if core_utils.FileExists(dev) {
		return dev
	}

	demo := filepath.Join(cwd, strings.Replace(file, ".yaml", ".demo.yaml", 1))
	if core_utils.FileExists(demo) {
		return demo
	}

	prod := filepath.Join(cwd, strings.Replace(file, ".yaml", ".prod.yaml", 1))
	if core_utils.FileExists(prod) {
		return prod
	}

	base := filepath.Join(cwd, file)
	if core_utils.FileExists(base) {
		return base
	}

	return ""
}

func GetDefaultYamlConfigPath() string {
	return GetYamlConfigPath(DefaultConfigFile)
}

// Deprecated: use GetDefaultYamlConfigPath instead
func GetDefaultConfigPath() string {
	return GetYamlConfigPath(DefaultConfigFile)
}

func LoadConfig[T interface{}](conf *T, file string) (*T, error) {
	content, err := GetConfigContent(file)
	if len(content) == 0 {
		return conf, fmt.Errorf("config file is empty: %s", file)
	} else if err != nil {
		return conf, err
	}

	err = core_utils.LoadStructFromBytes(
		conf,
		content,
	)

	return conf, err
}

func LoadConfigFromString[T interface{}](conf *T, content string) (*T, error) {
	core_utils.AllowNotice()

	_content, err := prepContent([]byte(content))
	if err != nil {
		return conf, err
	}

	err = core_utils.LoadStructFromBytes(
		conf,
		_content,
	)

	core_utils.DisallowNotice()
	return conf, err
}

func GetConfigContent(file string) (ConfContent, error) {
	var err error
	var _content ConfContent

	configPath := file
	if filepath.IsAbs(configPath) == false {
		var cwd string
		cwd, err = os.Getwd()
		if err != nil {
			return _content, err
		}

		configPath = filepath.Join(cwd, file)
	}

	configPath = core_utils.AbsPath(configPath)

	core_utils.AllowNotice()
	core_utils.Notice("Loading config from " + configPath)

	yamlFileContent, err := core_utils.GetFileContent(configPath)
	if err != nil {
		return _content, err
	}

	_content, err = prepContent(yamlFileContent.Bytes())
	if err != nil {
		return _content, err
	}

	core_utils.DisallowNotice()

	return _content, nil
}

func prepContent(bytes []byte) (ConfContent, error) {
	var err error
	var _content ConfContent
	_content = bytes
	_content, err = _content.ReplaceCustomVars()
	if err != nil {
		return _content, err
	}
	_content, err = _content.ReplaceImports()
	if err != nil {
		return _content, err
	}

	return _content, nil
}

func (c ConfContent) ReplaceImports() (ConfContent, error) {
	_content := string(c)

	imports := extractImports(_content)
	for key, value := range imports {
		replaceSComment := fmt.Sprintf("#@import(%s)", key)
		replaceS := fmt.Sprintf("@import(%s)", key)
		impContent, err := value.GetImportContent()
		if err != nil {
			return c, err
		}

		_content = strings.Replace(_content, replaceSComment, impContent, -1)
		_content = strings.Replace(_content, replaceS, impContent, -1)
	}

	return ConfContent(_content), nil
}

type importDef struct {
	key   string
	value string
}

func (c importDef) GetImportContent() (string, error) {
	// TODO: implement
	panic("not implemented")
	//return "not-implemented", nil
}

func extractImports(input string) map[string]importDef {
	re := regexp.MustCompile(
		`@import\((.*?)\)`,
	)

	matches := re.FindAllStringSubmatch(input, -1)
	extractedTags := make(map[string]string)

	prefixes := []string{"http", "/", "./", "../"}
	for _, match := range matches {
		if len(match) == 2 {
			importPath := match[1]

			imported := false
			for _, prefix := range prefixes {
				if strings.Index(importPath, prefix) == 0 {
					extractedTags[importPath] = importPath
					imported = true
					break
				}
			}
			if imported {
				continue
			}

			importFile := importPath
			if !strings.Contains(importFile, ".yaml") {
				importFile += ".yaml"
			}

			extractedTags[importPath] = fmt.Sprintf("${def}/%s", importFile)

		}
	}

	result := make(map[string]importDef)
	for key, value := range extractedTags {
		result[key] = importDef{key, value}
	}

	return result
}

func (c ConfContent) ToStringLines() []string {
	_contentS := string(c)
	//_contentS = strings.Replace(_contentS, "\r", "", -1)
	//_contentS = strings.Replace(_contentS, "\t", "    ", -1)
	//_contentS = strings.Replace(_contentS, "  ", " ", -1)
	lines := strings.Split(_contentS, "\n")

	var resultLines []string
	for _, line := range lines {
		if line != "" {
			resultLines = append(resultLines, line)
		}
	}

	return resultLines
}

func (c ConfContent) ReplaceCustomVars() (ConfContent, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return c, err
	}

	normalizedCwd := strings.ReplaceAll(cwd, "\\", "/")

	_content := string(c)
	_content = strings.Replace(_content, "${cwd}", normalizedCwd, -1)

	return ConfContent(_content), nil
}
