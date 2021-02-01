package main

import (
	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/backup/apps/backup/demo_api"
)

var (
	BuildDate   = "unknown"
	BuildTag    = "unknown"
	BuildCommit = "unknown"
)

const serviceName = "demo"

func main() {
	versionOnly, envPath, cfgService, l := apps.Prepare(BuildDate, BuildTag, BuildCommit, serviceName, apps.AppsSubpathDefault)
	if versionOnly {
		return
	}

	// running starters

	label := "BACKUP/SQLITE/REST BUILD"
	joinerOp, err := starter.Run(demo_api.Components(envPath, true, false), cfgService, label)
	if err != nil {
		l.Fatal(err)
	}
	defer joinerOp.CloseAll()

	demo_api.WG.Wait()
}