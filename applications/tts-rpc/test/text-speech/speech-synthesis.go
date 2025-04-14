package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func synthesizeStartedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Println("Synthesis started.")
}

func synthesizingHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Printf("Synthesizing, audio chunk size %d.\n", len(event.Result.AudioData))
}

func synthesizedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Printf("Synthesized, audio length %d.\n", len(event.Result.AudioData))
}

func cancelledHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Println("Received a cancellation.")
}

func main() {
	// This example requires environment variables named "SPEECH_KEY" and "SPEECH_REGION"
	speechKey := "8aKNEoi9njCkkJSZIKX4mHSXDCBl8HEtRlrP4ev4Zuw7X4J93CL2JQQJ99BCACqBBLyXJ3w3AAAYACOGJKg4"
	speechRegion := "southeastasia"

	audioConfig, err := audio.NewAudioConfigFromDefaultSpeakerOutput()
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer audioConfig.Close()
	speechConfig, err := speech.NewSpeechConfigFromSubscription(speechKey, speechRegion)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer speechConfig.Close()
	speechConfig.SetSpeechSynthesisOutputFormat(common.Webm24Khz16Bit24KbpsMonoOpus)

	speechConfig.SetSpeechSynthesisVoiceName("en-US-AvaMultilingualNeural")

	speechSynthesizer, err := speech.NewSpeechSynthesizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return
	}
	defer speechSynthesizer.Close()

	speechSynthesizer.SynthesisStarted(synthesizeStartedHandler)
	speechSynthesizer.Synthesizing(synthesizingHandler)
	speechSynthesizer.SynthesisCompleted(synthesizedHandler)
	speechSynthesizer.SynthesisCanceled(cancelledHandler)

	for {
		fmt.Printf("Enter some text that you want to speak, or enter empty text to exit.\n> ")
		text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		if len(text) == 0 {
			break
		}

		// StartSpeakingTextAsync sends the result to channel when the synthesis starts.
		task := speechSynthesizer.StartSpeakingTextAsync(text)
		var outcome speech.SpeechSynthesisOutcome
		select {
		case outcome = <-task:
		case <-time.After(60 * time.Second):
			fmt.Println("Timed out")
			return
		}
		defer outcome.Close()
		if outcome.Error != nil {
			fmt.Println("Got an error: ", outcome.Error)
			return
		}

		// in most case we want to streaming receive the audio to lower the latency,
		// we can use AudioDataStream to do so.
		stream, err := speech.NewAudioDataStreamFromSpeechSynthesisResult(outcome.Result)
		defer stream.Close()
		if err != nil {
			fmt.Println("Got an error: ", err)
			return
		}

		var all_audio []byte
		audio_chunk := make([]byte, 2048)
		for {
			n, err := stream.Read(audio_chunk)

			if err == io.EOF {
				break
			}

			all_audio = append(all_audio, audio_chunk[:n]...)
		}

		// 将音频数据保存到文件
		filePath := "output_test.wav"
		err = os.WriteFile(filePath, all_audio, 0644)
		if err != nil {
			fmt.Println("Got an error: ", err)
		}

		fmt.Printf("Read [%d] bytes from audio data stream.\n", len(all_audio))
	}
}
