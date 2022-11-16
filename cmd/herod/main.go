package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/strangelove-ventures/hero/app"
	"github.com/strangelove-ventures/hero/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd(
		"hero",
		"cosmos",
		app.DefaultNodeHome,
		"hero-1",
		app.ModuleBasics,
		app.New,
		// this line is used by starport scaffolding # root/arguments
	)

	rootCmd.AddCommand(cmd.AddConsumerSectionCmd(app.DefaultNodeHome))

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
