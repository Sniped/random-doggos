package dog

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"net/http"
	"os"
	"randomdogs/util"
	"time"
)

var (
	ApiUri = "https://dog.ceo/api"
	// A map of dog breeds mapped by their respective endpoints
	Breeds = map[string]string{
		"Golden Retriever": "retriever/golden",
		"Labrador":         "labrador",
	}
	IntervalRangeStart = time.Hour * 6
	IntervalRangeStop  = time.Hour * 7
)

type APIResponse = struct {
	Message, Status string
}

type Picture = struct {
	URL, Breed string
}

func DoPeriodicDogSend(duration time.Duration, session *discordgo.Session) {
	time.AfterFunc(duration, func() {
		SendRandomDog(session)
	})
}

func SendRandomDog(session *discordgo.Session) {
	dogPicture, err := RetrieveRandomDogPicture()
	if err != nil {
		log.Fatal("Error while retrieving dog picture", err)
	}
	_, err = session.ChannelMessageSendComplex(os.Getenv("DOG_CHANNEL"), &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "A wild doggo has appeared!",
			Image: &discordgo.MessageEmbedImage{URL: dogPicture.URL},
			Color: util.BlueColorHexadecimal,
			Footer: &discordgo.MessageEmbedFooter{
				Text: dogPicture.Breed,
			},
		},
	})
	if err != nil {
		log.Fatal("Could not create dog message", err)
	}
	DoPeriodicDogSend(time.Duration(rand.Int63n(int64(IntervalRangeStop-IntervalRangeStart))+int64(IntervalRangeStart)), session)
}

func GetRandomDogBreed() string {
	dogBreedKeys := make([]string, len(Breeds))
	// add keys to the dogBreedKeys slice
	i := 0
	for k := range Breeds {
		dogBreedKeys[i] = k
		i++
	}
	// choose a random key and return that key
	return dogBreedKeys[rand.Intn(len(Breeds))]
}

func RetrieveRandomDogPicture() (Picture, error) {
	breed := GetRandomDogBreed()
	res, err := http.Get(ApiUri + "/breed/" + Breeds[breed] + "/images/random")
	if err != nil {
		return Picture{}, err
	}
	defer res.Body.Close()
	var body APIResponse
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		log.Fatal("Could not decode JSON response body while fetching dog picture", err)
	}
	return Picture{URL: body.Message, Breed: breed}, nil
}
