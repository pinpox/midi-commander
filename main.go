package main

import (
	"fmt"
	"os/exec"
	// "strings"
	"time"

	// Midi
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver

	// env
	_ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"io"
	"os"
)

var config map[string]string

func main() {

	deviceName := os.Getenv("MIDI_COMMANDER_DEVICE")
	envConfig := os.Getenv("MIDI_COMMANDER_CONFIG")

	fmt.Println("Using config:", envConfig)

	// Open our jsonFile
	jsonFile, err := os.Open(envConfig)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	jsonStr, _ := io.ReadAll(jsonFile)
	config = map[string]string{}
	json.Unmarshal([]byte(jsonStr), &config)

	defer midi.CloseDriver()

	fmt.Println("Found Midi devices")
	fmt.Println(midi.GetInPorts())

	in, err := midi.FindInPort(deviceName)
	if err != nil {
		fmt.Println("can't find device", deviceName)
		return
	}

	fmt.Println("Using: ", in.String())

	stop, err := midi.ListenTo(in, processMidi, midi.UseSysEx())

	for err == nil {
		time.Sleep(time.Second * 1)
	}

	fmt.Printf("ERROR: %s\n", err)
	stop()
}

func processMidi(msg midi.Message, timestampms int32) {

	var bt []byte
	var ch, key, vel uint8

	switch {
	case msg.GetSysEx(&bt):
		fmt.Printf("got sysex: % X\n SYSEX NOT SUPPORTED YET\n", bt)
	case msg.GetNoteStart(&ch, &key, &vel):
		// fmt.Printf("starting note %s on channel %v with velocity %v\n", midi.Note(key), ch, vel)
		ID := fmt.Sprintf("channel-%v-note-%s", ch, midi.Note(key))

		// fmt.Println("ID is: ", ID)

		val, ok := config[ID]
		if ok {

			cmd := exec.Command("sh", "-c", val)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Start()
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("No command configured for ID:", ID)
		}

	// case msg.GetNoteEnd(&ch, &key):
	// 	fmt.Printf("ending note %s on channel %v\n", midi.Note(key), ch)
	default:
		// ignore
	}
}
