package sdk_wrapper

import (
	"bytes"
	"encoding/json"
	"github.com/PerformLine/go-stockutil/colorutil"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	"image/color"
	"log"
	"net/http"
	"os"
	"strings"
)

type CustomSettings struct {
	robot_name string
}

var settings map[string]interface{}
var defaultCustomSettings = CustomSettings{""}
var customSettings = defaultCustomSettings

func RefreshSDKSettings() {
	settingsJSON := getSDKSettings()
	customSettingsJSON, err := getCustomSettings()
	if err != nil {
		customSettings = defaultCustomSettings
	}
	//println(string(settingsJSON))
	json.Unmarshal([]byte(settingsJSON), &settings)
	json.Unmarshal([]byte(customSettingsJSON), &customSettings)
}

func GetEyeColor() color.RGBA {
	eyeColor := color.RGBA{0, 0, 0, 0xff}

	presetEyeColor := settings["eye_color"].(float64)
	var customEyeColor map[string]interface{} = (settings["custom_eye_color"]).(map[string]interface{})
	customEyeColorEnabled := customEyeColor["enabled"].(bool)

	if customEyeColorEnabled {
		customEyeColorHue := customEyeColor["hue"].(float64)
		customEyeColorSaturation := customEyeColor["saturation"].(float64)
		//println(fmt.Sprintf("HUE/SAT: %0f,%f", customEyeColorHue, customEyeColorSaturation))

		eyeColor.R, eyeColor.G, eyeColor.B = colorutil.HslToRgb(customEyeColorHue, customEyeColorSaturation, 1.0)
	} else {
		switch presetEyeColor {
		case float64(vectorpb.EyeColor_value["TIP_OVER_TEAL"]):
			eyeColor = color.RGBA{0x29, 0xae, 0x70, 0xff}
			break
		case float64(vectorpb.EyeColor_value["OVERFIT_ORANGE"]):
			eyeColor = color.RGBA{0xfe, 0x76, 0x14, 0xff}
			break
		case float64(vectorpb.EyeColor_value["UNCANNY_YELLOW"]):
			eyeColor = color.RGBA{0xf7, 0xcb, 0x04, 0xff}
			break
		case float64(vectorpb.EyeColor_value["NON_LINEAR_LIME"]):
			eyeColor = color.RGBA{0xa8, 0xd3, 0x04, 0xff}
			break
		case float64(vectorpb.EyeColor_value["SINGULARITY_SAPPHIRE"]):
			eyeColor = color.RGBA{0x0d, 0x97, 0xf0, 0xff}
			break
		case float64(vectorpb.EyeColor_value["FALSE_POSITIVE_PURPLE"]):
			eyeColor = color.RGBA{0xc6, 0x61, 0xfc, 0xff}
			break
		case float64(vectorpb.EyeColor_value["CONFUSION_MATRIX_GREEN"]):
			// Color unknown
			eyeColor = color.RGBA{0x00, 0xff, 0x00, 0xff}
			break
		}
	}
	//println(fmt.Sprintf("%02x,%02x,%02x", eyeColor.R, eyeColor.G, eyeColor.B))
	return eyeColor
}

func GetTemperatureUnit() string {
	unit := WEATHER_UNIT_CELSIUS
	isFaranheit := settings["temp_is_fahrenheit"].(bool)
	if isFaranheit {
		unit = WEATHER_UNIT_FARANHEIT
	}
	return unit
}

func SetCustomEyeColor(hue string, sat string) {
	payload := `{"custom_eye_color": {"enabled": true, "hue": ` + hue + `, "saturation": ` + sat + `} }`
	setSettingSDKStringHelper(payload)
}

func SetPresetEyeColor(value string) {
	payload := `{"custom_eye_color": {"enabled": false}, "eye_color": ` + value + `}`
	setSettingSDKStringHelper(payload)
}

func SetLocale(locale string) {
	SetSettingSDKstring("locale", locale)
}

func SetLocation(location string) {
	SetSettingSDKstring("default_location", location)
}

func SetTimeZone(timezone string) {
	SetSettingSDKstring("time_zone", timezone)
}

func SetRobotName(name string) {
	customSettings.robot_name = name
	saveCustomSettings()
}

func SetSettingSDKstring(setting string, value string) {
	payload := `{"` + setting + `": "` + value + `" }`
	setSettingSDKStringHelper(payload)
	RefreshSDKSettings()
}

/********************************************************************************************************/
/*                                                PRIVATE FUNCTIONS                                     */
/********************************************************************************************************/

func getSDKSettings() []byte {
	resp, err := Robot.Conn.PullJdocs(ctx, &vectorpb.PullJdocsRequest{
		JdocTypes: []vectorpb.JdocType{vectorpb.JdocType_ROBOT_SETTINGS},
	})
	if err != nil {
		return []byte(err.Error())
	}
	json := resp.NamedJdocs[0].Doc.JsonDoc
	return []byte(json)
}

func setSettingSDKStringHelper(payload string) {
	if !strings.Contains(Robot.Cfg.Token, "error") {
		url := "https://" + Robot.Cfg.Target + "/v1/update_settings"
		var updateJSON = []byte(`{"update_settings": true, "settings": ` + payload + ` }`)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(updateJSON))
		req.Header.Set("Authorization", "Bearer "+Robot.Cfg.Token)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Transport: transCfg}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	} else {
		log.Println("GUID not there")
	}
	RefreshSDKSettings()
}

func getCustomSettings() ([]byte, error) {
	json, err := os.ReadFile(GetMyStoragePath("custom_settings.json"))
	if err == nil {
		return nil, err
	}
	return []byte(json), nil
}

func saveCustomSettings() {
	file, _ := json.MarshalIndent(customSettings, "", " ")
	_ = os.WriteFile(GetMyStoragePath("custom_settings.json"), file, 0644)
}
