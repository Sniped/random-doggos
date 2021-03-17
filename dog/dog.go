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

type DogAPIResponse = struct {
	Message, Status string
}

func DoPeriodicDogSend(duration time.Duration) {
	time.AfterFunc(duration, func() {
		SendRandomDog(true)
	})
}

func SendRandomDog(periodic bool) {
	dogPictureUrl, err := RetrieveRandomDogPicture()
	if err != nil {
		log.Fatal("Error while retrieving dog picture", err)
	}
	_, err = util.Session.ChannelMessageSendComplex(os.Getenv("DOG_CHANNEL"), &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "A wild doggo has appeared!",
			Image: &discordgo.MessageEmbedImage{URL: dogPictureUrl},
			Color: util.BlueColorHexadecimal,
		},
	})
	if err != nil {
		log.Fatal("Could not create dog message", err)
	}
	if periodic {
		DoPeriodicDogSend(time.Duration(rand.Int63n(int64(IntervalRangeStop-IntervalRangeStart)) + int64(IntervalRangeStart)))
	}
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

func RetrieveRandomDogPicture() (string, error) {
	res, err := http.Get(ApiUri + "/breed/" + Breeds[GetRandomDogBreed()] + "/images/random")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	var body DogAPIResponse
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		log.Fatal("Could not decode JSON response body while fetching dog picture", err)
	}
	return body.Message, nil
}
