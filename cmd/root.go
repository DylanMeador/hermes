package cmd

import (
	"github.com/DylanMeador/hermes/cmd/airhorn"
	"github.com/DylanMeador/hermes/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"math/rand"
	"strings"
)

var quotes = []string{
	"Look... behind you.",
	"This will be fun!",
	"The joke's on you!",
	"Here we go!",
	"March, march, march, march!",
	"Now you see me, now you don't!",
	"Just a little bit closer!",
	"Why so serious?",
	"For my next trick, I'll make you disappear!",
	"How about a magic trick?",
}

var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}
`

func Cmd(s *discordgo.Session, m *discordgo.MessageCreate) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hermes",
		Short: quotes[rand.Intn(len(quotes))],
	}

	args := strings.Split(m.Content, " ")
	cmd.SetArgs(args[1:])
	cmd.SetOut(responseWriter{s, m})
	cmd.SetErr(errWriter{s, m})
	cmd.SetUsageTemplate(usageTemplate)

	cmd.AddCommand(airhorn.Cmd())

	return cmd
}

func Execute(s *discordgo.Session, m *discordgo.MessageCreate) error {
	return Cmd(s, m).ExecuteContext(discord.GenerateDiscordContext(s, m))
}

type responseWriter struct {
	s *discordgo.Session
	m *discordgo.MessageCreate
}

type errWriter struct {
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

func (rw errWriter) Write(p []byte) (int, error) {
	message := string(p)
	message = strings.TrimSpace(message)

	if len(message) == 0 {
		return len(p), nil
	}

	message = "> " + rw.m.Content + "\n" + message
	_, err := rw.s.ChannelMessageSend(rw.m.ChannelID, message)
	return len(p), err
}
