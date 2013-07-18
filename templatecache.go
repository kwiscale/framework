package kwiscale

import "io/ioutil"

var cache map[string]string

func getCachedTemplate(path string) string {
	if cache == nil {
		cache = make(map[string]string)
	}

	if cache[path] == "" {
		content, _ := ioutil.ReadFile(GetConfig().Templates + "/" + path)
		cache[path] = string(content)
	}

	return cache[path]
}
