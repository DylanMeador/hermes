package discord

import (
	"context"
	"github.com/DylanMeador/hermes/pkg/errors"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
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

func (hc *HermesCommand) GetUserIDFromMention(mention string) (*discordgo.User, error) {
	userID := mention

	if !strings.HasPrefix(userID, "<@") || !strings.HasSuffix(userID, ">") {
		return nil, errors.CommandArgumentErr
	} else {
		userID = strings.TrimPrefix(userID, "<@")
		userID = strings.TrimPrefix(userID, "!") // ! is optional in @mention I think
		userID = strings.TrimSuffix(userID, ">")
	}

	user, err := hc.Session.User(userID)
	if err != nil {
		log.Println(err)
		return nil, errors.CommandArgumentErr
	}

	return user, nil
}