package sdk

import (
	"bytes"
	sdk_wrapper "github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/voice"
	"os/exec"
)

const shellToUse = "bash"

func Init(serial string, language string, tmpPath string, dataPath string, nvmPath string) {
	sdk_wrapper.InitSDK(serial)
	sdk_wrapper.SetSDKPaths(tmpPath, dataPath, nvmPath)
	voice.InitLanguages(language)
	sdk_wrapper.RefreshSDKSettings()
}

func InitDefault(serial string) {
	Init(serial, voice.LANGUAGE_ENGLISH, "/tmp/", "data", "nvm/")
}

func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(shellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
