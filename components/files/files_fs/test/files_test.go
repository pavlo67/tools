package test

import (
	"testing"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/tools/components/files/files_fs"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/files"
	"github.com/pavlo67/tools/components/files/files_scenarios"
)

func TestFilesFS(t *testing.T) {
	_, cfgService := apps.PrepareTests(t, "test_service", "../../../../apps/", "test")
	require.NotNil(t, cfgService)

	var cfg config.Access
	err := cfgService.Value("files_fs", &cfg)
	require.NoErrorf(t, err, "%#v", cfgService)

	bucketID := files.BucketID("test_bucket")
	components := []starter.Starter{
		{files_fs.Starter(), common.Map{"buckets": files_fs.Buckets{bucketID: cfg.Path}}},
	}

	joinerOp, err := starter.Run(components, cfgService, "CLI BUILD FOR TEST")
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	files_scenarios.FilesTestScenario(t, joinerOp, files.InterfaceKey, bucketID)
}