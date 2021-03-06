// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*  Filename:    webapi.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-08-12 17:11:47.70504 -0700 PDT
 *  Description: Main source file in webapi
 */

package webapi

import (
	"code.google.com/p/gorilla/context"

	"net/http"
)

// See http://gorilla-web.appspot.com/pkg/context
var Context = context.DefaultContext

// An http.Handler sets Params(r) before the api's handler is called.
func Handler(api Interface) http.Handler { return handler{api} }

type handler struct{ Interface }

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, err := h.Params(r)
	if err != nil {
		if h := err.Header(); h != nil {
			for k, v := range h {
				w.Header().Set(k, v)
			}
		}
		http.Error(w, string(err.Body()), err.Code())
		return
	}
	setParams(r, p)
	defer Context.Clear(r)
	h.Interface.ServeHTTP(w, r)
}

// An interface for web server API endpoints.
type Interface interface {
	http.Handler
	ParamsParser
}

// Create an object that implements Interface. This should only be used if
// implementing Interface is not an option.
func New(pfn ParamsFunc, h http.Handler) Interface {
	return &webAPI{pfn, h}
}

type webAPI struct {
	pfunc ParamsFunc // couldn't resist
	h     http.Handler
}

func (api *webAPI) Params(r *http.Request) (P, Error)                { return api.pfunc(r) }
func (api *webAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) { api.h.ServeHTTP(w, r) }
