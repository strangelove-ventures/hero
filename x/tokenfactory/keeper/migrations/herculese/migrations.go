package herculese

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey) {
	store := ctx.KVStore(storeKey)
	store.Delete([]byte("admin"))
}
