package tickets

import (
	"bot/config"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strings"
	"time"
)

func Thisisbetterthanithougtobehonest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Interaction.Type != discordgo.InteractionMessageComponent {
		return
	}

	if i.MessageComponentData().CustomID == "ticket_category" {
		selected := i.MessageComponentData().Values[0]
		categoryName := map[string]string{
			"general":   "General Support",
			"technical": "Technical Support",
			"payment":   "Payment Support",
			"about":     "Bot Support",
		}[selected]

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: fmt.Sprintf("ticket_modal:%s", categoryName),
				Title:    "Open a Support Ticket",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "name",
							Label:       "Your Name",
							Style:       discordgo.TextInputShort,
							Required:    true,
							Placeholder: "Enter your name",
						},
					}},
					discordgo.ActionsRow{Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "question",
							Label:       "What do you need help with?",
							Style:       discordgo.TextInputParagraph,
							Required:    true,
							Placeholder: "Describe your issue.....",
						},
					}},
				},
			},
		})

		go func() {
			_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Components: &[]discordgo.MessageComponent{
					discordgo.ActionsRow{
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
					},
				},
			})
		}()
	}
}

func Nothingtoexplainevensomeonewithoutknowledgewouldcheckit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Interaction.Type != discordgo.InteractionModalSubmit {
		return
	}

	customIDParts := strings.Split(i.ModalSubmitData().CustomID, ":")
	if len(customIDParts) < 2 {
		fmt.Println("Error: No category found in modal CustomID")
		return
	}

	categoryName := customIDParts[1]
	user := i.Member.User
	guildID := i.GuildID
	userID := user.ID
	username := user.Username

	categoryID, err := createcategoryifnotexsistsykilikethatidkwhy(s, guildID, categoryName)
	if err != nil {
		fmt.Println("Error getting/creating category:", err)
		return
	}

	permissionOverwrites := []*discordgo.PermissionOverwrite{
		{ID: guildID, Type: discordgo.PermissionOverwriteTypeRole, Deny: discordgo.PermissionViewChannel},
		{ID: userID, Type: discordgo.PermissionOverwriteTypeMember, Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory},
	}

	if roles, exists := config.CategoryRoleMap[categoryName]; exists {
		for _, roleID := range roles {
			permissionOverwrites = append(permissionOverwrites, &discordgo.PermissionOverwrite{
				ID:    roleID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory,
			})
		}
	}

	channel, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:                 fmt.Sprintf("ticket-%s", username),
		Type:                 discordgo.ChannelTypeGuildText,
		ParentID:             categoryID,
		PermissionOverwrites: permissionOverwrites,
	})

	if err != nil {
		fmt.Println("Error creating ticket channel:", err)
		return
	}

	closeButton := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Close Ticket",
				Style:    discordgo.DangerButton,
				CustomID: "close_ticket",
			},
		},
	}

	_, _ = s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("✅ Your ticket has been created in **%s**.", categoryName),
		Components: []discordgo.MessageComponent{
			closeButton,
		},
	})

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ Your ticket has been created in **%s**: <#%s>", categoryName, channel.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func Handlestheclosebuttonandhandlessomelogicforthetranscriptzkyk(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Interaction.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := i.MessageComponentData().CustomID
	switch customID {
	case "close_ticket":
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "⚠️ Are you sure you want to close this ticket? This action is irreversible.",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Yes, Close",
								Style:    discordgo.DangerButton,
								CustomID: "confirm_close_ticket",
							},
							discordgo.Button{
								Label:    "Cancel",
								Style:    discordgo.SecondaryButton,
								CustomID: "cancel_close_ticket",
							},
						},
					},
				},
			},
		})

	case "confirm_close_ticket":
		channelID := i.ChannelID
		user := i.Member.User
		timestamp := time.Now().Format("2006-01-02")
		filename := fmt.Sprintf("%s-%s.txt", user.Username, timestamp)

		messages, err := fetchMessages(s, channelID)
		if err != nil {
			log.Printf("Error fetching messages: %v", err)
			return
		}

		err = willsavethetranscriptandsavemybrainmaybe(filename, messages)
		if err != nil {
			log.Printf("Error saving transcript: %v", err)
			return
		}

		err = establishedonthesendtranscriptappoarchyk(s, filename, config.TranscriptChannel)
		if err != nil {
			log.Printf("Error sending transcript: %v", err)
			return
		}

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    "✅ This ticket is now closed. Deleting channel...",
				Components: []discordgo.MessageComponent{},
			},
		})

		go func() {
			time.Sleep(3 * time.Second)
			if _, err := s.ChannelDelete(channelID); err != nil {
				log.Printf("Error deleting channel %s: %v", channelID, err)
			}
		}()

	case "cancel_close_ticket":
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Ticket closure canceled.",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Close Ticket",
								Style:    discordgo.DangerButton,
								CustomID: "close_ticket",
							},
						},
					},
				},
			},
		})
	}
}

func fetchMessages(s *discordgo.Session, channelID string) ([]*discordgo.Message, error) {
	var messages []*discordgo.Message
	var lastMessageID string

	for {

		fetchedMessages, err := s.ChannelMessages(channelID, 100, lastMessageID, "", "")
		if err != nil {
			return nil, err
		}

		if len(fetchedMessages) == 0 {
			break
		}

		messages = append(messages, fetchedMessages...)
		lastMessageID = fetchedMessages[len(fetchedMessages)-1].ID
		if len(fetchedMessages) < 100 {
			break
		}
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func willsavethetranscriptandsavemybrainmaybe(filename string, messages []*discordgo.Message) error {
	var messageContents []string

	for _, msg := range messages {
		messageContents = append(messageContents, fmt.Sprintf("%s: %s", msg.Author.Username, msg.Content))
	}

	content := strings.Join(messageContents, "\n")

	return os.WriteFile(filename, []byte(content), 0644)
}

func establishedonthesendtranscriptappoarchyk(s *discordgo.Session, filename string, channelID string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	_, err = s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: "📜 Ticket transcript",
		Files: []*discordgo.File{
			{
				Name:   filename,
				Reader: file,
			},
		},
	})

	return err
}
