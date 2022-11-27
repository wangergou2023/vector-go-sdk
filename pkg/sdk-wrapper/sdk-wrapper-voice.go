package sdk_wrapper

import (
	"github.com/bregydoc/gtranslate"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
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

var language string = LANGUAGE_ENGLISH

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
		fName := "TTS-" + GetRobotSerial()
		speech := htgotts.Speech{Folder: "/tmp/audio", Language: language}
		speech.CreateSpeechFile(text, fName)
		PlaySound("/tmp/audio/"+fName+".mp3", volume)
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
		translated = inputLanguage
	}
	return translated
}
