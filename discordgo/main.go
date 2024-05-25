package main

import (
	"github.com/bwmarrin/discordgo"

	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	token, exists := os.LookupEnv("TOKEN")
	if !exists {
		log.Fatal("No token provided. Exiting")
	}
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error while initializing Discord session:", err)
	}
	err = session.Open()
	if err != nil {
		log.Fatal("Error opening connection,", err)
	}

	fmt.Println("Joining whatsapp")
	_, err = session.ChannelVoiceJoin("468801766889357313", "778600399745712149", true, true)
	if err != nil {
		fmt.Println(err)
	}

	session.AddHandler(messageCreate)

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	
	defer session.Close()
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	session.ChannelMessageSend(message.ChannelID, "You're stupid")
}
