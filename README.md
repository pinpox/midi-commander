# midi-commander

Execute shell scripts when a MIDI note is played.


## Configuration

The following environment variables are expected to be set:

| Variable                | Description                   |
|-------------------------|-------------------------------|
| `MIDI_COMMANDER_DEVICE` | Name of the device to be used |
| `MIDI_COMMANDER_CONFIG` | Path to JSON config           |

The config has the following format:

```json
{
    "channel-0-note-C4": "./test-script1.sh"
    "channel-0-note-D4": "./test-script2.sh"
}
```
