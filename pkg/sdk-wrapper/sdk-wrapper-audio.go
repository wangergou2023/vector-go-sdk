package sdk_wrapper

import (
	"errors"
	"fmt"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"os"
	"os/exec"
	"strings"
	"time"
)

const SYSTEMSOUND_WIN = "data/audio/win.pcm"
const VOLUME_LEVEL_MAXIMUM = 5
const VOLUME_LEVEL_MINIMUM = 1

var audioStreamClient vectorpb.ExternalInterface_AudioFeedClient
var audioStreamEnable bool = false

func EnableAudioStream() {
	audioStreamClient, _ = Robot.Conn.AudioFeed(ctx, &vectorpb.AudioFeedRequest{})
	audioStreamEnable = true
}

func DisableAudioStream() {
	audioStreamEnable = false
	audioStreamClient = nil
}

func ProcessAudioStream() {
	// TODO!!!
}

// Returns values in the range 1-5
func GetMasterVolume() int {
	return int(settings["master_volume"].(float64))
}

// Returns values in the range 0-100
func GetAudioVolume() int {
	audioVol := 100 * GetMasterVolume() / VOLUME_LEVEL_MAXIMUM
	return audioVol
}

func SetMasterVolume(volume int) error {
	if volume <= VOLUME_LEVEL_MAXIMUM && volume >= VOLUME_LEVEL_MINIMUM {
		_, err := Robot.Conn.SetMasterVolume(
			ctx,
			&vectorpb.MasterVolumeRequest{
				VolumeLevel: vectorpb.MasterVolumeLevel(volume),
			},
		)
		if err != nil {
			RefreshSDKSettings()
		}
		return err
	}
	return fmt.Errorf("Invalid volume level")
}

// Plays amy sound file (mp3, wav, ecc) using FFMpeg to convert it to the right format

func PlaySound(filename string) string {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		println("File not found!")
		return "failure"
	}

	var pcmFile []byte
	tmpFileName := getTempFilename("pcm", true)
	if strings.Contains(filename, ".pcm") || strings.Contains(filename, ".raw") {
		fmt.Println("Assuming already pcm")
		pcmFile, _ = os.ReadFile(filename)
	} else {
		conOutput, conError := exec.Command("ffmpeg", "-y", "-i", filename, "-f", "s16le", "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1", tmpFileName).Output()
		if conError != nil {
			fmt.Println(conError)
		}
		fmt.Println("FFMPEG output: " + string(conOutput))
		pcmFile, _ = os.ReadFile(tmpFileName)
	}
	var audioChunks [][]byte
	for len(pcmFile) >= 1024 {
		audioChunks = append(audioChunks, pcmFile[:1024])
		pcmFile = pcmFile[1024:]
	}
	var audioClient vectorpb.ExternalInterface_ExternalAudioStreamPlaybackClient
	audioClient, _ = Robot.Conn.ExternalAudioStreamPlayback(
		ctx,
	)
	audioClient.SendMsg(&vectorpb.ExternalAudioStreamRequest{
		AudioRequestType: &vectorpb.ExternalAudioStreamRequest_AudioStreamPrepare{
			AudioStreamPrepare: &vectorpb.ExternalAudioStreamPrepare{
				AudioFrameRate: 16000,
				AudioVolume:    uint32(GetAudioVolume()),
			},
		},
	})
	fmt.Println(len(audioChunks))
	for _, chunk := range audioChunks {
		audioClient.SendMsg(&vectorpb.ExternalAudioStreamRequest{
			AudioRequestType: &vectorpb.ExternalAudioStreamRequest_AudioStreamChunk{
				AudioStreamChunk: &vectorpb.ExternalAudioStreamChunk{
					AudioChunkSizeBytes: 1024,
					AudioChunkSamples:   chunk,
				},
			},
		})
		time.Sleep(time.Millisecond * 30)
	}
	audioClient.SendMsg(&vectorpb.ExternalAudioStreamRequest{
		AudioRequestType: &vectorpb.ExternalAudioStreamRequest_AudioStreamComplete{
			AudioStreamComplete: &vectorpb.ExternalAudioStreamComplete{},
		},
	})
	os.Remove(tmpFileName)

	return "success"
}
