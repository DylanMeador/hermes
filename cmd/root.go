package cmd

import (
	"github.com/DylanMeador/hermes/cmd/airhorn"
	"github.com/DylanMeador/hermes/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

func Cmd(s *discordgo.Session, m *discordgo.MessageCreate) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hermes",
		Short: "How about a magic trick?",
		Long:  "For my next trick, I'll make you disappear!",
	}

	args := strings.Split(m.Content, " ")
	cmd.SetArgs(args[1:])
	cmd.SetOut(reponseWriter{s, m})

	cmd.AddCommand(airhorn.Cmd())

	return cmd
}

func Execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := Cmd(s, m).ExecuteContext(discord.GenerateDiscordContext(s, m)); err != nil {
		log.Println(err)
	}
}

type reponseWriter struct {
	s *discordgo.Session
	m *discordgo.MessageCreate
}

func(rw reponseWriter) Write(p []byte) (int, error) {
	message := string(p)
	message = strings.TrimSpace(message)

	if len(message) == 0 {
		return len(p), nil
	}

	_, err := rw.s.ChannelMessageSend(rw.m.ChannelID, message)
	return len(p), err
}
