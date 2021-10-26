package responses

// BuiltInfo represents the build information for this application
//
// swagger:response BuiltInfo
type BuiltInfo struct {
	// The Version of running application
	Version string
	// The BuildDate running application
	BuildDate string
	// The BuildCommit of running application
	BuildCommit string
}
