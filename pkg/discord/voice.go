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
	"time"
)

var soundCache = make(map[sounds.Sound][][]byte)
var mux sync.Mutex

func PlaySoundBytes(hc *HermesCommand, channelID string, soundBytes [][]byte) error {
	mux.Lock()
	defer mux.Unlock()

	guildID := hc.Message.GuildID
	vc, err := hc.Session.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		return err
	}
	defer vc.Disconnect()

	return playSoundInChannel(vc, soundBytes)
}

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
		soundBytes, err := loadSound(sound)
		if err != nil {
			log.Println("Error loading sound: ", err)
			return err
		}
		log.Println("playing: " + sound)
		err = playSoundInChannel(vc, soundBytes)
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

func RecordSounds(hc *HermesCommand, channelID string, duration time.Duration) ([][]byte, error) {
	mux.Lock()
	defer mux.Unlock()

	guildID := hc.Message.GuildID
	vc, err := hc.Session.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		return nil, err
	}
	defer vc.Disconnect()

	soundBytes := recordChannelSounds(vc, duration)
	return soundBytes, nil
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
				sb.WriteString(" If so, try the following command when you are in a voice channel:\n")
				sb.WriteString("> hermes unmute -u " + hc.Message.Author.Mention())

				_, err = hc.Session.ChannelMessageSend(hc.Message.ChannelID, sb.String())
				return err
			}
			return nil
		}
	}

	return err
}

func playSoundInChannel(vc *discordgo.VoiceConnection, soundBytes [][]byte) error {
	err := vc.Speaking(true)
	if err != nil {
		return err
	}

	for _, buff := range soundBytes {
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

	soundCache[sound] = buffer

	return soundCache[sound], nil
}

func recordChannelSounds(vc *discordgo.VoiceConnection, duration time.Duration) [][]byte {
	var buffer [][]byte
	ticker := time.NewTicker(duration)

	for {
		select {
		case packet := <-vc.OpusRecv:
			buffer = append(buffer, packet.Opus)
		case <-ticker.C:
			return buffer
		}
	}
}