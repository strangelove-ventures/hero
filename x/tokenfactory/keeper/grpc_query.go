package keeper

import (
	"github.com/strangelove-ventures/hero/x/tokenfactory/types"
)

var _ types.QueryServer = Keeper{}
