package cmd

import (
	"github.com/DylanMeador/hermes/cmd/annoy"
	"github.com/DylanMeador/hermes/cmd/blame"
	"github.com/DylanMeador/hermes/cmd/bugtest"
	"github.com/DylanMeador/hermes/cmd/echo"
	"github.com/DylanMeador/hermes/cmd/hype"
	"github.com/DylanMeador/hermes/cmd/invite"
	"github.com/DylanMeador/hermes/cmd/trick"
	unmute "github.com/DylanMeador/hermes/cmd/unumute"
	"github.com/DylanMeador/hermes/pkg/discord"
	"github.com/DylanMeador/hermes/pkg/emojis"
	"github.com/DylanMeador/hermes/pkg/errors"
	"github.com/DylanMeador/hermes/pkg/gifs"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
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
	cmd.DisableSuggestions = true

	cmd.AddCommand(annoy.Cmd())
	cmd.AddCommand(blame.Cmd())
	cmd.AddCommand(bugtest.Cmd())
	cmd.AddCommand(echo.Cmd())
	cmd.AddCommand(hype.Cmd())
	cmd.AddCommand(invite.Cmd())
	cmd.AddCommand(trick.Cmd())
	cmd.AddCommand(unmute.Cmd())

	return cmd
}

func Execute(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Println(m.Author.String() + ": " + m.Content)

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		bug(err, s, m)
	}

	if channel.Type == discordgo.ChannelTypeDM {
		log.Println("dm received")
		s.ChannelMessageSend(m.ChannelID, gifs.ALL_DM[rand.Intn(len(gifs.ALL_DM))])
		return
	}

	commandText := strings.TrimPrefix(m.Content, "||")
	commandText = strings.TrimSuffix(m.Content, "||")

	args := strings.Split(commandText, " ")[1:]
	cmd := Cmd(responseWriter{s, m})
	cmd.SetArgs(args)

	targetCommand, _, commandErr := cmd.Find(args)
	flagErr := targetCommand.Flags().Parse(args)

	isHiddenCommand := m.Content != commandText

	if isHiddenCommand {
		err := s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		if err != nil {
			bug(err, s, m)
		}

		// spoiler free mode!
		cmd.Hidden = true
		cmd.SetUsageTemplate(gifs.SPOILER)
	}

	var errEmojis []string
	if flagErr != nil && flagErr != pflag.ErrHelp {
		errEmojis = append(errEmojis, emojis.FLAG)
	}
	if commandErr != nil {
		errEmojis = append(errEmojis, emojis.C, emojis.O, emojis.M, emojis.M2, emojis.A, emojis.N, emojis.D)
	}

	if len(errEmojis) > 0 {
		if !isHiddenCommand {
			addReactions(s, m, emojis.POOP)
			addReactions(s, m, errEmojis...)
		}
	} else {
		if err := cmd.ExecuteContext(discord.GenerateDiscordContext(s, m, isHiddenCommand)); err != nil {
			if err == errors.CommandArgumentErr {
				if !isHiddenCommand {
					addReactions(s, m, emojis.POOP)
					addReactions(s, m, emojis.A, emojis.R, emojis.G)
				}
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
	_, err = s.ChannelMessageSend(m.ChannelID, gifs.JOKER_BRAVO+"\n"+gifs.BUG)
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
