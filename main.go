package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/brimstone/go-saverequest"
)

type dataItem struct {
	mime  string
	value string
}

var data map[string]dataItem

func list() ([]string, error) {
	keys := make([]string, 0)
	for k := range data {
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

func getDataItem(path []string) (dataItem, error) {
	if object, ok := data[strings.Join(path, "/")]; ok {
		return object, nil
	}
	return dataItem{}, fmt.Errorf("No such key")
}

func handleData(w http.ResponseWriter, r *http.Request) {

	//saverequest.Save(r)

	path := strings.Split(r.URL.Path, "/")
	path = path[1:]

	if r.Method == "GET" {
		if path[len(path)-1] == "" {
			if object, err := getDataItem(append(path[:len(path)-1], ".index")); err == nil {
				returnDataItem(w, object)
				return
			} else {
				keys, _ := list()
				out, _ := json.Marshal(keys)
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				fmt.Fprintf(w, string(out))
				return
			}
		} else if object, err := getDataItem(path); err == nil {
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

		data[strings.Join(path, "/")] = dataItem{
			value: string(body),
			mime:  r.Header.Get("Content-Type"),
		}
		return
	}
}

func init() {
	data = make(map[string]dataItem)
}

func main() {

	saverequest.WriteRequests = true
	http.HandleFunc("/data/", handleData)
	http.ListenAndServe(":8000", nil)
}
