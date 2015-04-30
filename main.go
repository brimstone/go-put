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
		keys = append(keys, k)
	}
	return keys, nil
}

func handleData(w http.ResponseWriter, r *http.Request) {

	//saverequest.Save(r)

	path := strings.Split(r.URL.Path, "/")
	path = path[1:]

	if r.Method == "GET" {
		if string(path[len(path)-1]) == "/" {
			keys, _ := list()
			out, _ := json.Marshal(keys)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			fmt.Fprintf(w, string(out))
		} else if object, ok := data[strings.Join(path, "/")]; ok {
			w.Header().Add("Content-Type", object.mime+"; charset=utf-8")
			fmt.Fprint(w, object.value)
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

func main() {

	data = make(map[string]dataItem)
	saverequest.WriteRequests = true
	http.HandleFunc("/data/", handleData)
	http.ListenAndServe(":8000", nil)
}
