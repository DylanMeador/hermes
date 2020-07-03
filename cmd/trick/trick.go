package trick

import (
	"github.com/DylanMeador/hermes/pkg/discord"
	"github.com/DylanMeador/hermes/pkg/gifs"
	"github.com/DylanMeador/hermes/pkg/sounds"
	"github.com/spf13/cobra"
	"math/rand"
	"time"
)

type args struct {
	forceDisappear bool
	forceJoke      bool
	forceMagic     bool
}

func Cmd() *cobra.Command {
	a := &args{}

	cmd := &cobra.Command{
		Use:   "trick",
		Short: "why so serious",
		Long:  "Crafted long ago as a plaything for a lonely prince, the enchanted marionette Shaco now delights in murder and mayhem.",
		RunE:  a.run,
	}

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

	voiceChannelID, err := hc.GetCommandUserVoiceChannelID()
	if err != nil {
		return err
	}

	if voiceChannelID != "" {
		err = a.performRandomTrick(hc, voiceChannelID)
	} else {
		_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, gifs.JOKER_SAD)
		_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, "I can only play tricks on people in voice channels.")
	}

	return err
}

func (a *args) performRandomTrick(hc *discord.HermesCommand, channelID string) error {
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

				time.Sleep(time.Second * 2)
			} else if calls == 2 {
				time.Sleep(time.Second * 4)
			} else {
				return discord.Unmute(hc, hc.Message.Author.ID)
			}

			return nil
		}
	}

	return discord.PlaySounds(hc, channelID, postSoundCb, soundsToPlay...)
}
