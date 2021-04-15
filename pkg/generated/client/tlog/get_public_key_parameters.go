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

package tlog

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGetPublicKeyParams creates a new GetPublicKeyParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetPublicKeyParams() *GetPublicKeyParams {
	return &GetPublicKeyParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetPublicKeyParamsWithTimeout creates a new GetPublicKeyParams object
// with the ability to set a timeout on a request.
func NewGetPublicKeyParamsWithTimeout(timeout time.Duration) *GetPublicKeyParams {
	return &GetPublicKeyParams{
		timeout: timeout,
	}
}

// NewGetPublicKeyParamsWithContext creates a new GetPublicKeyParams object
// with the ability to set a context for a request.
func NewGetPublicKeyParamsWithContext(ctx context.Context) *GetPublicKeyParams {
	return &GetPublicKeyParams{
		Context: ctx,
	}
}

// NewGetPublicKeyParamsWithHTTPClient creates a new GetPublicKeyParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetPublicKeyParamsWithHTTPClient(client *http.Client) *GetPublicKeyParams {
	return &GetPublicKeyParams{
		HTTPClient: client,
	}
}

/* GetPublicKeyParams contains all the parameters to send to the API endpoint
   for the get public key operation.

   Typically these are written to a http.Request.
*/
type GetPublicKeyParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get public key params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetPublicKeyParams) WithDefaults() *GetPublicKeyParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get public key params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetPublicKeyParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get public key params
func (o *GetPublicKeyParams) WithTimeout(timeout time.Duration) *GetPublicKeyParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get public key params
func (o *GetPublicKeyParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get public key params
func (o *GetPublicKeyParams) WithContext(ctx context.Context) *GetPublicKeyParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get public key params
func (o *GetPublicKeyParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get public key params
func (o *GetPublicKeyParams) WithHTTPClient(client *http.Client) *GetPublicKeyParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get public key params
func (o *GetPublicKeyParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *GetPublicKeyParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
