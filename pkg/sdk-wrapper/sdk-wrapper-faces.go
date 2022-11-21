package sdk_wrapper

import "github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"

func FaceEnrollmentListAll() []*vectorpb.LoadedKnownFace {
	response, _ := Robot.Conn.RequestEnrolledNames(
		ctx,
		&vectorpb.RequestEnrolledNamesRequest{},
	)
	return response.GetFaces()
}

func FaceEnrollmentChangeName(faceId int32, oldName string, newName string) string {
	response, _ := Robot.Conn.UpdateEnrolledFaceByID(
		ctx,
		&vectorpb.UpdateEnrolledFaceByIDRequest{
			FaceId:  faceId,
			OldName: oldName,
			NewName: newName,
		},
	)
	return response.Status.String()
}

// Start face enrolling for person
func FaceEnrollmentStart(personName string, id int32) string {
	response, _ := Robot.Conn.SetFaceToEnroll(
		ctx,
		&vectorpb.SetFaceToEnrollRequest{
			Name:        personName,
			ObservedId:  0,
			SaveId:      id,
			SaveToRobot: true,
			SayName:     true,
			UseMusic:    true,
		},
	)
	return response.Status.String()
}

// Cancels operation
func FaceEnrollmentCancel() string {
	response, _ := Robot.Conn.CancelFaceEnrollment(
		ctx,
		&vectorpb.CancelFaceEnrollmentRequest{},
	)
	return response.Status.String()
}

func FaceEnrollmentDeleteAll() string {
	response, _ := Robot.Conn.EraseAllEnrolledFaces(
		ctx,
		&vectorpb.EraseAllEnrolledFacesRequest{},
	)
	return response.Status.String()
}

func FaceEnrollmentDeleteById(id int32) string {
	response, _ := Robot.Conn.EraseEnrolledFaceByID(
		ctx,
		&vectorpb.EraseEnrolledFaceByIDRequest{
			FaceId: id,
		},
	)
	return response.Status.String()
}
