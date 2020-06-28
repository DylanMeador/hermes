package main

import (
	"fmt"
	"github.com/DylanMeador/hermes/cmd"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	token := "NzI2NjQyNTU0OTUwMzg1NjY1.XvgQtA.Sr3yaJ1No_6bnJWLBMnCXNFum1g"

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(messageUpdate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Here we go!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "hermes") {
		go func() {
			log.Println(m.Author.Username + ": " + m.Content)
			err := cmd.Execute(s, m)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "hermes") {
		s.ChannelMessageSend(m.ChannelID, "The joke's on you!")
	}
}