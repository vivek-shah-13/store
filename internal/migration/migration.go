package migration

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

const DefaultMigrationStatePath = "migration_state.json"

type MigrationState struct {
	Orgs []*OrgMigrationState
}

type OrgMigrationState struct {
	Name               string
	LastRanMigrationID int
}

func LoadMigrationState(ctx context.Context, path string) (*MigrationState, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &MigrationState{}, nil
		}

		return nil, errors.Wrapf(err, "failed to read %s", path)
	}

	var state MigrationState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal migration state data")
	}

	return &state, nil
}

func SaveMigrationState(ctx context.Context, state *MigrationState, path string) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return errors.Wrapf(err, "failed to marshal migration state")
	}

	return ioutil.WriteFile(path, data, 0644)
}
