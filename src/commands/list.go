package commands

import (
	"wakeywakey/database"
	"wakeywakey/utils"

	"github.com/bwmarrin/discordgo"
)

var ListDevices = discordgo.ApplicationCommand{
	Name:        "list",
	Description: "Lists all registered PCs for Wake-on-LAN.",
}

func HandleListDevices(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
	userId := i.Member.User.ID

	deviceList, err := database.GetAllEntriesByUserId(userId)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
				Embeds: []*discordgo.MessageEmbed{
					utils.EmbedError("Failed to retrieve device list", "Could not retrieve your registered devices: " + err.Error()),
				},
			},
		})
		return
	}

	if len(deviceList) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
				Embeds: []*discordgo.MessageEmbed{
					utils.EmbedSuccess("No registered devices", "You have no registered devices for Wake-on-LAN. Use the /register command to add one."),
				},
			},
		})
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "Registered Devices",
		Description: "Here are your registered devices for Wake-on-LAN:",
		Fields: []*discordgo.MessageEmbedField{},
	}

	for _, device := range deviceList {
		field := &discordgo.MessageEmbedField{
			Name: device.Alias,
			Value: device.MacAddress,
			Inline: false,
		}
		embed.Fields = append(embed.Fields, field)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 1 << 6,
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})

}
