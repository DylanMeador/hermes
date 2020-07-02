package unmute

import (
	"github.com/DylanMeador/hermes/pkg/discord"
	"github.com/DylanMeador/hermes/pkg/emojis"
	"github.com/DylanMeador/hermes/pkg/gifs"
	"github.com/spf13/cobra"
)

type args struct {
	userID string
}

func Cmd() *cobra.Command {
	a := &args{}

	cmd := &cobra.Command{
		Use:   "unmute",
		Short: "undo a trick that was a little too cruel",
		RunE:  a.run,
	}

	cmd.PersistentFlags().StringVarP(&a.userID, "user", "u", "", "@user to unmute")

	return cmd
}

func (a *args) run(command *cobra.Command, args []string) error {
	hc := discord.GetHermesCommandFromContext(command.Context())

	userID := a.userID

	if userID == "" {
		userID = hc.Message.Author.ID
	} else {
		user, err := hc.GetUserIDFromMention(userID)
		if err != nil {
			return err
		}

		userID = user.ID

		if hc.Session.State.User.ID == userID {
			if !hc.IsHidden {
				hc.Session.MessageReactionAdd(hc.Message.ChannelID, hc.Message.ID, emojis.CURSING)
			}

			_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, gifs.ITS_A_TRAP)
			return err
		}
	}

	guild, err := hc.Session.State.Guild(hc.Message.GuildID)
	if err != nil {
		return err
	}

	isInVoiceChannel := false
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			isInVoiceChannel = true
		}
	}

	if isInVoiceChannel {
		err := discord.Unmute(hc, userID)
		if err != nil {
			return err
		}

		if !hc.IsHidden {
			return hc.Session.MessageReactionAdd(hc.Message.ChannelID, hc.Message.ID, emojis.THUMBSUP)
		}
	} else {
		if !hc.IsHidden {
			hc.Session.MessageReactionAdd(hc.Message.ChannelID, hc.Message.ID, emojis.FACEPALM)
		}
		_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, "I can only unmute someone in a voice channel silly")
	}

	return err
}
