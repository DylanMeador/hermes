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
	log.Println(m.Author.ID + ": " + m.Content)

	cmd := Cmd(s, m)
	args := strings.Split(m.Content, " ")[1:]

	if err := cmd.ValidateArgs(args); err != nil {
		err = s.MessageReactionAdd(m.ChannelID, m.Message.ID, emojis.POOP)
		if err != nil {
			log.Println(err)
		}
	} else {
		cmd.SetArgs(args)
		if err := cmd.ExecuteContext(discord.GenerateDiscordContext(s, m)); err != nil {
			log.Println(err)
			_, err = s.ChannelMessageSend(m.ChannelID, gifs.BUG)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return
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