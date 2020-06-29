package cmd

import (
	"github.com/DylanMeador/hermes/cmd/invite"
	"github.com/DylanMeador/hermes/cmd/trick"
	"github.com/DylanMeador/hermes/discord"
	"github.com/DylanMeador/hermes/emojis"
	"github.com/DylanMeador/hermes/errors"
	"github.com/DylanMeador/hermes/gifs"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"log"
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

func Cmd(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hermes",
		Short: "A conductor of souls into the afterlife.",
	}

	cmd.SetOut(out)
	cmd.SetHelpCommand(&cobra.Command{Use: "nope", Hidden: true})
	cmd.InitDefaultHelpFlag()
	cmd.SetUsageTemplate(usageTemplate)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	cmd.AddCommand(invite.Cmd())
	cmd.AddCommand(trick.Cmd())

	return cmd
}

func Execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Println(m.Author.String() + ": " + m.Content)

	args := strings.Split(m.Content, " ")[1:]
	cmd := Cmd(responseWriter{s, m})
	cmd.SetArgs(args)

	targetCommand, _, commandErr := cmd.Find(args)
	flagErr := targetCommand.Flags().Parse(args)

	var errEmojis []string
	if flagErr != nil && flagErr != pflag.ErrHelp {
		errEmojis = append(errEmojis, emojis.FLAG)
	}
	if commandErr != nil {
		errEmojis = append(errEmojis, emojis.C, emojis.O, emojis.M, emojis.M2, emojis.A, emojis.N, emojis.D)
	}

	if len(errEmojis) > 0 {
		addReactions(s, m, emojis.POOP)
		addReactions(s, m, errEmojis...)
	} else {
		if err := cmd.ExecuteContext(discord.GenerateDiscordContext(s, m)); err != nil {
			if err == errors.CommandArgumentErr {
				addReactions(s, m, emojis.POOP)
				addReactions(s, m, emojis.A, emojis.R, emojis.G)
			} else {
				bug(err, s, m)
			}
		}
	}
}

func addReactions(s *discordgo.Session, m *discordgo.MessageCreate, emojiIDs ...string) {
	for _, emojiID := range emojiIDs {
		err := s.MessageReactionAdd(m.ChannelID, m.Message.ID, emojiID)
		if err != nil {
			bug(err, s, m)
			break
		}
	}
}

func bug(err error, s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Println(err)
	_, err = s.ChannelMessageSend(m.ChannelID, gifs.BUG)
	if err != nil {
		log.Println(err)
	}

	const hacker = "#hacker"
	if !strings.Contains(m.Author.Username, hacker) {
		err = s.GuildMemberNickname(m.GuildID, m.Author.ID, m.Author.Username+hacker)
		if err != nil {
			log.Println(err)
		}
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
