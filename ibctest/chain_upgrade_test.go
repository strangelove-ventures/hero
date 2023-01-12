package ibctest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/icza/dyno"
	tokenfactorytypes "github.com/strangelove-ventures/hero/x/tokenfactory/types"
	"github.com/strangelove-ventures/ibctest/v3"
	"github.com/strangelove-ventures/ibctest/v3/chain/cosmos"
	"github.com/strangelove-ventures/ibctest/v3/ibc"
	"github.com/strangelove-ventures/ibctest/v3/test"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	upgradeVersion = "dan-upgrade"

	haltHeightDelta    = uint64(10) // will propose upgrade this many blocks in the future
	blocksAfterUpgrade = uint64(10)

	adminKeyName            = "admin"
	ownerKeyName            = "owner"
	masterMinterKeyName     = "masterminter"
	minterKeyName           = "minter"
	minterControllerKeyName = "mintercontroller"
	blacklisterKeyName      = "blacklister"
	pauserKeyName           = "pauser"
	userKeyName             = "user"
	user2KeyName            = "user2"
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

type AdminModuleState struct {
	Admins []string `json:"admins"`
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

func TestHeroUpgrade(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	ctx := context.Background()

	client, network := ibctest.DockerSetup(t)

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
				Repository: "ghcr.io/strangelove-ventures/hero",
				Version:    "genesis",
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

	admin := ibctest.BuildWallet(kr, adminKeyName, chainCfg)
	masterMinter := ibctest.BuildWallet(kr, masterMinterKeyName, chainCfg)
	minter := ibctest.BuildWallet(kr, minterKeyName, chainCfg)
	owner := ibctest.BuildWallet(kr, ownerKeyName, chainCfg)
	minterController := ibctest.BuildWallet(kr, minterControllerKeyName, chainCfg)
	blacklister := ibctest.BuildWallet(kr, blacklisterKeyName, chainCfg)
	pauser := ibctest.BuildWallet(kr, pauserKeyName, chainCfg)
	user := ibctest.BuildWallet(kr, userKeyName, chainCfg)
	user2 := ibctest.BuildWallet(kr, user2KeyName, chainCfg)
	alice := ibctest.BuildWallet(kr, aliceKeyName, chainCfg)

	heroValidator := hero.Validators[0]

	err = heroValidator.RecoverKey(ctx, adminKeyName, admin.Mnemonic)
	require.NoError(t, err, "failed to restore masterminter key")

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

	err = heroValidator.RecoverKey(ctx, user2KeyName, user2.Mnemonic)
	require.NoError(t, err, "failed to restore user key")

	err = heroValidator.RecoverKey(ctx, aliceKeyName, alice.Mnemonic)
	require.NoError(t, err, "failed to restore alice key")

	err = heroValidator.InitFullNodeFiles(ctx)
	require.NoError(t, err, "failed to initialize hero validator config")

	genesisWallets := []ibc.WalletAmount{
		{
			Address: admin.Address,
			Denom:   chainCfg.Denom,
			Amount:  10_000,
		},
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
			Address: user2.Address,
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

	genBz, err = modifyGenesisHero(genBz, owner.Address, admin.Address)
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

	_, err = heroValidator.ExecTx(ctx, masterMinterKeyName,
		"tokenfactory", "configure-minter-controller", minterController.Address, minter.Address,
	)
	require.NoError(t, err, "failed to execute configure minter controller tx")

	_, err = heroValidator.ExecTx(ctx, minterControllerKeyName,
		"tokenfactory", "configure-minter", minter.Address, "1000uusdc",
	)
	require.NoError(t, err, "failed to execute configure minter tx")

	_, err = heroValidator.ExecTx(ctx, minterKeyName,
		"tokenfactory", "mint", user.Address, "100uusdc",
	)
	require.NoError(t, err, "failed to execute mint to user tx")

	userBalance, err := hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	// user address= 100
	require.Equal(t, int64(100), userBalance, "failed to mint uusdc to user")

	_, err = heroValidator.ExecTx(ctx, ownerKeyName,
		"tokenfactory", "update-blacklister", blacklister.Address,
	)
	require.NoError(t, err, "failed to set blacklister")

	_, err = heroValidator.ExecTx(ctx, blacklisterKeyName,
		"tokenfactory", "blacklist", user.Address,
	)
	require.NoError(t, err, "failed to blacklist user address")

	_, err = heroValidator.ExecTx(ctx, minterKeyName,
		"tokenfactory", "mint", user.Address, "100uusdc",
	)
	require.NoError(t, err, "failed to execute mint to user tx")

	userBalance, err = hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	// user address= 100
	require.Equal(t, int64(100), userBalance, "user balance should not have incremented while blacklisted")

	_, err = heroValidator.ExecTx(ctx, minterKeyName,
		"tokenfactory", "mint", user2.Address, "100uusdc",
	)
	require.NoError(t, err, "failed to execute mint to user2 tx")

	// BEFORE THE CHAIN UPGRADE THIS TX SHOULD GO THROUGH (Sending funds from non-blacklisted user to a blacklisted user)
	err = heroValidator.SendFunds(ctx, user2KeyName, ibc.WalletAmount{
		Address: user.Address,
		Denom:   "uusdc",
		Amount:  50,
	})
	require.NoError(t, err, "Because this is before the chain upgade, The tx to a blacklisted user SHOULD have been successful")

	userBalance, err = hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	// user address= 150
	require.Equal(t, int64(150), userBalance, "Because this is before the chain upgade, user balance should have incremented while blacklisted")

	err = heroValidator.SendFunds(ctx, user2KeyName, ibc.WalletAmount{
		Address: user.Address,
		Denom:   "token",
		Amount:  100,
	})
	require.NoError(t, err, "The tx should have been successfull as that is not the minting denom")

	userBalance, err = hero.GetBalance(ctx, user.Address, "token")
	require.NoError(t, err, "failed to get user balance")

	require.Equal(t, int64(10_100), userBalance, "user balance should have incremented")

	_, err = heroValidator.ExecTx(ctx, blacklisterKeyName,
		"tokenfactory", "unblacklist", user.Address,
	)
	require.NoError(t, err, "failed to unblacklist user address")

	_, err = heroValidator.ExecTx(ctx, minterKeyName,
		"tokenfactory", "mint", user.Address, "100uusdc",
	)
	require.NoError(t, err, "failed to execute mint to user tx")

	userBalance, err = hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	// user address= 250
	require.Equal(t, int64(250), userBalance, "user balance should have increased now that they are no longer blacklisted")

	_, err = heroValidator.ExecTx(ctx, ownerKeyName,
		"tokenfactory", "update-pauser", pauser.Address,
	)
	require.NoError(t, err, "failed to update pauser")

	_, err = heroValidator.ExecTx(ctx, pauserKeyName,
		"tokenfactory", "pause",
	)
	require.NoError(t, err, "failed to pause mints")

	_, err = heroValidator.ExecTx(ctx, minterKeyName,
		"tokenfactory", "mint", user.Address, "100uusdc",
	)
	require.NoError(t, err, "failed to execute mint to user tx")

	userBalance, err = hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	// user address= 250
	require.Equal(t, int64(250), userBalance, "user balance should not have increased while chain is paused")

	_, err = heroValidator.ExecTx(ctx, userKeyName,
		"bank", "send", user.Address, alice.Address, "100uusdc",
	)
	require.Error(t, err, "transaction was successful while chain was paused")

	userBalance, err = hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	// user address= 250
	require.Equal(t, int64(250), userBalance, "user balance should not have changed while chain is paused")

	aliceBalance, err := hero.GetBalance(ctx, alice.Address, "uusdc")
	require.NoError(t, err, "failed to get alice balance")

	require.Equal(t, int64(0), aliceBalance, "alice balance should not have increased while chain is paused")

	_, err = heroValidator.ExecTx(ctx, pauserKeyName,
		"tokenfactory", "unpause",
	)
	require.NoError(t, err, "failed to unpause mints")

	_, err = heroValidator.ExecTx(ctx, userKeyName,
		"bank", "send", user.Address, alice.Address, "100uusdc",
	)
	require.NoError(t, err, "failed to send tx bank from user to alice")

	userBalance, err = hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	// user address= 150
	require.Equal(t, int64(150), userBalance, "user balance should have changed because chain is unpaused")

	aliceBalance, err = hero.GetBalance(ctx, alice.Address, "uusdc")
	require.NoError(t, err, "failed to get alice balance")

	require.Equal(t, int64(100), aliceBalance, "alice balance should have increased because chain is not paused")

	height, err := hero.Height(ctx)
	require.NoError(t, err, "error fetching height before submit upgrade proposal")
	haltHeight := height + haltHeightDelta

	_, err = heroValidator.ExecTx(ctx, adminKeyName, "adminmodule", "submit-proposal", "software-upgrade", "herculese", "--title", "Herculese upgrade", "--description", "description here", "--upgrade-height", strconv.Itoa(int(haltHeight)), "--upgrade-info", "upgrade info")
	require.NoError(t, err, "error submiting proposal")

	timeoutCtx, timeoutCtxCancel := context.WithTimeout(ctx, time.Second*45)
	defer timeoutCtxCancel()

	height, err = hero.Height(ctx)
	require.NoError(t, err, "error fetching height before upgrade")

	// this should timeout due to chain halt at upgrade height.
	_ = test.WaitForBlocks(timeoutCtx, int(haltHeight-height)+1, hero)

	height, err = hero.Height(ctx)
	require.NoError(t, err, "error fetching height after hero should have halted")

	// make sure that hero is halted
	require.Equal(t, haltHeight, height, "height is not equal to halt height")

	// bring down nodes to prepare for upgrade
	err = hero.StopAllNodes(ctx)
	require.NoError(t, err, "error stopping node(s)")

	// upgrade version on all nodes
	hero.UpgradeVersion(ctx, client, upgradeVersion)

	// start all nodes back up.
	// validators reach consensus on first block after upgrade height
	// and hero block production resumes.
	err = hero.StartAllNodes(ctx)
	require.NoError(t, err, "error starting upgraded node(s)")

	timeoutCtx, timeoutCtxCancel = context.WithTimeout(ctx, time.Second*45)
	defer timeoutCtxCancel()

	err = test.WaitForBlocks(timeoutCtx, int(blocksAfterUpgrade), hero)
	require.NoError(t, err, "hero did not produce blocks after upgrade")

	height, err = hero.Height(ctx)
	require.NoError(t, err, "error fetching height after upgrade")

	require.GreaterOrEqual(t, height, haltHeight+blocksAfterUpgrade, "height did not increment enough after upgrade")

	err = hero.StopAllNodes(ctx)
	require.NoError(t, err, "error stopping node(s)")

	state, err := heroValidator.ExportState(ctx, int64(height))
	require.NoError(t, err, "failed to export state")
	t.Log("state: ", state)

	g := make(map[string]interface{})
	err = json.Unmarshal([]byte(state), &g)
	require.NoError(t, err, "failed to unmarshal state")

	tokenFactoryAdmin, err := dyno.Get(g, "app_state", "tokenfactory", "admin")
	require.NoError(t, err)
	tfa, ok := tokenFactoryAdmin.(TokenFactoryAddress)
	t.Log("tokenFactoryAdmin: ", tokenFactoryAdmin)
	t.Log("tfa: ", tfa)
	t.Log("ok: ", ok)

	// Validate new functionality from chain update
	_, err = heroValidator.ExecTx(ctx, blacklisterKeyName,
		"tokenfactory", "blacklist", user.Address,
	)
	require.NoError(t, err, "failed to blacklist user address")

	err = heroValidator.SendFunds(ctx, user2KeyName, ibc.WalletAmount{
		Address: user.Address,
		Denom:   "uusdc",
		Amount:  50,
	})
	require.Error(t, err, "Because this is after the chain upgade, The tx to a blacklisted user should not have been successful")

	userBalance, err = hero.GetBalance(ctx, user.Address, "uusdc")
	require.NoError(t, err, "failed to get user balance")

	require.Equal(t, int64(150), userBalance, "Because this is after the chain upgade, user balance should not have incremented while blacklisted")

	err = heroValidator.SendFunds(ctx, user2KeyName, ibc.WalletAmount{
		Address: user.Address,
		Denom:   "token",
		Amount:  100,
	})
	require.NoError(t, err, "The tx should have been successfull as that is no the minting denom")

	userBalance, err = hero.GetBalance(ctx, user.Address, "token")
	require.NoError(t, err, "failed to get user balance")

	require.Equal(t, int64(10_200), userBalance, "user balance should have incremented")
}

func modifyGenesisHero(genbz []byte, ownerAddress, adminAddress string) ([]byte, error) {
	g := make(map[string]interface{})
	if err := json.Unmarshal(genbz, &g); err != nil {
		return nil, fmt.Errorf("failed to unmarshal genesis file: %w", err)
	}
	admins := []string{adminAddress}
	if err := dyno.Set(g, admins, "app_state", "adminmodule", "admins"); err != nil {
		return nil, fmt.Errorf("failed to set admin module address in genesis json: %w", err)
	}

	if err := dyno.Set(g, TokenFactoryAddress{adminAddress}, "app_state", "tokenfactory", "admin"); err != nil {
		return nil, fmt.Errorf("failed to set admin address in genesis json: %w", err)
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
