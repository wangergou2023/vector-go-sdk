package vector

import (
	"context"

	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
)

// BehaviorControl assumes control of the vector robot for SDK usage.  Once control is
// assumed, a signal is sent on the start channel. To give control back to the bot, send
// a message to the stop channel.  Failing to do so may leave your bot in a funny, funny
// state.
func (v *Vector) BehaviorControl(ctx context.Context, start, stop chan bool) error {
	print("0...")
	r, err := v.Conn.BehaviorControl(
		ctx,
	)
	if err != nil {
		return err
	}

	print("1...")
	if err := r.Send(
		&vectorpb.BehaviorControlRequest{
			RequestType: &vectorpb.BehaviorControlRequest_ControlRequest{
				ControlRequest: &vectorpb.ControlRequest{
					Priority: vectorpb.ControlRequest_DEFAULT,
				},
			},
		},
	); err != nil {
		return err
	}

	print("SENT")
	for {
		ctrlresp, err := r.Recv()
		if err != nil {
			return err
		}
		print("RECEIVED")
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
				return err
			}
			return nil
		default:
			continue
		}
	}

}
