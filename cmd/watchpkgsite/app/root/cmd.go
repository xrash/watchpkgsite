package root

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

type cliopts struct {
	interval time.Duration
	logger   bool
	logfile  string
	workdir  string
	addr     string
}

type rootCommand struct {
	options cliopts
}

func (c *rootCommand) run(cmd *cobra.Command, args []string) {
	os.Exit(run(c.options, args))
}

func CreateCmd() *cobra.Command {
	c := &rootCommand{}

	cmd := &cobra.Command{
		Use:   "root",
		Short: "Run pksite and watch for changes",
		Long:  `Run pkgsite and watch for changes`,
		Run:   c.run,
	}

	cmd.Flags().DurationVarP(
		&c.options.interval,
		"interval",
		"",
		time.Second*30,
		"interval to look for updates, e.g. 10s, 30s, 3m, 1h etc.",
	)

	cmd.Flags().BoolVarP(
		&c.options.logger,
		"logger",
		"",
		true,
		"",
	)

	cmd.Flags().StringVarP(
		&c.options.logfile,
		"logfile",
		"",
		"",
		"",
	)

	cmd.Flags().StringVarP(
		&c.options.workdir,
		"workdir",
		"",
		".",
		"",
	)

	cmd.Flags().StringVarP(
		&c.options.addr,
		"addr",
		"",
		":8080",
		"",
	)

	return cmd
}
