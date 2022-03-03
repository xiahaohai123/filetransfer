package util

import "strings"

func ExtractUrlParam(url string) map[string]string {
	index := strings.Index(url, "?")
	if index == -1 {
		return map[string]string{}
	}
	paramUrl := url[index+1:]
	params := strings.Split(paramUrl, "&")
	paramMap := map[string]string{}
	for _, param := range params {
		keyValue := strings.Split(param, "=")
		if len(keyValue) == 2 {
			paramMap[keyValue[0]] = keyValue[1]
		}
	}
	return paramMap
}
