package server_http_v2_jschmhr

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/pavlo67/common/common/auth"

	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
)

// REST ----------------------------------------------------------------------------------------------------

var _ server_http.WrapperHTTP = WrapperHTTPREST

func WrapperHTTPREST(serverOpV2 server_http.OperatorV2, serverPath string, data interface{}) (string, string, server_http.HandlerHTTP, error) {
	var ep *server_http.Endpoint

	switch v := data.(type) {
	case server_http.Endpoint:
		ep = &v
	case *server_http.Endpoint:
		ep = v
	}

	if ep == nil {
		return "", "", nil, fmt.Errorf("wrong data for WrapperHTTPREST: %#v", data)
	}

	handler := func(w http.ResponseWriter, r *http.Request, paramsHR httprouter.Params) {
		//options, err := s.onRequest.Identity(r)
		//if err != nil {
		//	l.Error(err)
		//}
		var identity *auth.Identity

		var params server_http.PathParams
		if len(paramsHR) > 0 {
			params = server_http.PathParams{}
			for _, p := range paramsHR {
				params[p.Key] = p.Value
			}
		}

		w.Header().Set("Access-Control-Allow-Origin", server_http.CORSAllowOrigin)
		w.Header().Set("Access-Control-Allow-Headers", server_http.CORSAllowHeaders)
		w.Header().Set("Access-Control-Allow-Methods", server_http.CORSAllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", server_http.CORSAllowCredentials)

		responseData, err := ep.WorkerHTTP(serverOpV2, r, params, identity)
		if err != nil {
			l.Error(err)
		}

		if responseData.MIMEType != "" {
			w.Header().Set("Content-Type", responseData.MIMEType)
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(responseData.Data)))
		if responseData.FileName != "" {
			w.Header().Set("Content-Disposition", "attachment; filename="+responseData.FileName)
		}

		if responseData.Status > 0 {
			w.WriteHeader(responseData.Status)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if _, err := w.Write(responseData.Data); err != nil {
			l.Error("can't write response", err)
		}
	}

	method := strings.ToUpper(ep.Method)
	path := ep.PathTemplate(serverPath)

	return method, path, handler, nil
}