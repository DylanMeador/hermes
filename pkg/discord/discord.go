package discord

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

const (
	HERMESCOMMAND = "HERMESCOMMAND"
)

type HermesCommand struct {
	Session *discordgo.Session
	Message *discordgo.Message
}


func GetHermesCommandFromContext(ctx context.Context) *HermesCommand {
	return ctx.Value(HERMESCOMMAND).(*HermesCommand)
}

func GenerateDiscordContext(s *discordgo.Session, m *discordgo.MessageCreate) context.Context {
	return context.WithValue(context.Background(), HERMESCOMMAND, &HermesCommand{s, m.Message})
}

func (hc *HermesCommand) GetCommandUserVoiceChannelID() (string, error) {
	g, err := hc.Session.State.Guild(hc.Message.GuildID)
	if err != nil {
		return "", err
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == hc.Message.Author.ID {
			return vs.ChannelID, nil
		}
	}
	return "", nil
}