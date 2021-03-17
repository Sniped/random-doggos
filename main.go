package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"randomdogs/commands"
	"randomdogs/dog"
)

var session *discordgo.Session

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	session, err = discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))

	session.AddHandler(func(s *discordgo.Session, e *discordgo.Ready) {
		log.Print("Bot is ready on " + session.State.User.Username + "#" + session.State.User.Discriminator)
	})

	err = session.Open()
	if err != nil {
		log.Fatal("Cannot open discord bot session", err)
	}

	guildId := os.Getenv("GUILD_ID")
	if guildId == "" {
		log.Fatal("No guild ID specified in env")
	}

	dogChannelId := os.Getenv("DOG_CHANNEL")
	if dogChannelId == "" {
		log.Fatal("There is no dog channel ID specified in env")
	}

	commands.RegisterCommands(session, guildId)

	dog.DoPeriodicDogSend(dog.IntervalRangeStart, session)

	defer session.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
