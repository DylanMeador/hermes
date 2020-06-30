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