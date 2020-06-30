package bugtest

import (
	"errors"
	"github.com/spf13/cobra"
)

type args struct {
}

func Cmd() *cobra.Command {
	a := &args{}

	cmd := &cobra.Command{
		Use:   "bugtest",
		Short: "generates an error to test bug responses",
		RunE:  a.run,
	}

	cmd.Hidden = true

	return cmd
}

func (a *args) run(command *cobra.Command, args []string) error {
	return errors.New("not a real error")
}

