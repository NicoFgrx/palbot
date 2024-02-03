package main

import (
	"fmt"
	"log"
	"os"

	"runtime"

	"github.com/NicoFgrx/palbot-monitoring/api"
	"github.com/NicoFgrx/palbot-monitoring/bot"
	"github.com/joho/godotenv"
)

func Config() api.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[-] Error loading .env file")
	}
	url := os.Getenv("PROXMOX_VE_ENDPOINT")            // must be https://x.x.x.x:8006/api2/json
	id := os.Getenv("PROXMOX_VE_API_TOKEN_ID")         // must be user@pve!id
	secret := os.Getenv("PROXMOX_VE_API_TOKEN_SECRET") // must be xxx-xxx-xx

	fmt.Println("[+] Creating client")

	client := api.NewClient(url, id, secret)

	return *client
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("[ERROR] Error loading .env file")
	}

	// Setup Proxmox api client
	client := Config()

	node := os.Getenv("PROXMOX_NODE")
	lxc := os.Getenv("PROXMOX_LXC_VMID")

	// Setup Bot
	botToken, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		log.Fatal("[ERROR] Must set Discord token as env variable: DISCORD_TOKEN")
	}

	channelID, ok := os.LookupEnv("DISCORD_CHANNEL")
	if !ok {
		log.Fatal("[ERROR] Must set Discord channel ID as env variable: DISCORD_CHANNEL")
	}
	fmt.Printf("[INFO] Currently running on %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// Save API keys & start bot
	bot.BotToken = botToken
	bot.ChannelID = channelID
	bot.Client = client
	bot.Node = node
	bot.Lxc = lxc
	bot.Run()

}
