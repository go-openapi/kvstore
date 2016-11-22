package handlers

import (
	"net/http"

	"github.com/casualjim/patmosdb"
	"github.com/casualjim/patmosdb/gen/models"
	"github.com/casualjim/patmosdb/gen/restapi/operations/kv"
	"github.com/casualjim/patmosdb/persist"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

// NewDeleteEntry handles a request for deleting an entry
func NewDeleteEntry(rt *patmosdb.Runtime) kv.DeleteEntryHandler {
	return &deleteEntry{rt: rt}
}

// deleteEntry handles a request for deleting an entry
type deleteEntry struct {
	rt *patmosdb.Runtime
}

// Handle the delete entry request
func (d *deleteEntry) Handle(params kv.DeleteEntryParams) middleware.Responder {
	rid := swag.StringValue(params.XRequestID)
	noContent := kv.NewDeleteEntryNoContent().WithXRequestID(rid)

	if err := d.rt.DB().Delete(params.Key); err != nil {
		if err == persist.ErrNotFound {
			return noContent
		}
		return kv.NewDeleteEntryDefault(http.StatusInternalServerError).WithXRequestID(rid).WithPayload(&models.Error{Message: swag.String(err.Error())})
	}
	return noContent
}
