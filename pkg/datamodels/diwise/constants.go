package diwise

const urnPrefix string = "urn:ngsi-ld:"

const (
	// ExerciseTrailIDPrefix contains the mandatory prefix for Exercise Trail ID
	ExerciseTrailIDPrefix string = urnPrefix + ExerciseTrailTypeName + ":"
	// ExerciseTrailTypeName is a type name constant for ExerciseTrail
	ExerciseTrailTypeName string = "ExerciseTrail"
	//RoadSurfaceObservedTypeName is a type name constant for RoadSurfaceObserved
	RoadSurfaceObservedTypeName string = "RoadSurfaceObserved"
	// RoadSurfaceObservedIDPrefix contains the mandatory prefix for Road Surfaces Observed
	RoadSurfaceObservedIDPrefix string = urnPrefix + RoadSurfaceObservedTypeName + ":"
	// TrailSegmentIDPrefix contains the mandatory prefix for Trail Segment ID
	TrailSegmentIDPrefix string = urnPrefix + TrailSegmentTypeName + ":"
	// TrailSegmentTypeName is a type name constant for TrailSegment
	TrailSegmentTypeName string = "TrailSegment"
)
