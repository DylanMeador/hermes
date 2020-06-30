package trick

import (
	"github.com/DylanMeador/hermes/pkg/discord"
	"github.com/DylanMeador/hermes/pkg/errors"
	"github.com/DylanMeador/hermes/pkg/gifs"
	"github.com/DylanMeador/hermes/pkg/sounds"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"strings"
	"time"
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

type args struct {
	channelName    string
	forceDisappear bool
	forceJoke      bool
	forceMagic     bool
}

func Cmd() *cobra.Command {
	a := &args{}

	cmd := &cobra.Command{
		Use:   "trick",
		Short: "the demon jester, Shaco",
		Long:  "Crafted long ago as a plaything for a lonely prince, the enchanted marionette Shaco now delights in murder and mayhem.",
		RunE:  a.run,
	}

	cmd.PersistentFlags().StringVarP(&a.channelName, "channel", "c", "", "the voice channel to play tricks in")
	cmd.PersistentFlags().BoolVarP(&a.forceDisappear, "disappear", "d", false, "force the disappear voice to be played and user removed from voice channel")
	cmd.PersistentFlags().BoolVarP(&a.forceJoke, "joke", "j", false, "force the joke voice to be played and user muted then unmuted")
	cmd.PersistentFlags().BoolVarP(&a.forceMagic, "magic", "m", false, "force the magic voice to be played and magic gif posted")
	cmd.PersistentFlags().MarkHidden("disappear")
	cmd.PersistentFlags().MarkHidden("joke")
	cmd.PersistentFlags().MarkHidden("magic")

	return cmd
}

func (a *args) run(command *cobra.Command, args []string) error {
	hc := discord.GetHermesCommandFromContext(command.Context())

	g, err := hc.Session.State.Guild(hc.Message.GuildID)
	if err != nil {
		return err
	}

	var channelID string
	userIsInChannel := false
	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == hc.Message.Author.ID {
			channelID = vs.ChannelID
			userIsInChannel = true
			break
		}
	}
	if len(a.channelName) > 0 {
		for _, c := range g.Channels {
			if c.Type == discordgo.ChannelTypeGuildVoice && strings.EqualFold(c.Name, a.channelName) {
				channelID = c.ID
				break
			}
		}
		if len(channelID) == 0 {
			log.Println("Channel " + a.channelName + " does not exist.")
			return errors.CommandArgumentErr
		}
	}

	if channelID != "" {
		err = a.playRandomShacoSound(hc, channelID, userIsInChannel)
	} else {
		_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, quotes[rand.Intn(len(quotes))])
	}

	return err
}

func (a *args) playRandomShacoSound(hc *discord.HermesCommand, channelID string, userInChannel bool) error {
	allTrickSounds := []sounds.Sound{sounds.SHACO_MAGIC_TRICK, sounds.SHACO_MAKE_YOU_DISAPPEAR, sounds.SHACO_JOKES_ON_YOU}
	trickSound := allTrickSounds[rand.Intn(len(allTrickSounds))]

	if a.forceDisappear {
		trickSound = sounds.SHACO_MAKE_YOU_DISAPPEAR
	}
	if a.forceJoke {
		trickSound = sounds.SHACO_JOKES_ON_YOU
	}
	if a.forceMagic {
		trickSound = sounds.SHACO_MAGIC_TRICK
	}

	soundsToPlay := []sounds.Sound{trickSound}
	var postSoundCb func() error

	if userInChannel {
		// 	How about a magic trick?
		if trickSound == sounds.SHACO_MAGIC_TRICK {
			postSoundCb = func() error {
				_, err := hc.Session.ChannelMessageSend(hc.Message.ChannelID, gifs.ALL_MAGIC[rand.Intn(len(gifs.ALL_MAGIC))])
				return err
			}
		}

		// For my next trick, I'll make you disappear!
		if trickSound == sounds.SHACO_MAKE_YOU_DISAPPEAR {
			postSoundCb = func() error {
				return discord.RemoveFromChannel(hc, hc.Message.Author.ID)
			}
		}

		// The joke's on you!
		if trickSound == sounds.SHACO_JOKES_ON_YOU {
			soundsToPlay = append(soundsToPlay, sounds.SHACO_LAUGH2)
			soundsToPlay = append(soundsToPlay, sounds.SHACO_LAUGH3)

			calls := 0
			postSoundCb = func() error {
				calls = calls + 1
				if calls == 1 {
					err := discord.Mute(hc, hc.Message.Author.ID)
					if err != nil {
						return err
					}

					time.Sleep(time.Second * 5)
				} else if calls == 2 {
					time.Sleep(time.Second * 15)
				} else {
					return discord.Unmute(hc, hc.Message.Author.ID)
				}

				return nil
			}
		}
	}

	return discord.PlaySounds(hc, channelID, postSoundCb, soundsToPlay...)
}
