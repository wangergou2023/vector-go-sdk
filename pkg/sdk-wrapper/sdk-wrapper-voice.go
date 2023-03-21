/*
  To run this on a raspberry pi:
  sudo apt-get install gcc libasound2 libasound2-dev
  sudo apt install espeak -y
*/

package sdk_wrapper

import (
	"github.com/bregydoc/gtranslate"
	"github.com/fforchino/vector-go-sdk/pkg/vectorpb"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	"net/url"
	"os"
	"path"
	"strings"
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
const TTS_ENGINE_VOICESERVER = 2
const TTS_ENGINE_MAX = TTS_ENGINE_VOICESERVER

var language string = LANGUAGE_ENGLISH
var eSpeakLang string = "en"
var ttsEngine = TTS_ENGINE_VOICESERVER

func InitLanguages(language string) {
	// Here we should get the master volume level...
	SetLanguage(language)
}

func SetLanguage(lang string) {
	language = lang
	eSpeakLang = "en"
	if language == LANGUAGE_ITALIAN {
		eSpeakLang = "it"
	} else if language == LANGUAGE_SPANISH {
		eSpeakLang = "es"
	} else if language == LANGUAGE_FRENCH {
		eSpeakLang = "fr"
	} else if language == LANGUAGE_GERMAN {
		eSpeakLang = "de"
	} else if language == LANGUAGE_PORTUGUESE {
		eSpeakLang = "pt"
	} else if language == LANGUAGE_DUTCH {
		eSpeakLang = "nl"
	} else if language == LANGUAGE_RUSSIAN {
		eSpeakLang = "ru"
	} else if language == LANGUAGE_JAPANESE {
		eSpeakLang = "jp"
	} else if language == LANGUAGE_CHINESE {
		eSpeakLang = "zh"
	}
}

func GetLanguage() string {
	return language
}

func GetLanguageISO2() string {
	locale := GetLocale()
	return strings.Split(locale, "-")[0]
}

func SetTTSEngine(TTSEngine int) {
	if TTSEngine >= 0 && TTSEngine <= TTS_ENGINE_MAX {
		ttsEngine = TTSEngine
	}
}

func SayText(text string) {
	useNativeTTS := false // Use web TTS also for English
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
			speech := htgotts.Speech{Folder: GetTempPath(), Language: language, Handler: &handlers.Native{}}
			speech.CreateSpeechFile(text, fName)
			currentVolume := GetMasterVolume()
			SetMasterVolume(VOLUME_LEVEL_MAXIMUM)
			PlaySound(path.Join(GetTempPath(), fName+".mp3"))
			SetMasterVolume(currentVolume)
			os.Remove(path.Join(GetTempPath(), fName+".mp3"))
		} else if ttsEngine == TTS_ENGINE_ESPEAK {
			// Speex, more robotic. Chinese, Japanese and Russian are not directly supported
			fName := path.Join(GetTempPath(), "TTS-"+GetRobotSerial()+".wav")
			cmdData := "espeak " + "\"" + text + "\"" + " -l " + eSpeakLang + " -w " + fName + " echo 20 75 pitch 82 74"
			//println(cmdData)
			_, _, err := shellout(cmdData)
			if err != nil {
				println("ESPEAK ERROR " + err.Error())
			} else {
				PlaySound(fName)
				os.Remove(fName)
			}
		} else if ttsEngine == TTS_ENGINE_VOICESERVER {
			// Uses FakeYou voices
			fName := path.Join(GetTempPath(), "TTS-"+GetRobotSerial()+".wav")

			// Yoda
			theUrl := "https://www.wondergarden.app/voiceserver/index.php/getText?text=" + url.QueryEscape(text) + "&lang=en-US&voice=TM:fmspb239ea3a"

			println("Will save file " + fName)
			println("Request url: " + theUrl)
			err, fWebName := DownloadFile(fName, theUrl)
			if err == nil {
				println("Name from web: " + fWebName)
				if strings.HasSuffix(fWebName, ".mp3") {
					os.Rename(fName, path.Join(GetTempPath(), "TTS-"+GetRobotSerial()+".mp3"))
				}
				currentVolume := GetMasterVolume()
				SetMasterVolume(VOLUME_LEVEL_MAXIMUM)
				PlaySound(fName)
				SetMasterVolume(currentVolume)
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

/**********************************************************************************************************************/
/*                                              PRIVATE FUNCTIONS                                                     */
/**********************************************************************************************************************/

func refreshLanguage() {
	loc := GetLocale()
	if strings.HasPrefix(loc, "it") {
		SetLanguage(LANGUAGE_ITALIAN)
	} else if strings.HasPrefix(loc, "es") {
		SetLanguage(LANGUAGE_SPANISH)
	} else if strings.HasPrefix(loc, "fr") {
		SetLanguage(LANGUAGE_FRENCH)
	} else if strings.HasPrefix(loc, "de") {
		SetLanguage(LANGUAGE_GERMAN)
	} else {
		SetLanguage(LANGUAGE_ENGLISH)
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
