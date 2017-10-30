package main

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/meditativeape/urlshort/impl"
	"github.com/meditativeape/urlshort/util"
	"io/ioutil"
	"net/http"
)

func main() {
	yamlPath := flag.String("yaml", "", "Path to a YAML config file")
	jsonPath := flag.String("json", "", "Path to a JSON config file")
	useBolt := flag.Bool("bolt", false, "Whether to read from the Bolt database")
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

	// If bolt flag is set, build a BoltHandler using the current handler as
	// the fallback
	if *useBolt {
		db, err := bolt.Open("db/bolt.db", 0600, nil)
		check(err)
		defer db.Close()
		initBoltDb(db)
		handler = impl.BoltHandler(db, handler)
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

func initBoltDb(db *bolt.DB) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucketName := []byte(util.BucketName)
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			bucket, err := tx.CreateBucket(bucketName)
			check(err)
			err = bucket.Put([]byte("/gophercises"),
				[]byte("https://gophercises.com/exercises/?flash=Welcome%20to%20Gophercises%21"))
			check(err)
		}
		return nil
	})
	check(err)
}
