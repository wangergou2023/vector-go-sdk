package motors

import (
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
)

func DriveWheelsForward(lw float32, rw float32, lwtwo float32, rwtwo float32) {
	_, _ = sdk_wrapper.Robot.Conn.DriveWheels(
		sdk_wrapper.Ctx,
		&vectorpb.DriveWheelsRequest{
			LeftWheelMmps:   lw,
			RightWheelMmps:  rw,
			LeftWheelMmps2:  lwtwo,
			RightWheelMmps2: rwtwo,
		},
	)
}

func MoveLift(speed float32) {
	_, _ = sdk_wrapper.Robot.Conn.MoveLift(
		sdk_wrapper.Ctx,
		&vectorpb.MoveLiftRequest{
			SpeedRadPerSec: speed,
		},
	)
}

func MoveHead(speed float32) {
	_, _ = sdk_wrapper.Robot.Conn.MoveHead(
		sdk_wrapper.Ctx,
		&vectorpb.MoveHeadRequest{
			SpeedRadPerSec: speed,
		},
	)
}

func DriveOffCharger() {
	_, _ = sdk_wrapper.Robot.Conn.DriveOffCharger(
		sdk_wrapper.Ctx,
		&vectorpb.DriveOffChargerRequest{},
	)
}

func DriveOnCharger() {
	_, _ = sdk_wrapper.Robot.Conn.DriveOnCharger(
		sdk_wrapper.Ctx,
		&vectorpb.DriveOnChargerRequest{},
	)
}
