package cmd

import (
	"github.com/DylanMeador/hermes/cmd/airhorn"
	"github.com/DylanMeador/hermes/discord"
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

	args := strings.Split(m.Content, " ")
	cmd.SetArgs(args[1:])
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

func Execute(s *discordgo.Session, m *discordgo.MessageCreate) error {
	log.Println(m.Author.ID + ": " + m.Content)

	err := Cmd(s, m).ExecuteContext(discord.GenerateDiscordContext(s, m))
	if err != nil {
		log.Println(err)
		response := "> " + m.Content + "\n"
		response += "The joke's on you!"
		_, err = s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Println(err)
		}
	}
	return err
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