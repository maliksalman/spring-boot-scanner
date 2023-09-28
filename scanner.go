package main

import (
	"code.cloudfoundry.org/cli/cf/commandregistry"
	"code.cloudfoundry.org/cli/cf/models"
	"code.cloudfoundry.org/cli/cf/trace"
	"log"
	"os"
)

func main() {
	printOrgsAndSpaces()
}

func printOrgsAndSpaces() {

	// do some setup
	traceLogger := trace.NewLogger(os.Stdout, false, "", "")
	deps := commandregistry.NewDependency(os.Stdout, traceLogger, os.Getenv("CF_DIAL_TIMEOUT"))
	defer deps.Config.Close()

	// get access to all the repos we need
	orgsRepo := deps.RepoLocator.GetOrganizationRepository()
	spacesRepo := deps.RepoLocator.GetSpaceRepository()
	appsRepo := deps.RepoLocator.GetAppSummaryRepository()

	// find all orgs
	orgs, _ := orgsRepo.ListOrgs(100)
	for _, org := range orgs {

		// set the current org as target
		deps.Config.SetOrganizationFields(org.OrganizationFields)

		// iterate over all spaces in current org
		spacesRepo.ListSpacesFromOrg(org.GUID, func(space models.Space) bool {
			log.Printf("Org: %s, Space: %s", org.Name, space.Name)

			// set the current space as target
			deps.Config.SetSpaceFields(space.SpaceFields)

			// find app-summaries in the current org/space
			apps, _ := appsRepo.GetSummariesInCurrentSpace()
			for _, app := range apps {
				// print app info
				log.Printf("*** App: Name=%s, BuildPack=%s, Instances=%d, State=%s", app.Name, app.DetectedBuildpack, app.RunningInstances, app.State)
			}
			return true
		})
	}
}
