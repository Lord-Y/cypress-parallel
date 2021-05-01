// Package routers expose all routes of the api
package routers

var (
	// https://github.com/cypress-io/cypress-example-kitchensink/tree/master/cypress/integration/examples
	specs = []string{
		"cypress/integration/examples/actions.spec.js",
		"cypress/integration/examples/aliasing.spec.js",
		"cypress/integration/examples/assertions.spec.js",
		"cypress/integration/examples/connectors.spec.js",
		"cypress/integration/examples/cookies.spec.js",
		"cypress/integration/examples/files.spec.js",
		"cypress/integration/examples/location.spec.js",
	}

	// https://github.com/cypress-io/cypress-docker-images/tree/master/included
	cypressVersions = []string{
		"7.2.0",
		"7.1.0",
		"7.0.1",
	}
)
