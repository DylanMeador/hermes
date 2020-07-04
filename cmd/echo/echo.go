package echo

import (
	"github.com/DylanMeador/hermes/pkg/discord"
	"github.com/DylanMeador/hermes/pkg/errors"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

type args struct {
	channelName string
	duration    time.Duration
}

func Cmd() *cobra.Command {
	a := &args{}

	cmd := &cobra.Command{
		Use:   "echo",
		Short: "is that the wind?",
		RunE:  a.run,
	}

	cmd.PersistentFlags().StringVarP(&a.channelName, "channel", "c", "", "the voice channel to echo")
	cmd.PersistentFlags().DurationVarP(&a.duration, "duration", "d", 5 * time.Second, "how long to listen before echoing")

	return cmd
}

func (a *args) run(command *cobra.Command, args []string) error {
	hc := discord.GetHermesCommandFromContext(command.Context())

	if a.duration > 30 * time.Second {
		return errors.CommandArgumentErr
	}

	g, err := hc.Session.State.Guild(hc.Message.GuildID)
	if err != nil {
		return err
	}

	voiceChannelID := ""

	if a.channelName == "" {
		voiceChannelID, err = hc.GetCommandUserVoiceChannelID()
		if err != nil {
			return err
		}
	} else {
		for _, c := range g.Channels {
			if c.Type == discordgo.ChannelTypeGuildVoice && c.Name == a.channelName {
				voiceChannelID = c.ID
				break
			}
		}

		// if you don't find exact match, check for same name with different case
		if voiceChannelID == "" {
			for _, c := range g.Channels {
				if c.Type == discordgo.ChannelTypeGuildVoice && strings.EqualFold(c.Name, a.channelName) {
					voiceChannelID = c.ID
					break
				}
			}
		}
	}

	if voiceChannelID != "" {
		soundBytes, err := discord.RecordSounds(hc, voiceChannelID, a.duration)
		if err != nil {
			return err
		}

		return discord.PlaySoundBytes(hc, voiceChannelID, soundBytes)
	}

	return errors.CommandArgumentErr
}
