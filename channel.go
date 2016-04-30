package telly

// Channel data structure. Populated with HDHomeRun lineup
//
type Channel struct {
	GuideNumber string
	GuideName   string
	HD          ConvertibleBoolean
	URL         string
}
