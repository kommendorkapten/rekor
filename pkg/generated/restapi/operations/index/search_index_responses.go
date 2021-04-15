// Code generated by go-swagger; DO NOT EDIT.

//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package index

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/sigstore/rekor/pkg/generated/models"
)

// SearchIndexOKCode is the HTTP code returned for type SearchIndexOK
const SearchIndexOKCode int = 200

/*SearchIndexOK Returns zero or more entry UUIDs from the transparency log based on search query

swagger:response searchIndexOK
*/
type SearchIndexOK struct {

	/*
	  In: Body
	*/
	Payload []string `json:"body,omitempty"`
}

// NewSearchIndexOK creates SearchIndexOK with default headers values
func NewSearchIndexOK() *SearchIndexOK {

	return &SearchIndexOK{}
}

// WithPayload adds the payload to the search index o k response
func (o *SearchIndexOK) WithPayload(payload []string) *SearchIndexOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search index o k response
func (o *SearchIndexOK) SetPayload(payload []string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchIndexOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]string, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// SearchIndexBadRequestCode is the HTTP code returned for type SearchIndexBadRequest
const SearchIndexBadRequestCode int = 400

/*SearchIndexBadRequest The content supplied to the server was invalid

swagger:response searchIndexBadRequest
*/
type SearchIndexBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSearchIndexBadRequest creates SearchIndexBadRequest with default headers values
func NewSearchIndexBadRequest() *SearchIndexBadRequest {

	return &SearchIndexBadRequest{}
}

// WithPayload adds the payload to the search index bad request response
func (o *SearchIndexBadRequest) WithPayload(payload *models.Error) *SearchIndexBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search index bad request response
func (o *SearchIndexBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchIndexBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*SearchIndexDefault There was an internal error in the server while processing the request

swagger:response searchIndexDefault
*/
type SearchIndexDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSearchIndexDefault creates SearchIndexDefault with default headers values
func NewSearchIndexDefault(code int) *SearchIndexDefault {
	if code <= 0 {
		code = 500
	}

	return &SearchIndexDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the search index default response
func (o *SearchIndexDefault) WithStatusCode(code int) *SearchIndexDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the search index default response
func (o *SearchIndexDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the search index default response
func (o *SearchIndexDefault) WithPayload(payload *models.Error) *SearchIndexDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the search index default response
func (o *SearchIndexDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SearchIndexDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
