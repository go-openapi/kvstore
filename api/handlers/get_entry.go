package handlers

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/go-openapi/kvstore"
	"github.com/go-openapi/kvstore/gen/restapi/operations/kv"
	"github.com/go-openapi/kvstore/persist"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

// NewGetEntry handles a request for getting an entry
func NewGetEntry(rt *kvstore.Runtime) kv.GetEntryHandler {
	return &getEntry{rt: rt}
}

type getEntry struct {
	rt *kvstore.Runtime
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

	lastModified := time.Unix(0, value.LastUpdated).UTC().Format(time.RFC822Z)
	curVerStr := swag.StringValue(params.IfNoneMatch)
	if curVerStr != "" { // If-None-Match is optional
		curVer, err := strconv.ParseUint(curVerStr, 10, 64)
		if err != nil {
			return kv.NewGetEntryDefault(400).WithXRequestID(rid).WithPayload(modelsError(err))
		}
		if curVer == value.Version {
			return kv.NewGetEntryNotModified().WithXRequestID(rid).WithLastModified(lastModified).WithETag(curVerStr)
		}
	}

	payload := ioutil.NopCloser(bytes.NewBuffer(value.Value))
	return kv.NewGetEntryOK().WithXRequestID(rid).WithPayload(payload).WithETag(strconv.FormatUint(value.Version, 10)).WithLastModified(lastModified)
}
