package ngsi

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestGeoCoordinatesParser(t *testing.T) {
	coords, err := parseGeometryCoordinates("[[2.4,2.1],[3.3,3.7]]")
	if err != nil {
		t.Error("Got error from coordinates parser", err)
	}

	if len(coords) != 4 {
		t.Error("Expected 2 coordinates, got", len(coords)/2)
	} else {
		if coords[0] != 2.4 || coords[1] != 2.1 {
			t.Error(fmt.Sprintf("First position should be (2.4,2.1) and not (%f,%f)", coords[0], coords[1]))
		}
	}
}

func TestCreateQueryFromParameters(t *testing.T) {
	is := is.New(t)

	req, _ := http.NewRequest("GET", createURL("/entities", "limit=2", "offset=5"), nil)

	query, err := newQueryFromParameters(req, []string{"T"}, []string{"a"}, "")
	is.NoErr(err)                                 // newQueryFromParameters failed
	is.Equal(query.PaginationLimit(), uint64(2))  // failed to parse correct pagination limit
	is.Equal(query.PaginationOffset(), uint64(5)) // failed to parse correct pagination offset
}

const (
	DateTimeAt    string = "2017-12-13T14:20:00Z"
	DateEndTimeAt string = "2017-12-13T14:40:00Z"
)

func TestCreateTemporalQueryAfterTime(t *testing.T) {
	is := is.New(t)

	timeAt, _ := time.Parse(time.RFC3339, DateTimeAt)
	req, _ := http.NewRequest("GET", createURL("/entities",
		fmt.Sprintf("timerel=%s", TemporalRelationAfterTime),
		fmt.Sprintf("timeAt=%s", DateTimeAt)),
		nil)

	query, err := newQueryFromParameters(req, []string{"T"}, []string{"a"}, "")
	is.NoErr(err) // newQueryFromParameters failed
	is.True(query.IsTemporalQuery())

	from, to := query.Temporal().TimeSpan()
	is.Equal(from, timeAt) // start of timespan should equal timeAt
	is.True(to.IsZero())   // end of timespan should not be set
}

func TestCreateTemporalQueryUsesCorrectDefaultTemporalProperty(t *testing.T) {
	is := is.New(t)

	req, _ := http.NewRequest("GET", createURL("/entities",
		fmt.Sprintf("timerel=%s", TemporalRelationAfterTime),
		fmt.Sprintf("timeAt=%s", DateTimeAt)),
		nil)

	query, err := newQueryFromParameters(req, []string{"T"}, []string{"a"}, "")
	is.NoErr(err) // newQueryFromParameters failed
	is.Equal(query.Temporal().Property(), "observedAt")
}

func TestCreateTemporalQueryAfterTimeIgnoresEndTime(t *testing.T) {
	is := is.New(t)

	req, _ := http.NewRequest("GET", createURL("/entities",
		fmt.Sprintf("timerel=%s", TemporalRelationAfterTime),
		fmt.Sprintf("timeAt=%s", DateTimeAt),
		fmt.Sprintf("endTimeAt=%s", DateEndTimeAt)),
		nil)

	query, err := newQueryFromParameters(req, []string{"T"}, []string{"a"}, "")
	is.NoErr(err) // newQueryFromParameters failed
	is.True(query.IsTemporalQuery())

	_, to := query.Temporal().TimeSpan()
	is.True(to.IsZero()) // end time should not be set
}

func TestCreateTemporalQueryReturnsErrorOnNonSupportedTimerel(t *testing.T) {
	is := is.New(t)

	req, _ := http.NewRequest("GET", createURL("/entities", "timerel=notsupported"), nil)

	_, err := newQueryFromParameters(req, []string{"T"}, []string{"a"}, "")
	is.True(err != nil) // should return an error
}

func TestCreateTemporalQueryReturnsErrorOnMissingTimeAt(t *testing.T) {
	is := is.New(t)

	req, _ := http.NewRequest("GET", createURL("/entities",
		fmt.Sprintf("timerel=%s", TemporalRelationAfterTime)),
		nil)

	_, err := newQueryFromParameters(req, []string{"T"}, []string{"a"}, "")
	is.True(err != nil) // should return an error
}

func TestCreateTemporalQueryReturnsErrorOnInvalidTimeAt(t *testing.T) {
	is := is.New(t)

	req, _ := http.NewRequest("GET", createURL("/entities",
		fmt.Sprintf("timerel=%s", TemporalRelationAfterTime),
		"timeAt=notparseable"),
		nil)

	_, err := newQueryFromParameters(req, []string{"T"}, []string{"a"}, "")
	is.True(err != nil) // should return an error
}

func TestCreateTemporalQueryBeforeTime(t *testing.T) {
	is := is.New(t)

	timeAt, _ := time.Parse(time.RFC3339, DateTimeAt)
	req, _ := http.NewRequest("GET", createURL("/entities",
		fmt.Sprintf("timerel=%s", TemporalRelationBeforeTime),
		fmt.Sprintf("timeAt=%s", DateTimeAt)),
		nil)

	query, err := newQueryFromParameters(req, []string{"T"}, []string{"a"}, "")
	is.NoErr(err) // newQueryFromParameters failed
	is.True(query.IsTemporalQuery())

	from, to := query.Temporal().TimeSpan()
	is.True(from.IsZero()) // start of timespan should not be set
	is.Equal(to, timeAt)   // end of timespan should equal timeAt
}

func TestCreateTemporalQueryBetweenTwoTimes(t *testing.T) {
	is := is.New(t)

	timeAt, _ := time.Parse(time.RFC3339, DateTimeAt)
	endTimeAt, _ := time.Parse(time.RFC3339, DateEndTimeAt)

	req, _ := http.NewRequest("GET", createURL("/entities",
		fmt.Sprintf("timerel=%s", TemporalRelationBetweenTimes),
		fmt.Sprintf("timeAt=%s", DateTimeAt),
		fmt.Sprintf("endTimeAt=%s", DateEndTimeAt)),
		nil)

	query, err := newQueryFromParameters(req, []string{"T"}, []string{"a"}, "")
	is.NoErr(err) // newQueryFromParameters failed
	is.True(query.IsTemporalQuery())

	from, to := query.Temporal().TimeSpan()
	is.Equal(from, timeAt)  // check start of time span
	is.Equal(to, endTimeAt) // check end of time span
}

func TestCreateTemporalQueryBetweenFailsOnInvalidEndTimeAt(t *testing.T) {
	is := is.New(t)

	req, _ := http.NewRequest("GET", createURL("/entities",
		fmt.Sprintf("timerel=%s", TemporalRelationBetweenTimes),
		fmt.Sprintf("timeAt=%s", DateTimeAt),
		"endTimeAt=notparseable"),
		nil)

	_, err := newQueryFromParameters(req, []string{"T"}, []string{"a"}, "")
	is.True(err != nil) // should return an error
}

func TestCreateTemporalQueryReturnsErrorOnMissingEndTimeAt(t *testing.T) {
	is := is.New(t)

	req, _ := http.NewRequest("GET", createURL("/entities",
		fmt.Sprintf("timerel=%s", TemporalRelationBetweenTimes),
		fmt.Sprintf("timeAt=%s", DateTimeAt)),
		nil)

	_, err := newQueryFromParameters(req, []string{"T"}, []string{"a"}, "")
	is.True(err != nil) // should return an error
}

func TestCreateTemporalQueryWithCustomTemporalProperty(t *testing.T) {
	is := is.New(t)

	temporalProperty := "modifiedAt"
	req, _ := http.NewRequest("GET", createURL("/entities",
		fmt.Sprintf("timerel=%s", TemporalRelationAfterTime),
		fmt.Sprintf("timeAt=%s", DateTimeAt),
		fmt.Sprintf("timeproperty=%s", temporalProperty)),
		nil)

	query, err := newQueryFromParameters(req, []string{"T"}, []string{"a"}, "")
	is.NoErr(err) // newQueryFromParameters failed
	is.Equal(query.Temporal().Property(), temporalProperty)
}
