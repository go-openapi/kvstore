package handlers

import (
	"net/http"

	"github.com/casualjim/patmosdb"
	"github.com/casualjim/patmosdb/gen/models"
	"github.com/casualjim/patmosdb/gen/restapi/operations/kv"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

func modelsError(err error) *models.Error {
	return &models.Error{
		Message: swag.String(err.Error()),
	}
}

// NewFindKeys handles a request for finding the known keys
func NewFindKeys(rt *patmosdb.Runtime) kv.FindKeysHandler {
	return &findKeys{rt: rt}
}

type findKeys struct {
	rt *patmosdb.Runtime
}

// Handle the find known keys request
func (d *findKeys) Handle(params kv.FindKeysParams) middleware.Responder {
	rid := swag.StringValue(params.XRequestID)

	values, err := d.rt.DB().FindByPrefix(swag.StringValue(params.Prefix))
	if err != nil {
		return kv.NewFindKeysDefault(http.StatusInternalServerError).WithXRequestID(rid).WithPayload(modelsError(err))
	}

	var keys []string
	for _, kva := range values {
		keys = append(keys, kva.Key)
	}
	return kv.NewFindKeysOK().WithXRequestID(rid).WithPayload(keys)
}
