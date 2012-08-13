// Copyright 2012, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*  Filename:    error.go
 *  Author:      Bryan Matsuo <bryan.matsuo [at] gmail.com>
 *  Created:     2012-08-12 17:11:47.70504 -0700 PDT
 *  Description: api parameter error
 */

/*
Errors

Fatal errors encountered during WebAPI.Params parsing should be returned as an
Error. Errors are used to create an *http.Response instead of calling
Interface.ServeHTTP.

Errors are implemented as an interface. This allows errors to marshal response
data in a convenient way.
*/
package webapi

type Error interface {
	Code() int
	Header() map[string]string
	Body() []byte
}

type apiError struct {
	code   int
	header map[string]string
	body   []byte
}

func NewError(code int, header map[string]string, body []byte) Error {
	if header == nil {
		header = map[string]string{"Content-Type": "text/plain; charset=utf-8"}
	}
	return apiError{ code, header, body }
}

func (e apiError) Code() int                 { return e.code }
func (e apiError) Header() map[string]string { return e.header }
func (e apiError) Body() []byte              { return e.body }
