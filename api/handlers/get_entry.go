package handlers

import (
	"bytes"
	"io/ioutil"

	"github.com/casualjim/patmosdb"
	"github.com/casualjim/patmosdb/gen/restapi/operations/kv"
	"github.com/casualjim/patmosdb/persist"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

// NewGetEntry handles a request for getting an entry
func NewGetEntry(rt *patmosdb.Runtime) kv.GetEntryHandler {
	return &getEntry{rt: rt}
}

type getEntry struct {
	rt *patmosdb.Runtime
}

// Handle the get entry request
func (d *getEntry) Handle(params kv.GetEntryParams) middleware.Responder {
	rid := swag.StringValue(params.XRequestID)

	value, err := d.rt.DB().Get(params.Key)
	if err != nil {
		if err == persist.ErrNotFound {
			return kv.NewGetEntryNotFound().WithXRequestID(rid).WithPayload(modelsError(err))
		}
		return kv.NewGetEntryDefault(0).WithXRequestID(rid).WithPayload(modelsError(err))
	}

	payload := ioutil.NopCloser(bytes.NewBuffer([]byte(value)))
	return kv.NewGetEntryOK().WithXRequestID(rid).WithPayload(payload)
}
