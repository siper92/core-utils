package core_utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ConfContent []byte

func GetDefaultConfigPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		Debug("Error getting current working directory: " + err.Error())
		return "./" + DefaultConfigFile
	}

	local := filepath.Join(cwd, strings.Replace(DefaultConfigFile, ".yaml", ".local.yaml", 1))
	if FileExists(local) {
		return local
	}

	dev := filepath.Join(cwd, strings.Replace(DefaultConfigFile, ".yaml", ".dev.yaml", 1))
	if FileExists(dev) {
		return dev
	}

	demo := filepath.Join(cwd, strings.Replace(DefaultConfigFile, ".yaml", ".demo.yaml", 1))
	if FileExists(demo) {
		return demo
	}

	prod := filepath.Join(cwd, strings.Replace(DefaultConfigFile, ".yaml", ".prod.yaml", 1))
	if FileExists(prod) {
		return prod
	}

	loadFiles := []string{
		".conf.local.yaml",
		".conf.prod.yaml",
		".conf.demo.yaml",
		".conf.dev.yaml",
		DefaultConfigFile,
	}

	for _, file := range loadFiles {
		path := filepath.Join(cwd, file)
		if FileExists(path) {
			return path
		}
	}

	Debug("No config file found in " + cwd)
	return filepath.Join(cwd, DefaultConfigFile)
}

func LoadConfig[T interface{}](conf *T, file string) (*T, error) {
	content, err := GetConfigContent(file)
	if len(content) == 0 {
		return conf, fmt.Errorf("config file is empty: %s", file)
	} else if err != nil {
		return conf, err
	}

	err = LoadStructFromBytes(
		conf,
		content,
	)

	return conf, err
}

func LoadConfigFromString[T interface{}](conf *T, content string) (*T, error) {
	AllowNotice()

	_content, err := prepContent([]byte(content))
	if err != nil {
		return conf, err
	}

	err = LoadStructFromBytes(
		conf,
		_content,
	)

	DisallowNotice()
	return conf, err
}

func GetConfigContent(file string) (ConfContent, error) {
	var err error
	var _content ConfContent

	configPath := file
	if configPath[0] != '/' {
		var cwd string
		cwd, err = os.Getwd()
		if err != nil {
			return _content, err
		}

		configPath = filepath.Join(cwd, file)
	}

	configPath = AbsPath(configPath)

	AllowNotice()
	Notice("Loading config from " + configPath)

	yamlFileContent, err := GetFileContent(configPath)
	if err != nil {
		return _content, err
	}

	_content, err = prepContent(yamlFileContent.Bytes())
	if err != nil {
		return _content, err
	}

	DisallowNotice()

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

	_content := string(c)
	_content = strings.Replace(_content, "${cwd}", cwd, -1)

	return ConfContent(_content), nil
}
