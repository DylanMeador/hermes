package discord

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

const (
	SESSION = "SESSION"
	MESSAGE = "MESSAGE"
)

func GetSessionAndMessageFromContext(ctx context.Context) (*discordgo.Session, *discordgo.MessageCreate) {
	return ctx.Value(SESSION).(*discordgo.Session), ctx.Value(MESSAGE).(*discordgo.MessageCreate)
}

func GenerateDiscordContext(s *discordgo.Session, m *discordgo.MessageCreate) context.Context {
	ctx := context.WithValue(context.Background(), SESSION, s)
	return context.WithValue(ctx, MESSAGE, m)
}