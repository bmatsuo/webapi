// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*  Filename:    params.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-08-12 17:11:47.70504 -0700 PDT
 *  Description: a type to hold api parameters
 */

/*
Params

API parameters are a set of arbitrary values indexed by strings. Parameters are
retreived for an request r with

	params := webapi.Params(r)

Params are generated by the Params method of type implementing Interface.
*/
package webapi

import (
	"net/http"
)

// See http://gorilla-web.appspot.com/pkg/context
var ContextKey = "webapi:params"

// API request parameters. Could be query, form post, json post, ...
type P map[string]interface{}

func setParams(r *http.Request, p P) { Context.Set(r, ContextKey, p) }
func deleteParams(r *http.Request)   { Context.Delete(r, ContextKey) }

// Retrieve the parameter map for a request.
func Params(r *http.Request) P {
	if p := Context.Get(r, ContextKey); p != nil {
		return p.(P)
	}
	return nil
}

// Generates parameter maps from requests.
type ParamsParser interface {
	Params(r *http.Request) (P, Error)
}

// Implements ParamsParser.
type ParamsFunc func(r *http.Request) (P, Error)

func (pfunc ParamsFunc) Params(r *http.Request) (P, Error) { return pfunc(r) }
