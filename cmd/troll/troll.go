package troll

import (
	"fmt"
	"github.com/DylanMeador/hermes/discord"
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

var soundCache map[string][][]byte
var mux sync.Mutex

type args struct {
	channelName string
}

func Cmd() *cobra.Command {
	a := &args{}

	cmd := &cobra.Command{
		Use:    "troll",
		Short:  "Shaco sounds...mostly",
		RunE:   a.run,
	}

	cmd.PersistentFlags().StringVarP(&a.channelName, "channel", "c", "", "the voice channel to play the troll in")

	return cmd
}

func (a *args) run(cmd *cobra.Command, args []string) error {
	s, m := discord.GetSessionAndMessageFromContext(cmd.Context())

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
	if len(a.channelName) > 0 {
		for _, c := range g.Channels {
			if c.Type == discordgo.ChannelTypeGuildVoice && strings.EqualFold(c.Name, a.channelName) {
				channelID = c.ID
				break
			}
		}
		if len(channelID) == 0 {
			cmd.PrintErrln("Channel " + a.channelName + " does not exist.")
			return nil
		}
	} else {
		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				channelID = vs.ChannelID
			}
		}
		if len(channelID) == 0 {
			cmd.PrintErrln("You are not in a voice channel.")
			return nil
		}
	}

	randomShacoSound := sounds.ALL_SHACO[rand.Intn(len(sounds.ALL_SHACO))]

	sound, err := loadSound(randomShacoSound)
	if err != nil {
		log.Println("Error loading sound: ", err)
	}

	fmt.Println("playing: " + randomShacoSound)

	err = playSound(s, g.ID, channelID, sound)
	if err != nil {
		return err
	}

	if randomShacoSound == sounds.SHACO_JOKE {
		//for key, vc := range s.VoiceConnections {
		//	fmt.Println(key)
		//	fmt.Println(vc.UserID)
		//	fmt.Println(m.Author.ID)
		//}
		//s.GuildMemberMove()
		return s.GuildMemberMove(m.GuildID, m.Author.ID, "")
	}

	return nil
}

// loadSound attempts to load an encoded sound file from disk.
func loadSound(path string) ([][]byte, error) {
	if sound, ok := soundCache[path]; ok {
		return sound, nil
	}

	mux.Lock()
	defer mux.Unlock()

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
