package audio

import (
	"fmt"
	"testing"
)

func TestAudioToOpusData(t *testing.T) {

	opusData, duration, err := AudioToOpusData("test.mp3", 16000, 1)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("opusData:", opusData)
	fmt.Println("duration:", duration)
}
