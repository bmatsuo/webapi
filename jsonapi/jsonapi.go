// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*  Filename:    jsonapi.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-08-12 17:11:47.70504 -0700 PDT
 *  Description: a json webapi implementation
 */

package jsonapi

import (
	"github.com/bmatsuo/webapi"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

var Header = map[string]string{"Content-Type": "application/json; charset=utf-8"}

var ReadError = webapi.NewError(
	http.StatusInternalServerError,
	Header,
	[]byte(`{"error":"unexpected_error"}`))

var InvalidError = webapi.NewError(
	http.StatusBadRequest,
	Header,
	[]byte(`{"error":"invalid_json"}`))

// Sets Headers in w.
func SetHeader(w http.ResponseWriter) {
	h := w.Header()
	for k, v := range Header {
		h.Set(k, v)
	}
}

// An implementation of webapi.Interface
type JsonAPI struct {
	ctor    func() interface{}
	handler http.Handler
}

// Shorthand for new(jsonapi.JsonAPI.
func New() *JsonAPI { return new(JsonAPI) }

// An http.Handler for the API.
func (api *JsonAPI) Handler() http.Handler { return webapi.Handler(api) }

// Set the API's http.Handler. Chainable.
func (api *JsonAPI) Handle(h http.Handler) *JsonAPI { api.handler = h; return api }

// Shorthand for api.Handle(http.HandlerFunc(fn)). Chainable.
func (api *JsonAPI) HandleFunc(fn func(w http.ResponseWriter, r *http.Request)) *JsonAPI {
	api.handler = http.HandlerFunc(fn)
	return api
}

// Marshal JSON POST data into the value returned by ctor.
func (api *JsonAPI) Ctor(ctor func() interface{}) *JsonAPI { api.ctor = ctor; return api }

// Implements webapi.Interace. Reads and unmarshals JSON POST data. If Ctor has
// not been given, data is unmarshaled into p. If Ctor was given, the marshaled
// value is stored in p["json"].
func (api *JsonAPI) Params(r *http.Request) (p webapi.P, e webapi.Error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		e = ReadError
		return
	}

	var v interface{}
	p = make(webapi.P)
	if v = &p; api.ctor != nil {
		v = api.ctor()
		p["json"] = v
	}
	if err := json.Unmarshal(body, v); err != nil {
		e = InvalidError
	}
	return
}

// Implements webapi.Interface. Calls and handler given with api.Handler or
// api.HandlerFunc.
func (api *JsonAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if api.handler != nil {
		api.handler.ServeHTTP(w, r)
	}
}
