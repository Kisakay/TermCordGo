package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	token       string
	channelID   string
	discord     *discordgo.Session
	stopChannel chan bool
	initialized bool
)

func main() {
	fmt.Println("\x1b[31;1mTermCord for CLI in Go by Kisakay\x1b[0m")
	fmt.Println("\x1b[31;1m---------------------------------\x1b[0m")

	fmt.Print("Enter your Discord bot token: ")
	fmt.Scanln(&token)

	fmt.Print("Enter the channel ID: ")
	fmt.Scanln(&channelID)

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
		return
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening Discord session: ", err)
		return
	}

	fmt.Println("Bot is now running. Type '/stop' to exit.")
	fmt.Println("\x1b[31;1mawaiting for new message...\x1b[0m")

	stopChannel = make(chan bool)
	<-stopChannel
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID != channelID {
		return
	}

	if !initialized {
		initialized = true
		return
	}

	fmt.Printf("[+] %s: %s\n", m.Author.Username, m.Content)

	if strings.HasPrefix(m.Content, "/stop") {
		stopChannel <- true
		return
	}

	fmt.Print("You: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()

	_, err := s.ChannelMessageSend(m.ChannelID, text)
	fmt.Println("\x1b[31;1mawaiting for new message...\x1b[0m")

	if err != nil {
		fmt.Println("Error sending message: ", err)
	}
}
