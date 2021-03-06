// Code generated by go-swagger; DO NOT EDIT.

package kv

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/go-openapi/kvstore/gen/models"
)

// GetEntryOKCode is the HTTP code returned for type GetEntryOK
const GetEntryOKCode int = 200

/*GetEntryOK entry was found

swagger:response getEntryOK
*/
type GetEntryOK struct {
	/*The version of this entry

	 */
	ETag string `json:"ETag"`
	/*The time this entry was last modified

	 */
	LastModified string `json:"Last-Modified"`
	/*The request id this is a response to

	 */
	XRequestID string `json:"X-Request-Id"`

	/*
	  In: Body
	*/
	Payload io.ReadCloser `json:"body,omitempty"`
}

// NewGetEntryOK creates GetEntryOK with default headers values
func NewGetEntryOK() *GetEntryOK {

	return &GetEntryOK{}
}

// WithETag adds the eTag to the get entry o k response
func (o *GetEntryOK) WithETag(eTag string) *GetEntryOK {
	o.ETag = eTag
	return o
}

// SetETag sets the eTag to the get entry o k response
func (o *GetEntryOK) SetETag(eTag string) {
	o.ETag = eTag
}

// WithLastModified adds the lastModified to the get entry o k response
func (o *GetEntryOK) WithLastModified(lastModified string) *GetEntryOK {
	o.LastModified = lastModified
	return o
}

// SetLastModified sets the lastModified to the get entry o k response
func (o *GetEntryOK) SetLastModified(lastModified string) {
	o.LastModified = lastModified
}

// WithXRequestID adds the xRequestId to the get entry o k response
func (o *GetEntryOK) WithXRequestID(xRequestID string) *GetEntryOK {
	o.XRequestID = xRequestID
	return o
}

// SetXRequestID sets the xRequestId to the get entry o k response
func (o *GetEntryOK) SetXRequestID(xRequestID string) {
	o.XRequestID = xRequestID
}

// WithPayload adds the payload to the get entry o k response
func (o *GetEntryOK) WithPayload(payload io.ReadCloser) *GetEntryOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get entry o k response
func (o *GetEntryOK) SetPayload(payload io.ReadCloser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEntryOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header ETag

	eTag := o.ETag
	if eTag != "" {
		rw.Header().Set("ETag", eTag)
	}

	// response header Last-Modified

	lastModified := o.LastModified
	if lastModified != "" {
		rw.Header().Set("Last-Modified", lastModified)
	}

	// response header X-Request-Id

	xRequestID := o.XRequestID
	if xRequestID != "" {
		rw.Header().Set("X-Request-Id", xRequestID)
	}

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetEntryNotModifiedCode is the HTTP code returned for type GetEntryNotModified
const GetEntryNotModifiedCode int = 304

/*GetEntryNotModified entry was found but not modified

swagger:response getEntryNotModified
*/
type GetEntryNotModified struct {
	/*The version of this entry

	 */
	ETag string `json:"ETag"`
	/*The time this entry was last modified

	 */
	LastModified string `json:"Last-Modified"`
	/*The request id this is a response to

	 */
	XRequestID string `json:"X-Request-Id"`
}

// NewGetEntryNotModified creates GetEntryNotModified with default headers values
func NewGetEntryNotModified() *GetEntryNotModified {

	return &GetEntryNotModified{}
}

// WithETag adds the eTag to the get entry not modified response
func (o *GetEntryNotModified) WithETag(eTag string) *GetEntryNotModified {
	o.ETag = eTag
	return o
}

// SetETag sets the eTag to the get entry not modified response
func (o *GetEntryNotModified) SetETag(eTag string) {
	o.ETag = eTag
}

// WithLastModified adds the lastModified to the get entry not modified response
func (o *GetEntryNotModified) WithLastModified(lastModified string) *GetEntryNotModified {
	o.LastModified = lastModified
	return o
}

// SetLastModified sets the lastModified to the get entry not modified response
func (o *GetEntryNotModified) SetLastModified(lastModified string) {
	o.LastModified = lastModified
}

// WithXRequestID adds the xRequestId to the get entry not modified response
func (o *GetEntryNotModified) WithXRequestID(xRequestID string) *GetEntryNotModified {
	o.XRequestID = xRequestID
	return o
}

// SetXRequestID sets the xRequestId to the get entry not modified response
func (o *GetEntryNotModified) SetXRequestID(xRequestID string) {
	o.XRequestID = xRequestID
}

// WriteResponse to the client
func (o *GetEntryNotModified) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header ETag

	eTag := o.ETag
	if eTag != "" {
		rw.Header().Set("ETag", eTag)
	}

	// response header Last-Modified

	lastModified := o.LastModified
	if lastModified != "" {
		rw.Header().Set("Last-Modified", lastModified)
	}

	// response header X-Request-Id

	xRequestID := o.XRequestID
	if xRequestID != "" {
		rw.Header().Set("X-Request-Id", xRequestID)
	}

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(304)
}

// GetEntryNotFoundCode is the HTTP code returned for type GetEntryNotFound
const GetEntryNotFoundCode int = 404

/*GetEntryNotFound The entry was not found

swagger:response getEntryNotFound
*/
type GetEntryNotFound struct {
	/*The request id this is a response to

	 */
	XRequestID string `json:"X-Request-Id"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetEntryNotFound creates GetEntryNotFound with default headers values
func NewGetEntryNotFound() *GetEntryNotFound {

	return &GetEntryNotFound{}
}

// WithXRequestID adds the xRequestId to the get entry not found response
func (o *GetEntryNotFound) WithXRequestID(xRequestID string) *GetEntryNotFound {
	o.XRequestID = xRequestID
	return o
}

// SetXRequestID sets the xRequestId to the get entry not found response
func (o *GetEntryNotFound) SetXRequestID(xRequestID string) {
	o.XRequestID = xRequestID
}

// WithPayload adds the payload to the get entry not found response
func (o *GetEntryNotFound) WithPayload(payload *models.Error) *GetEntryNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get entry not found response
func (o *GetEntryNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEntryNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header X-Request-Id

	xRequestID := o.XRequestID
	if xRequestID != "" {
		rw.Header().Set("X-Request-Id", xRequestID)
	}

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetEntryDefault Error

swagger:response getEntryDefault
*/
type GetEntryDefault struct {
	_statusCode int
	/*The request id this is a response to

	 */
	XRequestID string `json:"X-Request-Id"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetEntryDefault creates GetEntryDefault with default headers values
func NewGetEntryDefault(code int) *GetEntryDefault {
	if code <= 0 {
		code = 500
	}

	return &GetEntryDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get entry default response
func (o *GetEntryDefault) WithStatusCode(code int) *GetEntryDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get entry default response
func (o *GetEntryDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithXRequestID adds the xRequestId to the get entry default response
func (o *GetEntryDefault) WithXRequestID(xRequestID string) *GetEntryDefault {
	o.XRequestID = xRequestID
	return o
}

// SetXRequestID sets the xRequestId to the get entry default response
func (o *GetEntryDefault) SetXRequestID(xRequestID string) {
	o.XRequestID = xRequestID
}

// WithPayload adds the payload to the get entry default response
func (o *GetEntryDefault) WithPayload(payload *models.Error) *GetEntryDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get entry default response
func (o *GetEntryDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEntryDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header X-Request-Id

	xRequestID := o.XRequestID
	if xRequestID != "" {
		rw.Header().Set("X-Request-Id", xRequestID)
	}

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
