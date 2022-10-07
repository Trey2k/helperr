package discord

import (
	"github.com/Trey2k/helperr/common"
	"github.com/bwmarrin/discordgo"
)

func TestHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Testing....",
		},
	})
	if err != nil {
		common.ErrorLogger.Fatal(err)
	}
}
