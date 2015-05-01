package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/brimstone/go-saverequest"
)

type dataItem struct {
	mime     string
	value    string
	children map[string]dataItem
}

var root dataItem

func list(base dataItem) ([]string, error) {
	log.Printf("base: %#s\n", base)
	keys := make([]string, 0)
	for k := range base.children {
		path := strings.Split(k, "/")
		if path[len(path)-1][0] != '\x2e' {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func returnDataItem(w http.ResponseWriter, object dataItem) error {
	w.Header().Add("Content-Type", object.mime+"; charset=utf-8")
	fmt.Fprint(w, object.value)
	return nil
}

func getDataItem(base dataItem, path []string) (dataItem, error) {
	//log.Println("Trying to get", path)
	log.Printf("looking for child: %s, %#v\n", path[0], base)
	if object, ok := base.children[path[0]]; ok {
		if len(path) > 1 {
			return getDataItem(object, path[1:])
		}
		return object, nil
	}
	return dataItem{}, fmt.Errorf("No such key")
}

func makeParents(base *dataItem, path []string) *dataItem {
	// exit conditions: len(path) < 2
	if base.children == nil {
		base.children = make(map[string]dataItem)
	}
	if len(path) == 1 {
		return base
	}
	if object, ok := base.children[path[0]]; ok {
		return &object
	} else {
		kid := makeParents(&object, path[1:])
		base.children[path[0]] = *kid
		return kid
	}
	return nil
}

func storeDataItem(path []string, object dataItem) error {
	parent := makeParents(&root, path)
	parent.children[path[len(path)-1]] = object
	return nil
}

func handleData(w http.ResponseWriter, r *http.Request) {

	//saverequest.Save(r)

	path := strings.Split(r.URL.Path, "/")
	path = path[1:]

	if r.Method == "GET" {
		if path[len(path)-1] == "" {
			//log.Println("Looking for", path, ".index")
			if object, err := getDataItem(root, append(path[:len(path)-1], ".index")); err == nil {
				returnDataItem(w, object)
				return
			} else {
				//log.Println("Didn't find .index, generating one")
				if object, err := getDataItem(root, path[:len(path)-1]); err == nil {
					//log.Println("Returning auto generated index")
					keys, _ := list(object)
					out, _ := json.Marshal(keys)
					w.Header().Add("Content-Type", "application/json; charset=utf-8")
					fmt.Fprintf(w, string(out))
					return
				} else {
					http.NotFound(w, r)
					return
				}
			}
		} else if object, err := getDataItem(root, path); err == nil {
			returnDataItem(w, object)
			return
		}
		http.NotFound(w, r)
		return
	} else if r.Method == "POST" {

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Fprint(w, "There was an error reading the body of the request")
			return
		}

		err = storeDataItem(path, dataItem{
			value: string(body),
			mime:  r.Header.Get("Content-Type"),
		})
		if err != nil {
			fmt.Fprint(w, "There was an error saving the item.", err.Error())
			return
		}

		return
	}
}

func main() {

	saverequest.WriteRequests = true
	http.HandleFunc("/data/", handleData)
	http.ListenAndServe(":8000", nil)
}
