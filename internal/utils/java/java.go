package java

import (
	"regexp"
)

func ExtractParamsAndReturn(signature string) (params []string, returnType string) {
	paramRegex := regexp.MustCompile(`\((.*?)\)(.+)`)
	matches := paramRegex.FindStringSubmatch(signature)

	if len(matches) == 3 {
		paramsString := matches[1]
		returnType = matches[2]

		paramTypeRegex := regexp.MustCompile(`(?:\[(?:\[)*|)([ZBCSIJFD]|L[^;]+;)`)
		paramMatches := paramTypeRegex.FindAllStringSubmatch(paramsString, -1)

		for _, match := range paramMatches {
			params = append(params, match[1])
		}
	}

	return params, returnType
}
