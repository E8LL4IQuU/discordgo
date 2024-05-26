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

	// session.ChannelMessageSend(message.ChannelID, "You're stupid")
	fmt.Println(message.Content)

	// FIXME: use llama3-text
	payload := RequestPayload{
		Model: "llama3",
		Messages: []Message{
			{
				Role:    "user",
				Content: message.Content,
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
		fmt.Printf("Model: %s\n", response.Model)
		fmt.Printf("CreatedAt: %s\n", response.CreatedAt.Format(time.RFC3339))
		fmt.Printf("Message Role: %s\n", response.Message.Role)
		fmt.Printf("Message Content: %s\n", response.Message.Content)
		fmt.Printf("Done: %v\n", response.Done)
		fmt.Println()
	}
}

// TODO: LLAVA
