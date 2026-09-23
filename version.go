package es

import "github.com/maloquacious/semver"

var (
	version = semver.Version{
		Major:      0,
		Minor:      2,
		Patch:      0,
		PreRelease: "alpha",

		// Automatically populate build metadata with commit info
		Build: semver.Commit(), // Uses Git commit hash from build info
	}
)

func Version() semver.Version {
	return version
}
