package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/strangelove-ventures/hero/testutil/keeper"
	"github.com/strangelove-ventures/hero/testutil/nullify"
	"github.com/strangelove-ventures/hero/x/tokenfactory/keeper"
	"github.com/strangelove-ventures/hero/x/tokenfactory/types"
)

func createTestPauser(keeper *keeper.Keeper, ctx sdk.Context) types.Pauser {
	item := types.Pauser{}
	keeper.SetPauser(ctx, item)
	return item
}

func TestPauserGet(t *testing.T) {
	keeper, ctx := keepertest.TokenfactoryKeeper(t)
	item := createTestPauser(keeper, ctx)
	rst, found := keeper.GetPauser(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestPauserRemove(t *testing.T) {
	keeper, ctx := keepertest.TokenfactoryKeeper(t)
	createTestPauser(keeper, ctx)
	keeper.RemovePauser(ctx)
	_, found := keeper.GetPauser(ctx)
	require.False(t, found)
}
