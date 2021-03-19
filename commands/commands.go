package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"randomdogs/dog"
	"randomdogs/util"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "doggo",
			Description: "Sends a random dog that's either a Labrador or a Golden Retriever (untreated poo poo)",
		},
	}
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"doggo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.User != nil {
				return
			}
			dogPicture, err := dog.RetrieveRandomDogPicture()
			if err != nil {
				log.Fatal("Could not retrieve random dog picture", err)
			}
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "A doggo appears at your request!",
							Image: &discordgo.MessageEmbedImage{URL: dogPicture.URL},
							Color: util.GenerateRandomColor(),
							Footer: &discordgo.MessageEmbedFooter{
								Text:    "Requested by " + i.Member.User.Username + "#" + i.Member.User.Discriminator + " " + util.BulletPoint + " " + dogPicture.Breed,
								IconURL: i.Member.User.AvatarURL(""),
							},
						},
					},
				},
			})
			if err != nil {
				log.Fatal("Could not respond to interaction", err)
			}
		},
	}
)

func RegisterCommands(session *discordgo.Session, guildId string) {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handlerFunc, ok := CommandHandlers[i.Data.Name]; ok {
			handlerFunc(s, i)
		}
	})

	for _, v := range Commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, guildId, v)
		if err != nil {
			log.Fatalf("Cannot create command with name %v: %v", v.Name, err)
		}
	}
}
