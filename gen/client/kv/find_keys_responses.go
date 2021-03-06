// Code generated by go-swagger; DO NOT EDIT.

package kv

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/go-openapi/kvstore/gen/models"
)

// FindKeysReader is a Reader for the FindKeys structure.
type FindKeysReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *FindKeysReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewFindKeysOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		result := NewFindKeysDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewFindKeysOK creates a FindKeysOK with default headers values
func NewFindKeysOK() *FindKeysOK {
	return &FindKeysOK{}
}

/*FindKeysOK handles this case with default header values.

list the keys known to this datastore
*/
type FindKeysOK struct {
	/*The request id this is a response to
	 */
	XRequestID string

	Payload []string
}

func (o *FindKeysOK) Error() string {
	return fmt.Sprintf("[GET /kv][%d] findKeysOK  %+v", 200, o.Payload)
}

func (o *FindKeysOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response header X-Request-Id
	o.XRequestID = response.GetHeader("X-Request-Id")

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewFindKeysDefault creates a FindKeysDefault with default headers values
func NewFindKeysDefault(code int) *FindKeysDefault {
	return &FindKeysDefault{
		_statusCode: code,
	}
}

/*FindKeysDefault handles this case with default header values.

Error
*/
type FindKeysDefault struct {
	_statusCode int

	/*The request id this is a response to
	 */
	XRequestID string

	Payload *models.Error
}

// Code gets the status code for the find keys default response
func (o *FindKeysDefault) Code() int {
	return o._statusCode
}

func (o *FindKeysDefault) Error() string {
	return fmt.Sprintf("[GET /kv][%d] findKeys default  %+v", o._statusCode, o.Payload)
}

func (o *FindKeysDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response header X-Request-Id
	o.XRequestID = response.GetHeader("X-Request-Id")

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
