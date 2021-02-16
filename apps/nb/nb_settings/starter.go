package nb_settings

import (
	"fmt"
	"io/ioutil"

	"github.com/pavlo67/common/common/filelib"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &nbStarter{}
}

var _ starter.Operator = &nbStarter{}

type nbStarter struct {
	restPrefix  string
	pagesPrefix string
}

// --------------------------------------------------------------------------

var l logger.Operator

func (ns *nbStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ns *nbStarter) Prepare(cfg *config.Config, options common.Map) error {
	var cfgStorage common.Map
	if err := cfg.Value("nb_api", &cfgStorage); err != nil {
		return errors.CommonError(err, fmt.Sprintf("in config: %#v", cfg))
	}

	ns.restPrefix = options.StringDefault("rest_prefix", "")
	ns.pagesPrefix = options.StringDefault("pages_prefix", "")

	return nil
}

func (ns *nbStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	srvOp, _ := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if srvOp == nil {
		return fmt.Errorf("no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	srvPort, isHTTPS := srvOp.Addr() // isHTTPS

	if err := restConfig.CompleteWithJoiner(joinerOp, "", srvPort, ns.restPrefix); err != nil {
		return err
	}

	//if err := server_http.InitEndpointsWithSwaggerV2(srvOp, restConfig, !isHTTPS, restStaticPath, restServerSubpath, l); err != nil {
	//	return err
	//}

	if err := server_http.InitPages(srvOp, restConfig, l); err != nil {
		return err
	}

	restStaticPath := filelib.CurrentPath() + "../rest_static/"
	restServerSubpath := "/*filepath"

	swaggerStaticPath := restStaticPath + "api-docs/"
	swaggerStaticFilePath := swaggerStaticPath + "swaggerJSON.json"

	swaggerJSON, err := restConfig.SwaggerV2(isHTTPS)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(swaggerStaticFilePath, swaggerJSON, 0644); err != nil {
		return fmt.Errorf("on ioutil.WriteFile(%s, %s, 0755): %s", swaggerStaticFilePath, swaggerJSON, err)
	}
	l.Infof("%d bytes are written into %s", len(swaggerJSON), swaggerStaticFilePath)

	if err := srvOp.HandleFiles("rest_static", ns.restPrefix+restServerSubpath, server_http.StaticPath{LocalPath: restStaticPath}); err != nil {
		return err
	}

	if err := pagesConfig.CompleteWithJoiner(joinerOp, "", srvPort, ns.pagesPrefix); err != nil {
		return err
	}
	if err := server_http.InitPages(srvOp, pagesConfig, l); err != nil {
		return err
	}
	pagesStaticPath := filelib.CurrentPath() + "../pages_static/"
	pagesStaticServerSubpath := "/static"
	if pagesStaticPath != "" {
		if err := srvOp.HandleFiles("pages_static", ns.pagesPrefix+pagesStaticServerSubpath+"/*filepath", server_http.StaticPath{LocalPath: pagesStaticPath}); err != nil {
			return err
		}
	}

	WG.Add(1)

	go func() {
		defer WG.Done()
		if err := srvOp.Start(); err != nil {
			l.Error("on srvOp.Start(): ", err)
		}
	}()

	return nil
}

// TODO!!! customize it
// if isHTTPS {
//	go http.ListenAndServe(":80", http.HandlerFunc(server_http.Redirect))
// }
