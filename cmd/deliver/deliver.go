package deliver

import (
	"github.com/DylanMeador/hermes/pkg/discord"
	"github.com/DylanMeador/hermes/pkg/errors"
	"github.com/DylanMeador/hermes/pkg/sounds"
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
		Use:   "deliver",
		Short: "deliver a custom message to another voice channel",
		RunE:  a.run,
	}

	cmd.PersistentFlags().StringVarP(&a.channelName, "channel", "c", "", "the voice channel to play in")
	cmd.PersistentFlags().DurationVarP(&a.duration, "duration", "d", 5 * time.Second, "how long to listen before echoing")
	cmd.MarkFlagRequired("channel")

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

	recordVoiceChannelID, err := hc.GetCommandUserVoiceChannelID()
	if err != nil {
		return err
	}

	targetVoiceChannelID := ""
	for _, c := range g.Channels {
		if c.Type == discordgo.ChannelTypeGuildVoice && c.Name == a.channelName {
			targetVoiceChannelID = c.ID
			break
		}
	}

	// if you don't find exact match, check for same name with different case
	if targetVoiceChannelID == "" {
		for _, c := range g.Channels {
			if c.Type == discordgo.ChannelTypeGuildVoice && strings.EqualFold(c.Name, a.channelName) {
				targetVoiceChannelID = c.ID
				break
			}
		}
	}

	if targetVoiceChannelID == "" || recordVoiceChannelID == targetVoiceChannelID {
		return errors.CommandArgumentErr
	}

	err = discord.PlaySounds(hc, recordVoiceChannelID, nil, sounds.RECORDING_PROMPT, sounds.COUNTDOWN)
	if err != nil {
		return err
	}

	soundBytes, err := discord.RecordSounds(hc, recordVoiceChannelID, a.duration)
	if err != nil {
		return err
	}

	return discord.PlaySoundBytes(hc, targetVoiceChannelID, soundBytes)
}
