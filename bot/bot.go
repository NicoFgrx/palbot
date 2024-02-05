package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/NicoFgrx/palbot/api"
	"github.com/bwmarrin/discordgo"
)

const prefix string = "!palbot"

// Store Bot API Tokens:
var (
	BotToken  string
	ChannelID string
	Client    api.Client
	Lxc       string
	Node      string
)

func Run() {
	// Create new Discord Session
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal(err)
	}

	// Listen of createMessgae events
	discord.AddHandler(messageCreate)
	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}

	// start goroutine to check network traffic under 1k/s bandwith
	go checkTraffic(1000000, discord)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Prevent loop with bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// prefix check
	args := strings.Split(m.Content, " ")
	if args[0] != prefix {
		return
	}

	// Custom command
	if args[1] == "start" {
		// Get current status before start
		status, err := Client.Status(Node, Lxc)
		if err != nil {
			s.ChannelMessageSend(ChannelID, "Error while getting status of server")
			log.Fatal("Error while getting status of server")

		}

		if status.Data.Status == "running" {
			s.ChannelMessageSend(ChannelID, "Server is already running.. skipped")
			return
		}

		// Start the server
		err = Client.Start(Node, Lxc)
		if err != nil {
			s.ChannelMessageSend(ChannelID, "Error while starting server")
			log.Fatal("Error while starting server")

		}

		s.ChannelMessageSend(ChannelID, "Starting server... wait few minutes")
		return

	} else if args[1] == "shutdown" {
		// Get current status before stop
		status, err := Client.Status(Node, Lxc)
		if err != nil {
			s.ChannelMessageSend(ChannelID, "Error while getting status of server")
			log.Fatal("Error while getting status of server")

		}

		if status.Data.Status == "stopped" {
			s.ChannelMessageSend(ChannelID, "Server is already stopped.. skipped")
			return
		}

		// Shutdown the server
		err = Client.Shutdown(Node, Lxc)
		if err != nil {
			s.ChannelMessageSend(ChannelID, "Error while shutdown the server")
			log.Fatal("Error while shutdown server")

		}

		s.ChannelMessageSend(ChannelID, "Shutdown server... GO TO SLEEP")
		return

	} else if args[1] == "status" {
		// Get status
		status, err := Client.Status(Node, Lxc)
		if err != nil {
			s.ChannelMessageSend(ChannelID, "Error while getting status of server")
			log.Fatal("Error while getting status of server")

		}

		if status.Data.Status == "running" {
			//
			advanced_status := formatFieldsStatus(status)
			s.ChannelMessageSendEmbed(ChannelID, &discordgo.MessageEmbed{
				Title:       ":man_mage: palbot :man_mage:",
				Fields:      advanced_status,
				Description: fmt.Sprintf("Server is currently %s", status.Data.Status),
			})
			return
		}
		s.ChannelMessageSend(ChannelID, fmt.Sprintf("Server is currently %s", status.Data.Status))
		return
	} else {
		s.ChannelMessageSend(ChannelID, "Current available commands are : start, shutdown and status")
		return
	}

}

func formatFieldsStatus(data api.LXCStatusResponse) []*discordgo.MessageEmbedField {
	var result []*discordgo.MessageEmbedField

	resources := discordgo.MessageEmbedField{
		Name:   "Resources",
		Inline: true,
	}
	used := discordgo.MessageEmbedField{
		Name:   "Used (%)",
		Inline: true,
	}

	// vCPU
	resources.Value += "vCPU\n"
	used.Value += fmt.Sprintf("%d\n", int(data.Data.CPU*100))

	// RAM
	resources.Value += "RAM\n"
	ram_used := data.Data.Mem
	ram_max := data.Data.Maxmem
	ram_average := (ram_used * 100) / ram_max
	used.Value += fmt.Sprintf("%d\n", ram_average)

	result = append(result, &resources, &used)

	return result
}

// Check by period the limit netout given in parameters
func checkTraffic(limit int, discord *discordgo.Session) {

	// overide rate limit at 1kbis/s, testing purpose
	// limit = 1000000

	// Init ticket
	duration := 10 * time.Second // FIXME Second => Minute
	ticker := time.NewTicker(duration)
	// fmt.Print("DEBUG : in go routine\n")

	// check each 10 min
	for range ticker.C {

		// Get current netinput
		status, err := Client.Status(Node, Lxc)
		if err != nil {
			log.Fatal("Error while checking traffic input")
		}
		// fmt.Printf("DEBUG : in ticker loop, status=%s netout=%d\n", status.Data.Status, status.Data.Netout)

		// Check if  server is running and current netin is under limit
		if status.Data.Status == "running" && status.Data.Netout <= limit {
			// send message
			discord.ChannelMessageSend(ChannelID, fmt.Sprintf("Low network output detected, can we shut the server ?"))
		}

	}

}
