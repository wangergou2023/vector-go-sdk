package sdk_wrapper

import (
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
)

// Displays animation

func PlayAnimation(anim string, loops uint32, ignoreBodyTrack bool, ignoreHeadTrack bool, ignoreLiftTrack bool) {
	var a = &vectorpb.Animation{
		Name: anim, // ignore SSL warnings
	}
	_, _ = Robot.Conn.PlayAnimation(
		ctx,
		&vectorpb.PlayAnimationRequest{
			Animation:       a,
			Loops:           loops,
			IgnoreBodyTrack: ignoreBodyTrack,
			IgnoreHeadTrack: ignoreHeadTrack,
			IgnoreLiftTrack: ignoreLiftTrack,
		},
	)
}
