package sdk_wrapper

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/settings"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper/voice"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vector"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type SDKConfigData struct {
	TmpPath  string
	DataPath string
	NvmPath  string
}

var Robot *vector.Vector
var bcAssumption bool = false
var Ctx context.Context
var transCfg = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore SSL warnings
}

var eventStream vectorpb.ExternalInterface_EventStreamClient
var SDKConfig = SDKConfigData{"/tmp/", "data", "nvm"}

func InitSDK(serial string) {
	var err error
	voice.InitLanguages(voice.LANGUAGE_ENGLISH)
	Robot, err = vector.NewEP(serial)
	if err != nil {
		log.Fatal(err)
	}
	Ctx = context.Background()
	eventStream, err = Robot.Conn.EventStream(Ctx, &vectorpb.EventRequest{})
	settings.RefreshSDKSettings()
}

func SetNDKPaths(tmpPath string, dataPath string, nvmPath string) {
	SDKConfig.TmpPath = tmpPath
	SDKConfig.DataPath = dataPath
	SDKConfig.NvmPath = nvmPath
}

func WaitForEvent() *vectorpb.Event {
	evt, err := eventStream.Recv()
	if err != nil {
		return nil
	}
	return evt.GetEvent()
}

func AssumeBehaviorControl(priority string) {
	var controlRequest *vectorpb.BehaviorControlRequest
	if priority == "high" {
		controlRequest = &vectorpb.BehaviorControlRequest{
			RequestType: &vectorpb.BehaviorControlRequest_ControlRequest{
				ControlRequest: &vectorpb.ControlRequest{
					Priority: vectorpb.ControlRequest_OVERRIDE_BEHAVIORS,
				},
			},
		}
	} else {
		controlRequest = &vectorpb.BehaviorControlRequest{
			RequestType: &vectorpb.BehaviorControlRequest_ControlRequest{
				ControlRequest: &vectorpb.ControlRequest{
					Priority: vectorpb.ControlRequest_DEFAULT,
				},
			},
		}
	}
	go func() {
		start := make(chan bool)
		stop := make(chan bool)
		bcAssumption = true
		go func() {
			// * begin - modified from official vector-go-sdk
			r, err := Robot.Conn.BehaviorControl(
				Ctx,
			)
			if err != nil {
				log.Println(err)
				return
			}

			if err := r.Send(controlRequest); err != nil {
				log.Println(err)
				return
			}

			for {
				ctrlresp, err := r.Recv()
				if err != nil {
					log.Println(err)
					return
				}
				if ctrlresp.GetControlGrantedResponse() != nil {
					start <- true
					break
				}
			}

			for {
				select {
				case <-stop:
					if err := r.Send(
						&vectorpb.BehaviorControlRequest{
							RequestType: &vectorpb.BehaviorControlRequest_ControlRelease{
								ControlRelease: &vectorpb.ControlRelease{},
							},
						},
					); err != nil {
						log.Println(err)
						return
					}
					return
				default:
					continue
				}
			}
			// * end - modified from official vector-go-sdk
		}()
		for {
			select {
			case <-start:
				for {
					if bcAssumption {
						time.Sleep(time.Millisecond * 500)
					} else {
						break
					}
				}
				stop <- true
				return
			}
		}
	}()
}

func ReleaseBehaviorControl() {
	bcAssumption = false
}

func GetRobotSerial() string {
	return Robot.Cfg.SerialNo
}

func GetTemporaryFilename(tag string, extension string, fullpath bool) string {
	tmpPath := SDKConfig.TmpPath
	tmpFile := GetRobotSerial() + "_" + tag + fmt.Sprintf("_%d", time.Now().Unix()) + "." + extension
	if fullpath {
		tmpFile = tmpPath + tmpFile
	}
	return tmpFile
}

func GetMyStoragePath(filename string) string {
	nvmPath := SDKConfig.NvmPath
	nvmPath = filepath.Join(nvmPath, GetRobotSerial())
	_ = os.MkdirAll(nvmPath, os.ModePerm)
	nvmPath = filepath.Join(nvmPath, filename)
	return nvmPath
}

func GetDataPath(filename string) string {
	dataPath := SDKConfig.DataPath
	var chunks []string = strings.Split(filename, "/")
	for _, chunk := range chunks {
		dataPath = filepath.Join(dataPath, chunk)
	}
	return dataPath
}

/**********************************************************************************************************************/
/*                                              PRIVATE FUNCTIONS                                                     */
/**********************************************************************************************************************/

const shellToUse = "bash"

func shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(shellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
