package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bot/config"
	"bot/tickets"

	"github.com/bwmarrin/discordgo"
)

func main() {
	config.Ladeeinfachdiekonfigjungegerman() //its useless i agree but it looks cool

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	dg.AddHandler(tickets.Thisisbetterthanithougtobehonest)
	dg.AddHandler(tickets.Nothingtoexplainevensomeonewithoutknowledgewouldcheckit)
	dg.AddHandler(tickets.Handlestheclosebuttonandhandlessomelogicforthetranscriptzkyk)

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	fmt.Println("Bot is running")

	tickets.Wegotanembedandiwantyoutosendthatwithadropdown(dg, config.ChannelID)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Offline (ik its not but wont take longer than 1 millisecond).")
	_ = dg.Close()
}
