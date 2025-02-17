package tickets

import "github.com/bwmarrin/discordgo"

func Wegotanembedandiwantyoutosendthatwithadropdown(s *discordgo.Session, channelID string) {
	embed := &discordgo.MessageEmbed{
		Title:       "Need Support?",
		Description: "Select an option drive on brother (sounds way better in german)",
		Color:       0x00ff00,
	}

	dropdown := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				CustomID:    "ticket_category",
				Placeholder: "Choose a ticket category...",
				Options: []discordgo.SelectMenuOption{
					{Label: "General Questions", Value: "general"},
					{Label: "Technical Support", Value: "technical"},
					{Label: "Payment Support", Value: "payment"},
					{Label: "About the Bot", Value: "about"},
				},
			},
		},
	}

	_, _ = s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: []discordgo.MessageComponent{dropdown},
	})
}
