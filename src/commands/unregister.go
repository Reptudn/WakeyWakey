package commands

import (
	"fmt"
	"wakeywakey/database"
	"wakeywakey/utils"

	"github.com/bwmarrin/discordgo"
)

var UnregisterWake = discordgo.ApplicationCommand{
	Name:        "unregister",
	Description: "Unregisters a PC by its alias for Wake-on-LAN.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "alias",
			Description: "The alias of the PC to unregister.",
			Required: true,
		},
	},
}

func HandleUnregister(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options
	alias := options[0].StringValue()

	err := database.RemoveWakeEntryByAlias(i.Member.User.ID, alias)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
				Embeds: []*discordgo.MessageEmbed{
					utils.EmbedError("Unregistration Failed", "Failed to unregister PC: "+err.Error()),
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
				utils.EmbedSuccess("Unregistration Successful", "Successfully unregistered PC."),
			},
		},
	})
}

func HandleUnregisterAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userId := i.Member.User.ID

	entries, err := database.GetAllEntriesByUserId(userId)
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

	fmt.Println("Provided autocomplete choices for user " + i.Member.User.Username + " (" + i.Member.User.ID + ")")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}