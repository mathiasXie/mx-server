package main

import (
	"fmt"
	"os"

	audio_utils "github.com/mathiasXie/gin-web/utils/audio/utils"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run mp3_to_opus.go <input_mp3_file>")
		return
	}

	inputFile := os.Args[1]

	opusData, duration, err := audio_utils.AudioToOpusData(inputFile, 16000, 1)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("opusData:", opusData)
	fmt.Println("duration:", duration)
}
