package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)


func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	token, exists := os.LookupEnv("TOKEN")
	if !exists {
		log.Fatal("No token provided. Exiting")
	}

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error while initializing Discord session:", err)
	}
	sess.Open()

	fmt.Println("Joining whatsapp")
	_, err = sess.ChannelVoiceJoin("468801766889357313", "778600399745712149", true, true)
	if err != nil {
		fmt.Println(err)
	}
	
	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	defer sess.Close()
}
