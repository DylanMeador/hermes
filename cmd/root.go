package cmd

import (
	"github.com/DylanMeador/hermes/cmd/airhorn"
	"github.com/DylanMeador/hermes/discord"
	"github.com/DylanMeador/hermes/emojis"
	"github.com/DylanMeador/hermes/gifs"
	"github.com/DylanMeador/hermes/shaco"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"strings"
)

var usageTemplate = `Usage:{{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
{{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
{{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}`


func Cmd(s *discordgo.Session, m *discordgo.MessageCreate) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hermes",
		Short: shaco.Quotes[rand.Intn(len(shaco.Quotes))],
	}

	cmd.SetOut(responseWriter{s, m})
	cmd.SetHelpCommand(&cobra.Command{Use: "nope", Hidden: true})
	cmd.PersistentFlags().Bool("help", false, "none")
	cmd.PersistentFlags().MarkHidden("help")
	cmd.SetUsageTemplate(usageTemplate)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	cmd.AddCommand(airhorn.Cmd())

	return cmd
}

func Execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Println(m.Author.String() + ": " + m.Content)

	cmd := Cmd(s, m)

	args := strings.Split(m.Content, " ")[1:]
	cmd.SetArgs(args)

	if err := cmd.ParseFlags(args); err != nil {
		addReactions(s, m.ChannelID, m.Message.ID, emojis.POOP, emojis.FLAG)
	} else if err := cobra.OnlyValidArgs(cmd, args); err != nil {
		addReactions(s, m.ChannelID, m.Message.ID, emojis.POOP, emojis.C, emojis.O, emojis.M, emojis.M, emojis.A, emojis.N, emojis.D)
	} else {
		if err := cmd.ExecuteContext(discord.GenerateDiscordContext(s, m)); err != nil {
			bug(err, s, m.ChannelID)
		}
	}
}

func addReactions(s *discordgo.Session, channelID string, messageID string, emojiIDs... string) {
	for _, emojiID := range emojiIDs {
		err := s.MessageReactionAdd(channelID, messageID, emojiID)
		if err != nil {
			bug(err, s, channelID)
			break
		}
	}
}

func bug(err error, s *discordgo.Session, channelID string) {
	log.Println(err)
	_, err = s.ChannelMessageSend(channelID, gifs.BUG)
	if err != nil {
		log.Println(err)
	}
}

type responseWriter struct {
	s *discordgo.Session
	m *discordgo.MessageCreate
}

func (rw responseWriter) Write(p []byte) (int, error) {
	message := string(p)
	message = strings.TrimSpace(message)

	if len(message) == 0 {
		return len(p), nil
	}

	_, err := rw.s.ChannelMessageSend(rw.m.ChannelID, message)
	return len(p), err
}