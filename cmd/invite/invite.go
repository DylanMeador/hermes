package invite

import (
	"github.com/DylanMeador/hermes/discord"
	"github.com/spf13/cobra"
)

type args struct {
	channelName    string
	forceDisappear bool
	forceJoke      bool
}

func Cmd() *cobra.Command {
	a := &args{}

	cmd := &cobra.Command{
		Use:   "invite",
		Short: "generate the server invite link",
		RunE:  a.run,
	}

	return cmd
}

func (a *args) run(command *cobra.Command, args []string) error {
	s, m := discord.GetSessionAndMessageFromContext(command.Context())

	g, err := s.State.Guild(m.GuildID)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, g.VanityURLCode)
	return err
}