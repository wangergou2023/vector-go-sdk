package faces

import (
	"github.com/digital-dream-labs/vector-go-sdk/pkg/sdk-wrapper"
	"github.com/digital-dream-labs/vector-go-sdk/pkg/vectorpb"
)

func FaceEnrollmentListAll() []*vectorpb.LoadedKnownFace {
	response, _ := sdk_wrapper.Robot.Conn.RequestEnrolledNames(
		sdk_wrapper.Ctx,
		&vectorpb.RequestEnrolledNamesRequest{},
	)
	return response.GetFaces()
}

func FaceEnrollmentChangeName(faceId int32, oldName string, newName string) string {
	response, _ := sdk_wrapper.Robot.Conn.UpdateEnrolledFaceByID(
		sdk_wrapper.Ctx,
		&vectorpb.UpdateEnrolledFaceByIDRequest{
			FaceId:  faceId,
			OldName: oldName,
			NewName: newName,
		},
	)
	return response.Status.String()
}

// Start face enrolling for person with the given name
// It doesn't seem to work, the face seems enrolled but not saved

func FaceEnrollmentStart(personName string) string {
	faces := FaceEnrollmentListAll()
	var maxId int32 = 0

	for i := 0; i < len(faces); i++ {
		if faces[i].FaceId > maxId {
			maxId = faces[i].FaceId
		}
	}
	maxId++

	response, _ := sdk_wrapper.Robot.Conn.SetFaceToEnroll(
		sdk_wrapper.Ctx,
		&vectorpb.SetFaceToEnrollRequest{
			Name:        personName,
			ObservedId:  0,
			SaveId:      maxId,
			SaveToRobot: true,
			SayName:     true,
			UseMusic:    true,
		},
	)

	return response.Status.String()
}

// Cancels operation
func FaceEnrollmentCancel() string {
	response, _ := sdk_wrapper.Robot.Conn.CancelFaceEnrollment(
		sdk_wrapper.Ctx,
		&vectorpb.CancelFaceEnrollmentRequest{},
	)
	return response.Status.String()
}

func FaceEnrollmentDeleteAll() string {
	response, _ := sdk_wrapper.Robot.Conn.EraseAllEnrolledFaces(
		sdk_wrapper.Ctx,
		&vectorpb.EraseAllEnrolledFacesRequest{},
	)
	return response.Status.String()
}

func FaceEnrollmentDeleteById(id int32) string {
	response, _ := sdk_wrapper.Robot.Conn.EraseEnrolledFaceByID(
		sdk_wrapper.Ctx,
		&vectorpb.EraseEnrolledFaceByIDRequest{
			FaceId: id,
		},
	)
	return response.Status.String()
}
