package airhorn

import (
	"encoding/binary"
	"github.com/DylanMeador/hermes/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var soundBuffer1 = make([][]byte, 0)
var soundBuffer2 = make([][]byte, 0)

var once sync.Once

type args struct {
	channelName string
}

func Cmd() *cobra.Command {
	a := &args{}

	cmd := &cobra.Command{
		Use:    "airhorn",
		Short:  "An airhorn sound will play in your current channel",
		PreRun: a.preRun,
		RunE:   a.run,
	}

	cmd.PersistentFlags().StringVarP(&a.channelName, "channel", "c", "", "the voice channel to play the airhorn in")

	return cmd
}

func (a *args) preRun(cmd *cobra.Command, args []string) {
	// Load the sound file.
	once.Do(func() {
		var err error
		soundBuffer1, err = loadSound("sounds/blame.dca")
		if err != nil {
			log.Println("Error loading sound: ", err)
		}
		soundBuffer2, err = loadSound("sounds/shaco/joke.dca")
		if err != nil {
			log.Println("Error loading sound: ", err)
		}
	})
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

	return playSound(s, g.ID, channelID)
}

// loadSound attempts to load an encoded sound file from disk.
func loadSound(path string) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("Error opening dca file :", err)
		return nil, err
	}

	var opuslen int
	var buffer [][]byte

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return nil, err
			}
			return buffer, nil
		}

		if err != nil {
			log.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			log.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Append encoded pcm data to the soundBuffer.
		buffer = append(buffer, InBuf)
	}
}

// playSound plays the current soundBuffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {
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
	for _, buff := range soundBuffer1 {
		vc.OpusSend <- buff
	}

	for _, buff := range soundBuffer2 {
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
