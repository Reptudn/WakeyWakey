package commands

import (
	"wakeywakey/database"
	"wakeywakey/utils"

	"github.com/bwmarrin/discordgo"
)

var Wake = discordgo.ApplicationCommand{
	Name:        "wake",
	Description: "Sends a Wake-on-LAN packet to wake up your PC.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "alias",
			Description: "The Alias of the PC to wake up.",
			Required: true,
			Autocomplete: true,
		},
	},
}

func HandleWake(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options
	alias := options[0].StringValue()

	var err error
	macAddress, err := database.GetMacByAlias(i.Member.User.ID, alias)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
				Embeds: []*discordgo.MessageEmbed{
					utils.EmbedError("No PC registered", "No PC registered with alias '" + alias + "': " + err.Error()),
				},
			},
		})
		return
	}

	// First try my implementation
	err = utils.SendWakeOnLANPacket(macAddress)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
				Embeds: []*discordgo.MessageEmbed{
					utils.EmbedError("Failed to send Wake-on-LAN packet", "Could not send Wake-on-LAN packet to '" + alias + "': " + err.Error()),
				},
			},
		})
		return
	}

	// If that fails, try using the wakeonlan command
	err = utils.sendWakeOnLANPacketViaCommand(macAddress)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
				Embeds: []*discordgo.MessageEmbed{
					utils.EmbedError("Failed to send Wake-on-LAN packet via command", "Could not send Wake-on-LAN packet to '" + alias + "' via command: " + err.Error()),
				},
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 1 << 6,
			Embeds: []*discordgo.MessageEmbed{
				utils.EmbedSuccess("Wake-on-LAN packet sent", "Wake-on-LAN packet sent successfully to '" + alias + "'!"),
			},
		},
	})
}

func HandleWakeAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	entries, err := database.GetAllEntriesByUserId(i.Member.User.ID)
	if err != nil {
		return
	}

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, entry := range entries {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  entry.Alias,
			Value: entry.Alias,
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}