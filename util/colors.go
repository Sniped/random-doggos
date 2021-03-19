package util

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
)

func GenerateRandomColor() int {
	rgbSlice := make([]string, 3)
	for i := range rgbSlice {
		rgbSlice[i] = strconv.FormatInt(int64(rand.Intn(255)), 16)
	}
	hex, err := strconv.ParseInt(strings.Join(rgbSlice, ""), 16, 32)
	if err != nil {
		log.Fatalf("Error while generating random hex %v", err)
	}
	return int(hex)
}
