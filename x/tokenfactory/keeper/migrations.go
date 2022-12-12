package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/hero/x/tokenfactory/keeper/migrations/herculese"

	"github.com/cosmos/cosmos-sdk/types/module"
)

var _ module.MigrationHandler = Migrator{}.MigrateToHerculese

type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		keeper: keeper,
	}
}

func (m Migrator) MigrateToHerculese(ctx sdk.Context) error {
	herculese.MigrateStore(ctx, m.keeper.storeKey)
	return nil
}
