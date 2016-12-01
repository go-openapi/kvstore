package handlers

import (
	"errors"
	"io/ioutil"
	"strconv"

	"github.com/go-openapi/kvstore"
	"github.com/go-openapi/kvstore/gen/restapi/operations/kv"
	"github.com/go-openapi/kvstore/persist"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

// NewPutEntry handles a request for saving an entry
func NewPutEntry(rt *kvstore.Runtime) kv.PutEntryHandler {
	return &putEntry{rt: rt}
}

type putEntry struct {
	rt *kvstore.Runtime
}

// Handle the put entry request
func (d *putEntry) Handle(params kv.PutEntryParams) middleware.Responder {
	rid := swag.StringValue(params.XRequestID)
	key := params.Key
	var version uint64
	if params.IfMatch != nil && swag.StringValue(params.IfMatch) != "" {
		var err error
		version, err = strconv.ParseUint(swag.StringValue(params.IfMatch), 10, 64)
		if err != nil {
			return kv.NewPutEntryDefault(400).WithXRequestID(rid).WithPayload(modelsError(err))
		}
	}

	value, err := ioutil.ReadAll(params.Body)
	e := params.Body.Close()
	if err != nil {
		return kv.NewPutEntryDefault(0).WithXRequestID(rid).WithPayload(modelsError(err))
	}
	if e != nil {
		return kv.NewPutEntryDefault(0).WithXRequestID(rid).WithPayload(modelsError(e))
	}

	if err := d.rt.DB().Put(key, persist.Value{Value: value, Version: version}); err != nil {
		if err == persist.ErrVersionMismatch {
			return kv.NewPutEntryConflict().WithXRequestID(rid).WithPayload(modelsError(err))
		}
		if err == persist.ErrGone {
			return kv.NewPutEntryGone().WithXRequestID(rid).WithPayload(modelsError(errors.New("entry was deleted")))
		}
		if err == persist.ErrNotFound {
			return kv.NewPutEntryNotFound().WithXRequestID(rid).WithPayload(modelsError(err))
		}
		return kv.NewPutEntryDefault(0).WithXRequestID(rid).WithPayload(modelsError(err))
	}
	return kv.NewPutEntryNoContent().WithXRequestID(rid).WithETag(version)
}
