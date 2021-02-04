package vcs

import (
	"testing"
	"time"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/joiner"
	"github.com/stretchr/testify/require"
)

func TestHistoryCheckOn(t *testing.T) {
	time0 := time.Now()

	interfaceKey0 := joiner.InterfaceKey("0")

	actorKey0 := common.Key("0")

	actorKey1 := common.Key("1")

	hOld := History{
		{
			ActorKey: "",
			Key:      ActionKey("a0"),
			DoneAt:   time.Now(),
		},
		{
			ActorKey: actorKey0,
			Key:      ActionKey("a1"),
			DoneAt:   time.Now(),
		},
	}

	hNew0 := append(hOld, Action{
		ActorKey: actorKey0,
		Key:      "an0",
		DoneAt:   time0,
		Related: &joiner.Link{
			InterfaceKey: interfaceKey0,
			ID:           "123",
		},
	})

	hNew1 := append(hOld, Action{
		ActorKey: actorKey0,
		Key:      "an0",
		DoneAt:   time.Now(),
		Related: &joiner.Link{
			InterfaceKey: interfaceKey0,
			ID:           "123",
		},
	})

	err01 := hNew1.CheckOn(hNew0)
	require.Error(t, err01) // times are different

	hNew2 := append(hOld, Action{
		ActorKey: actorKey1,
		Key:      "an0",
		DoneAt:   time0,
		Related: &joiner.Link{
			InterfaceKey: interfaceKey0,
			ID:           "123",
		},
	})

	err02 := hNew2.CheckOn(hNew0)
	require.Error(t, err02) // actors are different

	hNew0duplicate := append(hOld, Action{
		ActorKey: actorKey0,
		Key:      "an0",
		DoneAt:   time0,
		Related: &joiner.Link{
			InterfaceKey: interfaceKey0,
			ID:           "123",
		},
	})

	err00duplicate := hNew0duplicate.CheckOn(hNew0)
	require.NoError(t, err00duplicate)

	hNew0duplicate1 := append(hNew0duplicate, Action{
		ActorKey: actorKey0,
		Key:      "an0",
		DoneAt:   time.Now(),
		Related: &joiner.Link{
			InterfaceKey: interfaceKey0,
			ID:           "123",
		},
	})

	err00duplicate1 := hNew0duplicate1.CheckOn(hNew0)
	require.NoError(t, err00duplicate1)

	err0duplicate0duplicate1 := hNew0duplicate1.CheckOn(hNew0duplicate)
	require.NoError(t, err0duplicate0duplicate1)

}