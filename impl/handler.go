package impl

import (
	// "fmt"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if pathsToUrls[path] != "" {
			header := w.Header()
			header.Set("Location", pathsToUrls[path])
			w.WriteHeader(302)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYAML, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathsToUrls := makeMap(parsedYAML)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseYAML(yml []byte) (*[]map[string]string, error) {
	parsedYAML := make([]map[string]string, 0)
	err := yaml.Unmarshal(yml, &parsedYAML)
	return &parsedYAML, err
}

func makeMap(parsed *[]map[string]string) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, v := range *parsed {
		pathsToUrls[v["path"]] = v["url"]
	}
	return pathsToUrls
}

// Similar to YAMLHandler, but parses JSON instead.
// JSON is expected to be in the format:
// [
//    {"path": "/some-path",
//	   "url": "https://www.some-url.com/demo"}
// ]
//
//
func JSONHandler(inputJson []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(inputJson)
	if err != nil {
		return nil, err
	}
	pathsToUrls := makeMap(parsedJSON)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseJSON(inputJson []byte) (*[]map[string]string, error) {
	parsedJSON := make([]map[string]string, 0)
	err := json.Unmarshal(inputJson, &parsedJSON)
	return &parsedJSON, err
}
