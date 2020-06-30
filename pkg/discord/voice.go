package discord

import (
	"github.com/DylanMeador/hermes/pkg/sounds"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

var soundCache map[sounds.Sound][][]byte
var mux sync.Mutex

func PlaySound(hc *HermesCommand, channelID string, sound sounds.Sound) error {
	return PlaySounds(hc, channelID, nil, sound)
}

func PlaySounds(hc *HermesCommand, channelID string, postSoundCb func() error, sounds ...sounds.Sound) error {
	mux.Lock()
	defer mux.Unlock()

	guildID := hc.Message.GuildID
	vc, err := hc.Session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}
	defer vc.Disconnect()

	for _, sound := range sounds {
		err = playSoundInChannel(vc, sound)
		if err != nil {
			return err
		}
		if postSoundCb != nil {
			err = postSoundCb()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func RemoveFromChannel(hc *HermesCommand, userID string) error {
	data := struct {
		ChannelID *string `json:"channel_id"`
	}{nil}

	guildID := hc.Message.GuildID
	guildMember := discordgo.EndpointGuildMember(guildID, userID)

	_, err := hc.Session.RequestWithBucketID("PATCH", guildMember, data, discordgo.EndpointGuildMember(guildID, ""))
	if rErr, ok := err.(*discordgo.RESTError); ok {
		// ignore error if user not in channel
		if rErr.Message != nil && rErr.Message.Code == 40032 {
			_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, "You escaped my tricks this time.")
			return err
		}
	}
	return err
}

func Mute(hc *HermesCommand, userID string) error {
	return mute(hc, userID, true)
}

func Unmute(hc *HermesCommand, userID string) error {
	return mute(hc, userID, false)
}

func mute(hc *HermesCommand, userID string, mute bool) error {
	data := struct {
		Mute bool `json:"mute"`
	}{mute}

	guildID := hc.Message.GuildID
	guildMember := discordgo.EndpointGuildMember(guildID, userID)
	_, err := hc.Session.RequestWithBucketID("PATCH", guildMember, data, discordgo.EndpointGuildMember(guildID, ""))

	if rErr, ok := err.(*discordgo.RESTError); ok {
		// if we failed to unmute, help the poor user out :-)
		if rErr.Message != nil && rErr.Message.Code == 40032 {
			if mute == false {
				var sb strings.Builder
				sb.WriteString("It seems you have avoided my tricks, but that may have caused you some pain.")
				sb.WriteString(" Maybe this will help.\n")
				sb.WriteString("> hermes unmute")

				_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, sb.String())
				return err
			}
			return nil
		}
	}

	return err
}

func playSoundInChannel(vc *discordgo.VoiceConnection, sound sounds.Sound) error {
	buffer, err := loadSound(sound)
	if err != nil {
		log.Println("Error loading sound: ", err)
		return err
	}

	err = vc.Speaking(true)
	if err != nil {
		return err
	}

	log.Println("playing: " + sound)

	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	return vc.Speaking(false)
}

func loadSound(sound sounds.Sound) ([][]byte, error) {
	if sound, ok := soundCache[sound]; ok {
		return sound, nil
	}

	file, err := os.Open(string(sound))
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
