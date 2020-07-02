package blame

import (
	"github.com/DylanMeador/hermes/pkg/discord"
	"github.com/DylanMeador/hermes/pkg/emojis"
	"github.com/DylanMeador/hermes/pkg/sounds"
	"github.com/spf13/cobra"
)

type args struct{
}

func Cmd() *cobra.Command {
	a := &args{}

	cmd := &cobra.Command{
		Use:   "blame",
		Short: "was it really your fault?",
		RunE:  a.run,
	}

	return cmd
}

func (a *args) run(command *cobra.Command, args []string) error {
	hc := discord.GetHermesCommandFromContext(command.Context())

	voiceChannelID, err := hc.GetCommandUserVoiceChannelID()
	if err != nil {
		return err
	}

	if !hc.IsHidden {
		hc.Session.MessageReactionAdd(hc.Message.ChannelID, hc.Message.ID, emojis.PRAY)
	}

	return discord.PlaySound(hc, voiceChannelID, sounds.BLAME)
}

