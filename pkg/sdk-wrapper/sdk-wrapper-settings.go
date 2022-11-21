package sdk_wrapper

import (
	"bytes"
	"log"
	"net/http"
	"strings"
)

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

func SetSettingSDKstring(setting string, value string) {
	payload := `{"` + setting + `": "` + value + `" }`
	setSettingSDKStringHelper(payload)
}

/********************************************************************************************************/
/*                                                PRIVATE FUNCTIONS                                     */
/********************************************************************************************************/

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

}
