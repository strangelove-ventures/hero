package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/strangelove-ventures/hero/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateOwner_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateOwner
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateOwner{
				From: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateOwner{
				From: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
