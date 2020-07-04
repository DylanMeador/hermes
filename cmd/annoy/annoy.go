package annoy

import (
	"github.com/DylanMeador/hermes/pkg/discord"
	"github.com/DylanMeador/hermes/pkg/errors"
	"github.com/DylanMeador/hermes/pkg/sounds"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"math/rand"
	"strings"
)

type args struct {
	channelName string
}

func Cmd() *cobra.Command {
	a := &args{}

	cmd := &cobra.Command{
		Use:   "annoy",
		Short: "this will be fun",
		RunE:  a.run,
	}

	cmd.PersistentFlags().StringVarP(&a.channelName, "channel", "c", "", "the voice channel to annoy")

	return cmd
}

func (a *args) run(command *cobra.Command, args []string) error {
	hc := discord.GetHermesCommandFromContext(command.Context())

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
		return discord.PlaySound(hc, voiceChannelID, sounds.ALL_SHACO[rand.Intn(len(sounds.ALL_SHACO))])
	}

	return errors.CommandArgumentErr
}
