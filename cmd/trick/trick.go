package trick

import (
	"fmt"
	"github.com/DylanMeador/hermes/discord"
	"github.com/DylanMeador/hermes/errors"
	"github.com/DylanMeador/hermes/sounds"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/spf13/cobra"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
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

var soundCache map[string][][]byte
var mux sync.Mutex

type args struct {
	channelName    string
	forceDisappear bool
	forceJoke      bool
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
	cmd.PersistentFlags().BoolVarP(&a.forceJoke, "joke", "j", false, "force the joke voice to be played and user muted")
	cmd.PersistentFlags().MarkHidden("disappear")
	cmd.PersistentFlags().MarkHidden("joke")

	return cmd
}

func (a *args) run(command *cobra.Command, args []string) error {
	s, m := discord.GetSessionAndMessageFromContext(command.Context())

	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return err
	}
	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return err
	}

	var channelID string
	userIsInChannel := false
	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
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
			command.PrintErrln("Channel " + a.channelName + " does not exist.")
			return errors.CommandArgumentErr
		}
	}

	if channelID != "" {
		return a.playRandomShacoSound(s, m, channelID, userIsInChannel)
	}

	_, err = s.ChannelMessageSend(m.ChannelID, quotes[rand.Intn(len(quotes))])
	return err
}

func (a *args) playRandomShacoSound(s *discordgo.Session, m *discordgo.MessageCreate, channelID string, userInChannel bool) error {
	mux.Lock()
	defer mux.Unlock()

	randomShacoSound := sounds.ALL_SHACO[rand.Intn(len(sounds.ALL_SHACO))]

	if a.forceDisappear {
		randomShacoSound = sounds.SHACO_JOKE
	}
	if a.forceJoke {
		randomShacoSound = sounds.SHACO_ATTACK3
	}

	sound, err := loadSound(randomShacoSound)
	if err != nil {
		log.Println("Error loading sound: ", err)
		return err
	}

	fmt.Println("playing: " + randomShacoSound)

	err = playSound(s, m.GuildID, channelID, sound)
	if err != nil {
		return err
	}

	if !userInChannel {
		return nil
	}

	if randomShacoSound == sounds.SHACO_JOKE {
		// For my next trick, I'll make you disappear!
		data := struct {
			ChannelID *string `json:"channel_id"`
		}{nil}

		guildMember := discordgo.EndpointGuildMember(m.GuildID, m.Author.ID)

		_, err = s.RequestWithBucketID("PATCH", guildMember, data, discordgo.EndpointGuildMember(m.GuildID, ""))
		if err != nil {
			return err
		}
	}else if randomShacoSound == sounds.SHACO_ATTACK3 {
		// The joke's on you!
		data := struct {
			Mute   bool `json:"mute"`
		}{ true }

		guildMember := discordgo.EndpointGuildMember(m.GuildID, m.Author.ID)
		_, err = s.RequestWithBucketID("PATCH", guildMember, data, discordgo.EndpointGuildMember(m.GuildID, ""))
		if err != nil {
			return err
		}

		sound, err = loadSound(sounds.SHACO_LAUGH2)
		if err != nil {
			log.Println("Error loading sound: ", err)
			return err
		}

		time.Sleep(time.Second * 5)

		err = playSound(s, m.GuildID, channelID, sound)
		if err != nil {
			return err
		}

		time.Sleep(time.Second * 15)

		data.Mute = false
		_, err = s.RequestWithBucketID("PATCH", guildMember, data, discordgo.EndpointGuildMember(m.GuildID, ""))
		if err != nil {
			return err
		}

		time.Sleep(time.Second * 2)

		sound, err = loadSound(sounds.SHACO_LAUGH3)
		if err != nil {
			log.Println("Error loading sound: ", err)
			return err
		}
		err = playSound(s, m.GuildID, channelID, sound)
		if err != nil {
			return err
		}
	}

	return nil
}

// loadSound attempts to load an encoded sound file from disk.
func loadSound(path string) ([][]byte, error) {
	if sound, ok := soundCache[path]; ok {
		return sound, nil
	}

	file, err := os.Open(path)
	if err != nil {
		log.Println("Error opening dca file :", err)
		return nil, err
	}

	var buffer [][]byte

	decoder := dca.NewDecoder(file)

	for {
		frame, err := decoder.OpusFrame()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		buffer = append(buffer, frame)
	}

	return buffer, nil
}

// playSound plays the current soundBuffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string, sound [][]byte) (err error) {
	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the soundBuffer data.
	for _, buff := range sound {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
