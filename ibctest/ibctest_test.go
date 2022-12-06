package ibctest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/icza/dyno"
	integration "github.com/strangelove-ventures/hero/ibctest"
	tokenfactorytypes "github.com/strangelove-ventures/hero/x/tokenfactory/types"
	"github.com/strangelove-ventures/ibctest/v3"
	"github.com/strangelove-ventures/ibctest/v3/chain/cosmos"
	"github.com/strangelove-ventures/ibctest/v3/ibc"
	"github.com/strangelove-ventures/ibctest/v3/test"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	ownerKeyName            = "owner"
	masterMinterKeyName     = "masterminter"
	minterKeyName           = "minter"
	minterControllerKeyName = "mintercontroller"
	blacklisterKeyName      = "blacklister"
	pauserKeyName           = "pauser"
	userKeyName             = "user"
	aliceKeyName            = "alice"

	mintingDenom = "uusdc"
)

var (
	denomMetadata = []DenomMetadata{
		{
			Display: "usdc",
			Base:    "uusdc",
			Name:    "USDC",
			Symbol:  "USDC",
			DenomUnits: []DenomUnit{
				{
					Denom: "uusdc",
					Aliases: []string{
						"microusdc",
					},
					Exponent: "0",
				},
				{
					Denom: "musdc",
					Aliases: []string{
						"milliusdc",
					},
					Exponent: "3",
				},
				{
					Denom:    "usdc",
					Exponent: "6",
				},
			},
		},
	}
)

type DenomMetadata struct {
	Display    string      `json:"display"`
	Base       string      `json:"base"`
	Name       string      `json:"name"`
	Symbol     string      `json:"symbol"`
	DenomUnits []DenomUnit `json:"denom_units"`
}

type DenomUnit struct {
	Denom    string   `json:"denom"`
	Aliases  []string `json:"aliases"`
	Exponent string   `json:"exponent"`
}

type TokenFactoryAddress struct {
	Address string `json:"address"`
}

type TokenFactoryPaused struct {
	Paused bool `json:"paused"`
}

type TokenFactoryDenom struct {
	Denom string `json:"denom"`
}

func HeroEncoding() *simappparams.EncodingConfig {
	cfg := cosmos.DefaultEncoding()

	// register custom types
	tokenfactorytypes.RegisterInterfaces(cfg.InterfaceRegistry)

	return &cfg
}

func TestHeroChain(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	ctx := context.Background()

	client, network := ibctest.DockerSetup(t)

	repo, version := integration.GetDockerImageInfo()

	chainCfg := ibc.ChainConfig{
		Type:           "cosmos",
		Name:           "hero",
		ChainID:        "hero-1",
		Bin:            "herod",
		Denom:          "token",
		Bech32Prefix:   "cosmos",
		GasPrices:      "0.0token",
		GasAdjustment:  1.1,
		TrustingPeriod: "504h",
		NoHostMount:    false,
		Images: []ibc.DockerImage{
			{
				Repository: repo,
				Version:    version,
				UidGid:     "1025:1025",
			},
		},
		EncodingConfig: HeroEncoding(),
	}

	nv := 1
	nf := 0

	cf := ibctest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*ibctest.ChainSpec{
		{
			ChainConfig:   chainCfg,
			NumValidators: &nv,
			NumFullNodes:  &nf,
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	hero := chains[0].(*cosmos.CosmosChain)

	err = hero.Initialize(ctx, t.Name(), client, network)
	require.NoError(t, err, "failed to initialize hero chain")

	kr := keyring.NewInMemory()

	masterMinter := ibctest.BuildWallet(kr, masterMinterKeyName, chainCfg)
	minter := ibctest.BuildWallet(kr, minterKeyName, chainCfg)
	owner := ibctest.BuildWallet(kr, ownerKeyName, chainCfg)
	minterController := ibctest.BuildWallet(kr, minterControllerKeyName, chainCfg)
	blacklister := ibctest.BuildWallet(kr, blacklisterKeyName, chainCfg)
	pauser := ibctest.BuildWallet(kr, pauserKeyName, chainCfg)
	user := ibctest.BuildWallet(kr, userKeyName, chainCfg)
	alice := ibctest.BuildWallet(kr, aliceKeyName, chainCfg)

	heroValidator := hero.Validators[0]

	err = heroValidator.RecoverKey(ctx, ownerKeyName, owner.Mnemonic)
	require.NoError(t, err, "failed to restore owner key")

	err = heroValidator.RecoverKey(ctx, masterMinterKeyName, masterMinter.Mnemonic)
	require.NoError(t, err, "failed to restore masterminter key")

	err = heroValidator.RecoverKey(ctx, minterControllerKeyName, minterController.Mnemonic)
	require.NoError(t, err, "failed to restore mintercontroller key")

	err = heroValidator.RecoverKey(ctx, minterKeyName, minter.Mnemonic)
	require.NoError(t, err, "failed to restore minter key")

	err = heroValidator.RecoverKey(ctx, blacklisterKeyName, blacklister.Mnemonic)
	require.NoError(t, err, "failed to restore blacklister key")

	err = heroValidator.RecoverKey(ctx, pauserKeyName, pauser.Mnemonic)
	require.NoError(t, err, "failed to restore pauser key")

	err = heroValidator.RecoverKey(ctx, userKeyName, user.Mnemonic)
	require.NoError(t, err, "failed to restore user key")

	err = heroValidator.RecoverKey(ctx, aliceKeyName, alice.Mnemonic)
	require.NoError(t, err, "failed to restore alice key")

	err = heroValidator.InitFullNodeFiles(ctx)
	require.NoError(t, err, "failed to initialize hero validator config")

	genesisWallets := []ibc.WalletAmount{
		{
			Address: owner.Address,
			Denom:   chainCfg.Denom,
			Amount:  10_000,
		},
		{
			Address: masterMinter.Address,
			Denom:   chainCfg.Denom,
			Amount:  10_000,
		},
		{
			Address: minter.Address,
			Denom:   chainCfg.Denom,
			Amount:  10_000,
		},
		{
			Address: minterController.Address,
			Denom:   chainCfg.Denom,
			Amount:  10_000,
		},
		{
			Address: blacklister.Address,
			Denom:   chainCfg.Denom,
			Amount:  10_000,
		},
		{
			Address: pauser.Address,
			Denom:   chainCfg.Denom,
			Amount:  10_000,
		},
		{
			Address: user.Address,
			Denom:   chainCfg.Denom,
			Amount:  10_000,
		},
		{
			Address: alice.Address,
			Denom:   chainCfg.Denom,
			Amount:  10_000,
		},
	}

	for _, wallet := range genesisWallets {
		err = heroValidator.AddGenesisAccount(ctx, wallet.Address, []types.Coin{types.NewCoin(wallet.Denom, types.NewIntFromUint64(uint64(wallet.Amount)))})
		require.NoError(t, err, "failed to add genesis account")
	}

	genBz, err := heroValidator.GenesisFileContent(ctx)
	require.NoError(t, err, "failed to read genesis file")

	genBz, err = modifyGenesisHero(genBz, owner.Address)
	require.NoError(t, err, "failed to modify genesis file")

	err = heroValidator.OverwriteGenesisFile(ctx, genBz)
	require.NoError(t, err, "failed to write genesis file")

	_, _, err = heroValidator.ExecBin(ctx, "add-consumer-section")
	require.NoError(t, err, "failed to add consumer section to hero validator genesis file")

	err = heroValidator.CreateNodeContainer(ctx)
	require.NoError(t, err, "failed to create hero validator container")

	err = heroValidator.StartContainer(ctx)
	require.NoError(t, err, "failed to create hero validator container")

	_, err = heroValidator.ExecTx(ctx, ownerKeyName,
		"tokenfactory", "update-master-minter", masterMinter.Address,
	)
	require.NoError(t, err, "failed to execute update master minter tx")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	_, err = heroValidator.ExecTx(ctx, masterMinterKeyName,
		"tokenfactory", "configure-minter-controller", minterController.Address, minter.Address,
	)
	require.NoError(t, err, "failed to execute configure minter controller tx")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	_, err = heroValidator.ExecTx(ctx, minterControllerKeyName,
		"tokenfactory", "configure-minter", minter.Address, "1000uusdc",
	)
	require.NoError(t, err, "failed to execute configure minter tx")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	_, err = heroValidator.ExecTx(ctx, minterKeyName,
		"tokenfactory", "mint", user.Address, "100uusdc",
	)
	require.NoError(t, err, "failed to execute mint to user tx")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	userBalance, err := hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	require.Equal(t, int64(100), userBalance, "failed to mint uusdc to user")

	_, err = heroValidator.ExecTx(ctx, ownerKeyName,
		"tokenfactory", "update-blacklister", blacklister.Address,
	)
	require.NoError(t, err, "failed to set blacklister")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	_, err = heroValidator.ExecTx(ctx, blacklisterKeyName,
		"tokenfactory", "blacklist", user.Address,
	)
	require.NoError(t, err, "failed to blacklist user address")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	_, err = heroValidator.ExecTx(ctx, minterKeyName,
		"tokenfactory", "mint", user.Address, "100uusdc",
	)
	require.NoError(t, err, "failed to execute mint to user tx")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	userBalance, err = hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	require.Equal(t, int64(100), userBalance, "user balance should not have incremented while blacklisted")

	_, err = heroValidator.ExecTx(ctx, blacklisterKeyName,
		"tokenfactory", "unblacklist", user.Address,
	)
	require.NoError(t, err, "failed to unblacklist user address")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	_, err = heroValidator.ExecTx(ctx, minterKeyName,
		"tokenfactory", "mint", user.Address, "100uusdc",
	)
	require.NoError(t, err, "failed to execute mint to user tx")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	userBalance, err = hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	require.Equal(t, int64(200), userBalance, "user balance should have increased now that they are no longer blacklisted")

	_, err = heroValidator.ExecTx(ctx, ownerKeyName,
		"tokenfactory", "update-pauser", pauser.Address,
	)
	require.NoError(t, err, "failed to update pauser")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	_, err = heroValidator.ExecTx(ctx, pauserKeyName,
		"tokenfactory", "pause",
	)
	require.NoError(t, err, "failed to pause mints")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	_, err = heroValidator.ExecTx(ctx, minterKeyName,
		"tokenfactory", "mint", user.Address, "100uusdc",
	)
	require.NoError(t, err, "failed to execute mint to user tx")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	userBalance, err = hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	require.Equal(t, int64(200), userBalance, "user balance should not have increased while chain is paused")

	_, err = heroValidator.ExecTx(ctx, userKeyName,
		"bank", "send", user.Address, alice.Address, "100uusdc",
	)
	require.Error(t, err, "transaction was successful while chain was paused")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	userBalance, err = hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	require.Equal(t, int64(200), userBalance, "user balance should not have changed while chain is paused")

	aliceBalance, err := hero.GetBalance(ctx, alice.Address, "uusdc")
	require.NoError(t, err, "failed to get alice balance")

	require.Equal(t, int64(0), aliceBalance, "alice balance should not have increased while chain is paused")

	_, err = heroValidator.ExecTx(ctx, pauserKeyName,
		"tokenfactory", "unpause",
	)
	require.NoError(t, err, "failed to unpause mints")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	_, err = heroValidator.ExecTx(ctx, userKeyName,
		"bank", "send", user.Address, alice.Address, "100uusdc",
	)
	require.NoError(t, err, "failed to send tx bank from user to alice")

	err = test.WaitForBlocks(ctx, 1, hero)
	require.NoError(t, err, "failed to wait for a block on hero chain")

	userBalance, err = hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	require.Equal(t, int64(100), userBalance, "user balance should not have changed while chain is paused")

	aliceBalance, err = hero.GetBalance(ctx, alice.Address, "uusdc")
	require.NoError(t, err, "failed to get alice balance")

	require.Equal(t, int64(100), aliceBalance, "alice balance should not have increased while chain is paused")

}

func modifyGenesisHero(genbz []byte, ownerAddress string) ([]byte, error) {
	g := make(map[string]interface{})
	if err := json.Unmarshal(genbz, &g); err != nil {
		return nil, fmt.Errorf("failed to unmarshal genesis file: %w", err)
	}
	if err := dyno.Set(g, TokenFactoryAddress{ownerAddress}, "app_state", "tokenfactory", "owner"); err != nil {
		return nil, fmt.Errorf("failed to set owner address in genesis json: %w", err)
	}
	if err := dyno.Set(g, TokenFactoryPaused{false}, "app_state", "tokenfactory", "paused"); err != nil {
		return nil, fmt.Errorf("failed to set paused in genesis json: %w", err)
	}
	if err := dyno.Set(g, TokenFactoryDenom{mintingDenom}, "app_state", "tokenfactory", "mintingDenom"); err != nil {
		return nil, fmt.Errorf("failed to set minting denom in genesis json: %w", err)
	}

	if err := dyno.Set(g, denomMetadata, "app_state", "bank", "denom_metadata"); err != nil {
		return nil, fmt.Errorf("failed to set denom metadata in genesis json: %w", err)
	}

	out, err := json.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal genesis bytes to json: %w", err)
	}
	return out, nil

}
