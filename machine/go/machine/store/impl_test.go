package store

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.skia.org/infra/go/testutils/unittest"
	"go.skia.org/infra/machine/go/machine"
	"go.skia.org/infra/machine/go/machineserver/config"
)

func TestMachineToStoreDescription_NoDimensions(t *testing.T) {
	unittest.SmallTest(t)
	d := machine.NewDescription()
	m := machineDescriptionToStoreDescription(d)
	assert.Equal(t, storeDescription{
		Mode:               d.Mode,
		LastUpdated:        d.LastUpdated,
		MachineDescription: d,
	}, m)
}

func TestMachineToStoreDescription_WithDimensions(t *testing.T) {
	unittest.SmallTest(t)
	d := machine.NewDescription()
	d.Dimensions[machine.DimOS] = []string{"Android"}
	d.Dimensions[machine.DimDeviceType] = []string{"sailfish"}
	d.Dimensions[machine.DimQuarantined] = []string{"Device sailfish too hot."}

	m := machineDescriptionToStoreDescription(d)
	assert.Equal(t, storeDescription{
		OS:                 []string{"Android"},
		DeviceType:         []string{"sailfish"},
		Quarantined:        []string{"Device sailfish too hot."},
		Mode:               d.Mode,
		LastUpdated:        d.LastUpdated,
		MachineDescription: d,
	}, m)
}

func setupForTest(t *testing.T) (context.Context, config.InstanceConfig) {
	require.NotEmpty(t, os.Getenv("FIRESTORE_EMULATOR_HOST"), "This test requires the firestore emulator.")
	cfg := config.InstanceConfig{
		Store: config.Store{
			Project:  "test-project",
			Instance: fmt.Sprintf("test-%s", uuid.New()),
		},
	}
	ctx := context.Background()
	return ctx, cfg
}

func TestNew(t *testing.T) {
	unittest.LargeTest(t)
	ctx, cfg := setupForTest(t)
	_, err := New(ctx, true, cfg)
	require.NoError(t, err)
}

func TestUpdate_CanUpdateEvenIfDescriptionDoesntExist(t *testing.T) {
	unittest.LargeTest(t)
	ctx, cfg := setupForTest(t)
	store, err := New(ctx, true, cfg)
	require.NoError(t, err)

	called := false
	err = store.Update(ctx, "skia-rpi2-rack2-shelf1-001", func(previous machine.Description) machine.Description {
		assert.Equal(t, machine.ModeAvailable, previous.Mode)
		ret := previous.Copy()
		ret.Mode = machine.ModeMaintenance
		called = true
		return ret
	})
	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, int64(1), store.updateCounter.Get())
	assert.Equal(t, int64(0), store.updateDataToErrorCounter.Get())

	snap, err := store.machinesCollection.Doc("skia-rpi2-rack2-shelf1-001").Get(ctx)
	require.NoError(t, err)
	var storedDescription machine.Description
	err = snap.DataTo(&storedDescription)
	require.NoError(t, err)
	assert.Equal(t, machine.ModeMaintenance, storedDescription.Mode)
	assert.NoError(t, store.firestoreClient.Close())
}

func TestUpdate_CanUpdateIfDescriptionExists(t *testing.T) {
	unittest.LargeTest(t)
	ctx, cfg := setupForTest(t)
	store, err := New(ctx, true, cfg)
	require.NoError(t, err)

	// First write a Description.
	err = store.Update(ctx, "skia-rpi2-rack2-shelf1-001", func(previous machine.Description) machine.Description {
		ret := previous.Copy()
		ret.Mode = machine.ModeMaintenance
		ret.Dimensions[machine.DimOS] = []string{"Android"}
		return ret
	})
	require.NoError(t, err)

	// Now confirm we get the Description we previously wrote on the next update.
	err = store.Update(ctx, "skia-rpi2-rack2-shelf1-001", func(previous machine.Description) machine.Description {
		assert.Equal(t, machine.ModeMaintenance, previous.Mode)
		assert.Equal(t, []string{"Android"}, previous.Dimensions["os"])
		assert.Empty(t, previous.Dimensions[machine.DimDeviceType])
		ret := previous.Copy()
		ret.Mode = machine.ModeAvailable
		return ret
	})
	require.NoError(t, err)
	assert.NoError(t, store.firestoreClient.Close())
}

func TestWatch_StartWatchBeforeMachineExists(t *testing.T) {
	unittest.LargeTest(t)
	ctx, cfg := setupForTest(t)
	store, err := New(ctx, true, cfg)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// First add the watch.
	ch := store.Watch(ctx, "skia-rpi2-rack2-shelf1-001")

	// Then create the document.
	err = store.Update(ctx, "skia-rpi2-rack2-shelf1-001", func(previous machine.Description) machine.Description {
		ret := previous.Copy()
		ret.Mode = machine.ModeMaintenance
		return ret
	})
	require.NoError(t, err)

	// Wait for first description.
	m := <-ch
	assert.Equal(t, machine.ModeMaintenance, m.Mode)
	assert.Equal(t, int64(1), store.watchReceiveSnapshotCounter.Get())
	assert.Equal(t, int64(0), store.watchDataToErrorCounter.Get())
	assert.NoError(t, store.firestoreClient.Close())
}

func TestWatch_StartWatchAfterMachineExists(t *testing.T) {
	unittest.LargeTest(t)
	ctx, cfg := setupForTest(t)
	store, err := New(ctx, true, cfg)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// First create the document.
	err = store.Update(ctx, "skia-rpi2-rack2-shelf1-001", func(previous machine.Description) machine.Description {
		ret := previous.Copy()
		ret.Mode = machine.ModeMaintenance
		ret.Dimensions[machine.DimOS] = []string{"Android"}
		ret.Annotation.Message = "Hello World!"
		return ret
	})
	require.NoError(t, err)

	// Then add the watch.
	ch := store.Watch(ctx, "skia-rpi2-rack2-shelf1-001")

	// Wait for first description.
	m := <-ch
	assert.Equal(t, machine.ModeMaintenance, m.Mode)
	assert.Equal(t, machine.SwarmingDimensions{machine.DimOS: {"Android"}}, m.Dimensions)
	assert.Equal(t, "Hello World!", m.Annotation.Message)
	assert.NoError(t, store.firestoreClient.Close())
}

func TestWatch_IsCancellable(t *testing.T) {
	unittest.LargeTest(t)
	ctx, cfg := setupForTest(t)
	store, err := New(ctx, true, cfg)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(ctx)

	// First create the document.
	err = store.Update(ctx, "skia-rpi2-rack2-shelf1-001", func(previous machine.Description) machine.Description {
		ret := previous.Copy()
		ret.Mode = machine.ModeMaintenance
		return ret
	})
	require.NoError(t, err)

	// Then add the watch.
	ch := store.Watch(ctx, "skia-rpi2-rack2-shelf1-001")

	cancel()

	// The test passes if we get past this loop since that means the channel was closed.
	for range ch {
	}
	assert.NoError(t, store.firestoreClient.Close())
}

func TestList_Success(t *testing.T) {
	unittest.LargeTest(t)
	ctx, cfg := setupForTest(t)
	store, err := New(ctx, true, cfg)
	require.NoError(t, err)

	// List on an empty collection is OK.
	descriptions, err := store.List(ctx)
	require.NoError(t, err)
	assert.Len(t, descriptions, 0)

	// Add a single description.
	err = store.Update(ctx, "skia-rpi2-rack2-shelf1-001", func(previous machine.Description) machine.Description {
		assert.Equal(t, machine.ModeAvailable, previous.Mode)
		ret := previous.Copy()
		ret.Mode = machine.ModeMaintenance
		ret.Dimensions["foo"] = []string{"bar", "baz"}
		return ret
	})
	require.NoError(t, err)

	// Confirm it appears in the list.
	descriptions, err = store.List(ctx)
	require.NoError(t, err)
	assert.Len(t, descriptions, 1)
	assert.Equal(t, machine.SwarmingDimensions{"foo": {"bar", "baz"}}, descriptions[0].Dimensions)

	// Add a second description.
	err = store.Update(ctx, "skia-rpi2-rack2-shelf1-002", func(previous machine.Description) machine.Description {
		assert.Equal(t, machine.ModeAvailable, previous.Mode)
		ret := previous.Copy()
		ret.Mode = machine.ModeMaintenance
		ret.Dimensions["foo"] = []string{"quux"}
		return ret
	})

	// Confirm they both show up in the list.
	descriptions, err = store.List(ctx)
	require.NoError(t, err)
	assert.Len(t, descriptions, 2)
}
