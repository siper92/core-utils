package core_utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ConfContent []byte

func GetConfigContent(file string) ConfContent {
	configPath := file
	if configPath[0] != '/' {
		cwd, err := os.Getwd()
		StopOnError(err)

		configPath = filepath.Join(cwd, file)
	}

	Notice("Loading config from " + configPath)

	yamlFileContent, err := GetFileContent(configPath)
	StopOnError(err)

	_content := ConfContent(yamlFileContent.Bytes())
	_content.ReplaceImports()

	return _content
}

func (c *ConfContent) ReplaceImports() {
	_content := string(*c)

	// find lines starting with "# <<< "
	for _, line := range c.ToStringLines() {
		if strings.HasPrefix(line, "#import") {
			_content = strings.Replace(_content, line+"\n", "", 1)
			if !strings.Contains(line, "@") {
				StopOnError(fmt.Errorf("import line does not contain @: %s", line))
			}

			// get import file
			_importFile := strings.Replace(line, "#import", "", 1)
			_importFile = strings.Replace(_importFile, " ", "", -1)

			// find string after @
			_importName := "@" + strings.Split(_importFile, "@")[1]

			_importFile = _importFile + ".yaml"

			_importContent, err := GetFileContent(_importFile)
			StopOnError(err)
			if len(_importContent.Bytes()) == 0 {
				StopOnError(fmt.Errorf("import file is empty: %s", _importFile))
			}

			for _, line2 := range c.ToStringLines() {
				if !strings.HasPrefix(line2, "#import") && strings.Contains(line2, _importName) {
					padding := ""
					for _, char := range line2 {
						if char == ' ' {
							padding += " "
						} else {
							break
						}
					}

					updatedLine := strings.Replace(_importContent.String(), "\n", "\n"+padding, -1)
					_content = strings.Replace(_content, line2+"\n", padding+updatedLine, 1)
				}
			}
		}
	}

	*c = ConfContent(_content)
}

func (c *ConfContent) ToStringLines() []string {
	_contentS := string(*c)
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
