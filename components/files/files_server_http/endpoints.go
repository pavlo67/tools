package files_server_http

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/errata"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/tools/components/files"
)

var saveEndpoint = server_http.Endpoint{
	Method:     "POST",
	PathParams: []string{"bucket_id", "path", "new_file_pattern"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		bucketID := files.BucketID(params["bucket_id"])
		path := params["path"]
		newFilePattern := params["new_file_pattern"]

		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return serverOp.ResponseRESTError(http.StatusBadRequest, errata.CommonError(err, "reading body"), req)
		}

		pathCorrected, err := filesOp.Save(bucketID, path, newFilePattern, data, options)
		if err != nil {
			return serverOp.ResponseRESTError(0, err, req)
		}

		return serverOp.ResponseRESTOk(0, pathCorrected)
	},
}

var readEndpoint = server_http.Endpoint{
	Method:     "GET",
	PathParams: []string{"bucket_id", "path"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		bucketID := files.BucketID(params["bucket_id"])
		path := params["path"]

		data, err := filesOp.Read(bucketID, path, options)
		if err != nil || data == nil {
			return serverOp.ResponseRESTError(0, err, req)
		}

		return serverOp.ResponseRESTOk(0, data)
	},
}

var removeEndpoint = server_http.Endpoint{
	Method:     "DELETE",
	PathParams: []string{"bucket_id", "path"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		bucketID := files.BucketID(params["bucket_id"])
		path := params["path"]

		if err := filesOp.Remove(bucketID, path, options); err != nil {
			return serverOp.ResponseRESTError(0, err, req)
		}

		return serverOp.ResponseRESTOk(0, nil)
	},
}

var listEndpoint = server_http.Endpoint{
	Method:     "GET",
	PathParams: []string{"bucket_id", "path", "depth"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		bucketID := files.BucketID(params["bucket_id"])
		path := params["path"]
		depth, err := strconv.Atoi(params["depth"])
		if err != nil {
			return serverOp.ResponseRESTError(0, errors.Wrapf(err, "can't read depth (%s)", params["depth"]), req)
		}

		filesInfo, err := filesOp.List(bucketID, path, depth, options)
		if err != nil {
			return serverOp.ResponseRESTError(0, err, req)
		}

		return serverOp.ResponseRESTOk(0, filesInfo)
	},
}

var statEndpoint = server_http.Endpoint{
	Method:     "GET",
	PathParams: []string{"bucket_id", "path", "depth"},
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.Params, options *crud.Options) (server.Response, error) {
		bucketID := files.BucketID(params["bucket_id"])
		path := params["path"]
		depth, err := strconv.Atoi(params["depth"])
		if err != nil {
			return serverOp.ResponseRESTError(0, errors.Wrapf(err, "can't read depth (%s)", params["depth"]), req)
		}

		fileInfo, err := filesOp.Stat(bucketID, path, depth, options)
		if err != nil {
			return serverOp.ResponseRESTError(0, err, req)
		}

		return serverOp.ResponseRESTOk(0, fileInfo)
	},
}
