package main

import (
	"flag"
	"fmt"
	"github.com/meditativeape/urlshort/impl"
	"io/ioutil"
	"net/http"
)

func main() {
	yamlPath := flag.String("yaml", "", "Path to a YAML config file")
	jsonPath := flag.String("json", "", "Path to a JSON config file")
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	handler := impl.MapHandler(pathsToUrls, mux)

	// If yaml flag is set, read the yaml file, and build a YAMLHandler using
	// the current handler as the fallback
	if len(*yamlPath) > 0 {
		yaml, err := ioutil.ReadFile(*yamlPath)
		check(err)
		handler, err = impl.YAMLHandler([]byte(yaml), handler)
		check(err)
	}

	// If json flag is set, read the json file, and build a JSONHandler using
	// the current handler as the fallback
	if len(*jsonPath) > 0 {
		json, err := ioutil.ReadFile(*jsonPath)
		check(err)
		handler, err = impl.JSONHandler([]byte(json), handler)
		check(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "No redirect for this URL!")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
