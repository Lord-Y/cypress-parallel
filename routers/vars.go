// Package routers expose all routes of the api
package routers

var (
	// https://github.com/cypress-io/cypress-example-kitchensink/tree/master/cypress/e2e/2-advanced-examples
	specs = []string{
		"cypress/e2e/2-advanced-examples/actions.cy.js",
		"cypress/e2e/2-advanced-examples/aliasing.cy.js",
		"cypress/e2e/2-advanced-examples/assertions.cy.js",
		"cypress/e2e/2-advanced-examples/connectors.cy.js",
		"cypress/e2e/2-advanced-examples/cookies.cy.js",
		"cypress/e2e/2-advanced-examples/files.cy.js",
		"cypress/e2e/2-advanced-examples/location.cy.js",
		"cypress/e2e/2-advanced-examples",
		"cypress/e2e/1-getting-started",
	}

	// https://github.com/cypress-io/cypress-docker-images/tree/master/included
	// cypress docker images https://hub.docker.com/r/cypress/included
	cypressVersions = []string{
		"10.10.0-0.3.0",
	}
)
