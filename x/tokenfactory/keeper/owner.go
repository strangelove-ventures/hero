package keeper

import (
	"github.com/strangelove-ventures/hero/x/tokenfactory/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetOwner set owner in the store
func (k Keeper) SetOwner(ctx sdk.Context, owner types.Owner) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&owner)
	store.Set(types.KeyPrefix(types.OwnerKey), b)
}

// GetOwner returns owner
func (k Keeper) GetOwner(ctx sdk.Context) (val types.Owner, found bool) {
	store := ctx.KVStore(k.storeKey)

	b := store.Get(types.KeyPrefix(types.OwnerKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveOwner removes owner from the store
func (k Keeper) RemoveOwner(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OwnerKey))
	store.Delete(types.KeyPrefix(types.OwnerKey))
}
