package main

import (
	"fmt"
	"os/exec"
	"strings"
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

func runCmdForID(ID string, values []string) {

	cmdExe, ok := config[ID]
	param := fmt.Sprintf("%s %s", cmdExe, strings.Join(values, " "))

	if ok {

		cmd := exec.Command("/bin/sh", "-c", param)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("No command configured for ID:", ID)
	}

}

func processMidi(msg midi.Message, timestampms int32) {

	var bt []byte
	var ch, key, controller, value, vel uint8

	switch {
	case msg.GetControlChange(&ch, &controller, &value):
		ID := fmt.Sprintf("cc-controller-%v-channel-%v", controller, ch)
		runCmdForID(ID, []string{
			fmt.Sprintf("%v", ch),
			fmt.Sprintf("%v", controller),
			fmt.Sprintf("%v", value),
		})
	case msg.GetSysEx(&bt):
		fmt.Printf("got sysex: % X\n SYSEX NOT SUPPORTED YET\n", bt)
	case msg.GetNoteStart(&ch, &key, &vel):
		ID := fmt.Sprintf("midi-channel-%v-note-%s", ch, midi.Note(key))
		runCmdForID(ID,
			[]string{
				fmt.Sprintf("%v", ch),
				fmt.Sprintf("%s", midi.Note(key)),
				fmt.Sprintf("%v", vel)})
	default:
		//ignore
		// fmt.Printf("got unknow message: %v ", msg)
	}
}
