package invite

import (
	"github.com/DylanMeador/hermes/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
)

type args struct {
	temporaryMembership bool
}

func Cmd() *cobra.Command {
	a := &args{}

	cmd := &cobra.Command{
		Use:   "invite",
		Short: "generate a server invite link",
		RunE:  a.run,
	}

	cmd.PersistentFlags().BoolVarP(&a.temporaryMembership, "temporary", "t", false, "users will only have temporary membership to the server")

	return cmd
}

func (a *args) run(command *cobra.Command, args []string) error {
	s, m := discord.GetSessionAndMessageFromContext(command.Context())

	invite := discordgo.Invite{
		Temporary: a.temporaryMembership,
	}
	i , err := s.ChannelInviteCreate(m.ChannelID, invite)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "https://discord.gg/" + i.Code)
	return err
}
