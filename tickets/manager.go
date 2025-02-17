package tickets

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func createcategoryifnotexsistsykilikethatidkwhy(s *discordgo.Session, guildID, categoryName string) (string, error) {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch guild channels: %v", err)
	}

	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildCategory && ch.Name == categoryName {
			return ch.ID, nil
		}
	}

	category, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name: categoryName,
		Type: discordgo.ChannelTypeGuildCategory,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create category %s: %v", categoryName, err)
	}

	return category.ID, nil
}
