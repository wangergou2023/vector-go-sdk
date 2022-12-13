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
const TTS_ENGINE_MAX = TTS_ENGINE_ESPEAK

var language string = LANGUAGE_ENGLISH
var eSpeakLang string = "english"
var ttsEngine = TTS_ENGINE_ESPEAK

func InitLanguages(language string) {
	// Here we should get the master volume level...
	SetLanguage(language)
}

func SetLanguage(lang string) {
	language = lang
	eSpeakLang = "english"
	if language == LANGUAGE_ITALIAN {
		eSpeakLang = "italian"
	} else if language == LANGUAGE_SPANISH {
		eSpeakLang = "spanish"
	} else if language == LANGUAGE_FRENCH {
		eSpeakLang = "french"
	} else if language == LANGUAGE_GERMAN {
		eSpeakLang = "german"
	} else if language == LANGUAGE_PORTUGUESE {
		eSpeakLang = "portuguese"
	} else if language == LANGUAGE_DUTCH {
		eSpeakLang = "dutch"
	} else if language == LANGUAGE_RUSSIAN {
		eSpeakLang = "russian"
	} else if language == LANGUAGE_JAPANESE {
		eSpeakLang = "japanese"
	} else if language == LANGUAGE_CHINESE {
		eSpeakLang = "chinese"
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
			speech := htgotts.Speech{Folder: GetTempPath(), Language: language, Handler: &handlers.Native{}}
			speech.CreateSpeechFile(text, fName)
			currentVolume := GetMasterVolume()
			SetMasterVolume(VOLUME_LEVEL_MAXIMUM)
			PlaySound(path.Join(GetTempPath(), fName+".mp3"))
			SetMasterVolume(currentVolume)
			os.Remove(path.Join(GetTempPath(), fName+".mp3"))
		} else {
			// Speex, more robotic. Chinese, Japanese and Russian are not directly supported
			fName := path.Join(GetTempPath(), "TTS-"+GetRobotSerial()+".wav")
			cmdData := "espeak " + "\"" + text + "\"" + " -l " + eSpeakLang + " -w " + fName + " echo 20 75 pitch 82 74"
			_, _, err := shellout(cmdData)
			if err != nil {
				println("ESPEAK ERROR " + err.Error())
			} else {
				PlaySound(fName)
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
