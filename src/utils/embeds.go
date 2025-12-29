package utils

import "github.com/bwmarrin/discordgo"

func EmbedError(title string, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       0xFF0000, // Red color for errors
	}
}

func EmbedSuccess(title string, description string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       0x00FF00, // Green color for success
	}
}