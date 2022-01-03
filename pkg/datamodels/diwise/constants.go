package diwise

const (
	// ExerciseTrailIDPrefix contains the mandatory prefix for Exercise Trail ID
	ExerciseTrailIDPrefix string = "urn:ngsi-ld:" + ExerciseTrailTypeName + ":"
	// ExerciseTrailTypeName contains the Exercise Trail type name to reduce hard coding
	ExerciseTrailTypeName string = "ExerciseTrail"
	// RoadSurfaceObservedIDPrefix contains the mandatory prefix for Road Surfaces Observed
	RoadSurfaceObservedIDPrefix string = "urn:ngsi-ld:RoadSurfaceObserved:"
	// TrailSegmentIDPrefix contains the mandatory prefix for Trail Segment ID
	TrailSegmentIDPrefix string = "urn:ngsi-ld:" + TrailSegmentTypeName + ":"
	// TrailSegmentTypeName contains the Trail Segment type name to reduce hard coding
	TrailSegmentTypeName string = "TrailSegment"
)
