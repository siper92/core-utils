package core_utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ConfContent []byte

func LoadConfig[T interface{}](conf T, file string) (T, error) {
	content, err := GetConfigContent(file)
	if len(content) == 0 {
		return conf, fmt.Errorf("config file is empty: %s", file)
	} else if err != nil {
		return conf, err
	}

	err = LoadStructFromBytes(
		&conf,
		content,
	)

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

	AllowNotice()
	Notice("Loading config from " + configPath)

	yamlFileContent, err := GetFileContent(configPath)
	if err != nil {
		return _content, err
	}

	_content = yamlFileContent.Bytes()
	_content, err = _content.ReplaceCustomVars()
	if err != nil {
		return _content, err
	}
	_content, err = _content.ReplaceImports()
	if err != nil {
		return _content, err
	}

	DisallowNotice()

	return _content, nil
}

func (c ConfContent) ReplaceImports() (ConfContent, error) {
	_content := string(c)

	// find lines starting with "# <<< "
	for _, line := range c.ToStringLines() {
		if strings.Contains(line, "@import(") {

			// get import file
			_importFile := strings.Replace(line, "#import", "", 1)
			_importFile = strings.Replace(_importFile, " ", "", -1)

			importContent := ""

			_content = strings.Replace(_content, line+"\n", importContent, 1)

		}
	}

	return ConfContent(_content), nil
}

func extractImports(input string) map[string]string {
	re := regexp.MustCompile(
		`@import\((.*?)\)`,
	)

	matches := re.FindAllStringSubmatch(input, -1)
	extractedTags := make(map[string]string)

	for _, match := range matches {
		if len(match) == 2 {
			importPath := match[1]

			if strings.Index(importPath, "/") == 0 {
				if !IsDir(importPath) {
					Warning("%s is not a directory", importPath)
					continue
				}

				extractedTags[importPath] = importPath
			} else if strings.Index(importPath, "./") == 0 ||
				strings.Index(importPath, "../") == 0 {
				extractedTags[importPath] = AbsPath(importPath)
			} else {
				importFile := importPath
				if !strings.Contains(importFile, ".yaml") {
					importFile += ".yaml"
				}

				extractedTags[importPath] = fmt.Sprintf("${default}/%s", importFile)
			}
		}
	}

	return extractedTags
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
