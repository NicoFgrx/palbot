package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	api "github.com/NicoFgrx/palbot/api"
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

	client := Config()
	node := "pve"
	lxc := "4130"

	// Get initial status
	fmt.Printf("[+] Get Status of %s LXC container\n", lxc)
	status, err := client.Status(node, lxc)
	if err != nil {
		log.Fatalf("[-] Error while getting status%s : %s", lxc, err)
	}
	fmt.Printf("[+] %s is currently %s\n", lxc, status.Data.Status)

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Shutdown LXC
	fmt.Printf("[+] Shutdown %s LXC container\n", lxc)
	err = client.Shutdown(node, lxc)
	if err != nil {
		log.Fatalf("[-] Error while shutdown %s : %s", lxc, err)
	}

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Get new status
	fmt.Printf("[+] Get Status of %s LXC container\n", lxc)
	status, err = client.Status(node, lxc)
	if err != nil {
		log.Fatalf("[-] Error while getting status%s : %s", lxc, err)
	}
	fmt.Printf("[+] %s is currently %s\n", lxc, status.Data.Status)

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Start LXC
	// Shutdown LXC
	fmt.Printf("[+] Start %s LXC container\n", lxc)
	err = client.Start(node, lxc)
	if err != nil {
		log.Fatalf("[-] Error while shutdown %s : %s", lxc, err)
	}

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Get new status
	fmt.Printf("[+] Get Status of %s LXC container\n", lxc)
	status, err = client.Status(node, lxc)
	if err != nil {
		log.Fatalf("[-] Error while getting status%s : %s", lxc, err)
	}
	fmt.Printf("[+] %s is currently %s\n", lxc, status.Data.Status)

}
