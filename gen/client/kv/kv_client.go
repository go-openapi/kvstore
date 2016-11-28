package kv

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new kv API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for kv API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
DeleteEntry delete entry API
*/
func (a *Client) DeleteEntry(params *DeleteEntryParams) (*DeleteEntryNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewDeleteEntryParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "deleteEntry",
		Method:             "DELETE",
		PathPattern:        "/kv/{key}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &DeleteEntryReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*DeleteEntryNoContent), nil

}

/*
FindKeys lists all the keys
*/
func (a *Client) FindKeys(params *FindKeysParams) (*FindKeysOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewFindKeysParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "findKeys",
		Method:             "GET",
		PathPattern:        "/kv",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &FindKeysReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*FindKeysOK), nil

}

/*
GetEntry get entry API
*/
func (a *Client) GetEntry(params *GetEntryParams, writer io.Writer) (*GetEntryOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetEntryParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getEntry",
		Method:             "GET",
		PathPattern:        "/kv/{key}",
		ProducesMediaTypes: []string{"application/octet-stream"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &GetEntryReader{formats: a.formats, writer: writer},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*GetEntryOK), nil

}

/*
PutEntry put entry API
*/
func (a *Client) PutEntry(params *PutEntryParams) (*PutEntryCreated, *PutEntryNoContent, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPutEntryParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "putEntry",
		Method:             "PUT",
		PathPattern:        "/kv/{key}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/octet-stream"},
		Schemes:            []string{"http"},
		Params:             params,
		Reader:             &PutEntryReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, nil, err
	}
	switch value := result.(type) {
	case *PutEntryCreated:
		return value, nil, nil
	case *PutEntryNoContent:
		return nil, value, nil
	}
	return nil, nil, nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}