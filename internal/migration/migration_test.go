package migration

import (
	"context"
	"testing"

	"gotest.tools/v3/assert"
)

func TestLoadSaveMigrationState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	path := t.TempDir() + "/migration_state.json"

	in := MigrationState{
		Orgs: []*OrgMigrationState{
			{
				Name:               "google",
				LastRanMigrationID: -1,
			},
			{
				Name:               "microsoft",
				LastRanMigrationID: 0,
			},
			{
				Name:               "default",
				LastRanMigrationID: 1,
			},
		},
	}

	if err := SaveMigrationState(ctx, &in, path); err != nil {
		t.Fatal(err)
	}

	out, err := LoadMigrationState(ctx, path)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, in, *out)
}

func TestSaveMigrationState_willOverwriteExistingData(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	path := t.TempDir() + "/migration_state.json"

	in1 := MigrationState{
		Orgs: []*OrgMigrationState{
			{
				Name:               "google",
				LastRanMigrationID: 1,
			},
		},
	}

	if err := SaveMigrationState(ctx, &in1, path); err != nil {
		t.Fatal(err)
	}

	in2 := MigrationState{
		Orgs: []*OrgMigrationState{
			{
				Name:               "google",
				LastRanMigrationID: 2,
			},
		},
	}

	if err := SaveMigrationState(ctx, &in2, path); err != nil {
		t.Fatal(err)
	}

	out, err := LoadMigrationState(ctx, path)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, in2, *out)
}

func TestLoadMigrationState_fileDoesNotExist_returnsEmptyState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	path := t.TempDir() + "/migration_state.json"

	out, err := LoadMigrationState(ctx, path)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, *out, MigrationState{})
}
