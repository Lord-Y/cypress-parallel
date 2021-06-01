// Package routers expose all routes of the api
package routers

var (
	// https://github.com/cypress-io/cypress-example-kitchensink/tree/master/cypress/integration/2-advanced-examples
	specs = []string{
		"cypress/integration/2-advanced-examples/actions.spec.js",
		"cypress/integration/2-advanced-examples/aliasing.spec.js",
		"cypress/integration/2-advanced-examples/assertions.spec.js",
		"cypress/integration/2-advanced-examples/connectors.spec.js",
		"cypress/integration/2-advanced-examples/cookies.spec.js",
		"cypress/integration/2-advanced-examples/files.spec.js",
		"cypress/integration/2-advanced-examples/location.spec.js",
		"cypress/integration/2-advanced-examples",
		"cypress/integration/1-getting-started",
	}

	// https://github.com/cypress-io/cypress-docker-images/tree/master/included
	// cypress docker images https://hub.docker.com/r/cypress/included
	cypressVersions = []string{
		"7.2.0-0.0.3",
	}
)
