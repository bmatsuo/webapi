// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


/*
An example JSON API server.

	curl -i \
		--header 'Content-Type: application/json' \
		--data '{"name":"foobar"}' \
		'http://localhost:8080/greeting.json'
	curl -i \
		--header 'Content-Type: application/json' \
		--data '{"X":3,"Y":4}' \
		'http://localhost:8080/add.json'
*/
package main

import (
	"github.com/bmatsuo/webapi"
	"github.com/bmatsuo/webapi/jsonapi"

	"encoding/json"
	"net/http"
)

func greeting(w http.ResponseWriter, r *http.Request) {
	params := webapi.Params(r)
	name := params["name"].(string)
	jsonapi.SetHeader(w)
	p, _ := json.Marshal(map[string]string{"result": "hello, " + name + "!"})
	w.Write(p)
}

func add(w http.ResponseWriter, r *http.Request) {
	params := webapi.Params(r)["json"].(*struct{ X, Y int })
	jsonapi.SetHeader(w)
	p, _ := json.Marshal(map[string]int{"result": params.X + params.Y})
	w.Write(p)
}

func main() {
	// Unmarshals json straight into the Params map.
	http.Handle("/greeting.json",
		jsonapi.New().
			HandleFunc(greeting).
			Handler())

	// Unmarshals json into a struct type.
	http.Handle("/add.json",
		jsonapi.New().
			Ctor(func() interface{} { return new(struct{ X, Y int }) }).
			HandleFunc(add).
			Handler())

	http.ListenAndServe(":8080", nil)
}
