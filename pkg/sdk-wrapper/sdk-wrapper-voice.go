/*
  To run this on a raspberry pi:
  sudo apt-get install gcc libasound2 libasound2-dev
  sudo apt install espeak -y
*/

package sdk_wrapper

import (
	"bytes"
	"github.com/bregydoc/gtranslate"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	"os"
	"os/exec"
)

const LANGUAGE_ENGLISH = voices.English
const LANGUAGE_ITALIAN = voices.Italian
const LANGUAGE_SPANISH = voices.Spanish
const LANGUAGE_FRENCH = voices.French
const LANGUAGE_GERMAN = voices.German
const LANGUAGE_PORTUGUESE = voices.Portuguese
const LANGUAGE_DUTCH = voices.Dutch
const LANGUAGE_RUSSIAN = voices.Russian
const LANGUAGE_JAPANESE = voices.Japanese
const LANGUAGE_CHINESE = voices.Chinese

const TTS_ENGINE_HTGO = 0
const TTS_ENGINE_ESPEAK = 1
const TTS_ENGINE_MAX = TTS_ENGINE_ESPEAK

var language string = LANGUAGE_ENGLISH
var ttsEngine = TTS_ENGINE_ESPEAK

var volume = 100

func Init() {
	// Here we should get the master volume level...
}

func SetLanguage(lang string) {
	language = lang
}

func GetLanguage() string {
	return language
}

func SetTTSEngine(TTSEngine int) {
	if TTSEngine >= 0 && TTSEngine < TTS_ENGINE_MAX {
		ttsEngine = TTSEngine
	}
}

func sayText(text string) {
	_, _ = Robot.Conn.SayText(
		ctx,
		&vectorpb.SayTextRequest{
			Text:           text,
			UseVectorVoice: true,
			DurationScalar: 1.0,
		},
	)
}

func SayText(text string) {
	useNativeTTS := true
	if language == LANGUAGE_ITALIAN ||
		language == LANGUAGE_SPANISH ||
		language == LANGUAGE_FRENCH ||
		language == LANGUAGE_GERMAN ||
		language == LANGUAGE_PORTUGUESE ||
		language == LANGUAGE_DUTCH ||
		language == LANGUAGE_RUSSIAN ||
		language == LANGUAGE_JAPANESE ||
		language == LANGUAGE_CHINESE {
		useNativeTTS = false
	}

	if useNativeTTS {
		sayText(text)
	} else {
		//println("Saying " + text + " in language=" + language)
		if ttsEngine == TTS_ENGINE_HTGO {
			// Uses Google voices
			fName := "TTS-" + GetRobotSerial()
			speech := htgotts.Speech{Folder: "/tmp", Language: language, Handler: &handlers.Native{}}
			speech.CreateSpeechFile(text, fName)
			PlaySound("/tmp/"+fName+".mp3", volume)
			os.Remove("/tmp/" + fName + ".mp3")
		} else {
			// Speex, more robotic
			fName := "/tmp/TTS-" + GetRobotSerial() + ".wav"
			cmdData := "\"" + text + "\"" + " -v italian -w " + fName + " echo 20 75 pitch 82 74"
			println(cmdData)
			cmd := exec.Command("espeak", cmdData)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			println(out.String())
			if err != nil {
				println("ESPEAK ERROR " + err.Error())
			} else {
				PlaySound(fName, volume)
				os.Remove(fName)
			}
		}
	}
}

func Translate(text string, inputLanguage string, outputLanguage string) string {
	translated, err := gtranslate.TranslateWithParams(
		text,
		gtranslate.TranslationParams{
			From: inputLanguage,
			To:   outputLanguage,
		},
	)
	if nil != err {
		println("GTRANSLATE ERROR " + err.Error())
		translated = inputLanguage
	}
	return translated
}
