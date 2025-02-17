package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	Token             string
	ChannelID         string
	TranscriptChannel string
	CategoryRoleMap   map[string][]string
)

func Ladeeinfachdiekonfigjungegerman() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Token = os.Getenv("DISCORD_TOKEN")
	ChannelID = os.Getenv("CHANNEL_ID")
	TranscriptChannel = os.Getenv("TRANSCRIPTCHANNELID")

	CategoryRoleMap = map[string][]string{
		"General Support":   strings.Split(os.Getenv("GENERAL_SUPPORT_ROLES"), ","),
		"Technical Support": strings.Split(os.Getenv("TECHNICAL_SUPPORT_ROLES"), ","),
		"Payment Support":   strings.Split(os.Getenv("PAYMENT_SUPPORT_ROLES"), ","),
		"Bot Support":       strings.Split(os.Getenv("BOT_SUPPORT_ROLES"), ","),
	}
}
