package sdk_wrapper

import (
	"errors"
	"fmt"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const SYSTEMSOUND_WIN = "data/audio/win.pcm"

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

// Plays amy sound file (mp3, wav, ecc) using FFMpeg to convert it to the right format

func PlaySound(filename string, volume int) string {
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
				AudioVolume:    uint32(volume),
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

// Sets master volume (1 to 5)

func SetMasterVolume(volume int) bool {
	ret := false
	if volume >= 1 && volume <= 5 {
		strVol := strconv.Itoa(volume)
		SetSettingSDKstring("master_volume", strVol)
		ret = true
	}
	return ret
}
