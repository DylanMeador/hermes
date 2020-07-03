package hype

import (
	"github.com/DylanMeador/hermes/pkg/discord"
	"github.com/DylanMeador/hermes/pkg/emojis"
	"github.com/DylanMeador/hermes/pkg/gifs"
	"github.com/DylanMeador/hermes/pkg/sounds"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

type args struct {
	mention string
}

func Cmd() *cobra.Command {
	a := &args{}

	cmd := &cobra.Command{
		Use:   "hype",
		Short: "hype ya homies",
		RunE:  a.run,
	}

	cmd.PersistentFlags().StringVarP(&a.mention, "user", "u", "", "@user that deserves the hype")

	return cmd
}

func (a *args) run(command *cobra.Command, args []string) error {
	hc := discord.GetHermesCommandFromContext(command.Context())

	user, err := hc.GetUserIDFromMention(a.mention)
	if err != nil {
		return err
	}

	if hc.Session.State.User.ID == user.ID {
		if !hc.IsHidden {
			hc.Session.MessageReactionAdd(hc.Message.ChannelID, hc.Message.ID, emojis.CURSING)
		}

		_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, gifs.BAD_JOKE)
		return err
	}

	_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, strings.Repeat(emojis.PARTYING, 15))
	if err != nil {
		return err
	}

	_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, strings.Repeat(emojis.CONFETTI, 15))
	if err != nil {
		return err
	}

	_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, gifs.TYLER1_DAB)
	if err != nil {
		return err
	}
	_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, gifs.COWBELL)
	if err != nil {
		return err
	}

	guild, err := hc.Session.State.Guild(hc.Message.GuildID)
	if err != nil {
		return err
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == user.ID {
			postSoundCb := func() error {
				time.Sleep(time.Millisecond * 100)
				return nil
			}
			discord.PlaySounds(hc, vs.ChannelID, postSoundCb, sounds.HORNHORNHORN, sounds.HYPE1, sounds.HYPE2)
		}
	}

	return err
}

