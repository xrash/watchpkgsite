package app

import (
	"fmt"
	"os"

	"github.com/xrash/watchpkgsite/cmd/watchpkgsite/app/root"
)

func App() {
	rootCmd := root.CreateCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
