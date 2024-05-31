package main

import (
	"io"
	"time"

	"github.com/bwmarrin/discordgo"

	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	type Payload struct {
		Name string
	}

	payload := &Payload{
		Name: "llama3",
	}

        jsonData, err := json.Marshal(payload)
        if err != nil {
                fmt.Println("Error marshalling JSON:", err)
                return
        }

        url := "http://127.0.0.1:11434/api/pull"

        request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
        if err != nil {
                log.Fatal("Error creating request:", err)
        }
        fmt.Println(request)

        request.Header.Set("Content-Type", "application/json")

        client := &http.Client{}
        response, err := client.Do(request)
        if err != nil {
                log.Fatal("Error sending request:", err)
        }

        defer response.Body.Close()
	// Defer execution until we get 200 OK

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

	if message.ChannelID != "845650803242958858" {
		fmt.Println("Message outside currently watching channels")
		fmt.Println(message.ChannelID)
		return
	}

	// Message represents the structure of the message in the JSON
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	// RequestPayload represents the structure of the JSON payload
	type RequestPayload struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
		Stream   bool      `json:"stream"`
	}

	// ResponseMessage represents the structure of the message in the response JSON
	type ResponseMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	// ResponsePayload represents the structure of the response JSON
	type ResponsePayload struct {
		Model     string          `json:"model"`
		CreatedAt time.Time       `json:"created_at"`
		Message   ResponseMessage `json:"message"`
		Done      bool            `json:"done"`
	}

	fmt.Println(message.Author.Username + ": " + message.Content)

	payload := RequestPayload{
		Model: "llama3",
		Messages: []Message{
			{
				Role:    "user",
				Content: "discord user(" + message.Author.Username + "): " + message.Content,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	url := "http://127.0.0.1:11434/api/chat"

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}
	fmt.Println(request)

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}

	defer response.Body.Close()

	// Create a JSON decoder
	decoder := json.NewDecoder(response.Body)

	// Read and process each JSON object in the response
	for {
		var response ResponsePayload
		if err := decoder.Decode(&response); err == io.EOF {
			// End of the response
			break
		} else if err != nil {
			fmt.Println("Error decoding JSON response:", err)
			return
		}

		// Print the response
		session.ChannelMessageSend(message.ChannelID, response.Message.Content)
	}
}

// TODO: Typing status
// TODO: Short-term memory
// TODO: LLAVA
