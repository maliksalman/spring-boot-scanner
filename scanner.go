package main

import (
	"code.cloudfoundry.org/cli/actor/v7action"
	"code.cloudfoundry.org/cli/cf/commandregistry"
	"code.cloudfoundry.org/cli/cf/models"
	"code.cloudfoundry.org/cli/cf/trace"
	"code.cloudfoundry.org/cli/command/v7/shared"
	"code.cloudfoundry.org/cli/util/configv3"
	"code.cloudfoundry.org/cli/util/ui"
	"code.cloudfoundry.org/clock"
	"log"
	"os"
)

func main() {

	// actor setup
	cfConfig, _ := configv3.GetCFConfig()
	commandUI, _ := ui.NewUI(cfConfig)
	ccClient, uaaClient, routingClient, _ := shared.GetNewClientsAndConnectToCF(cfConfig, commandUI, "")
	v7action.NewActor(ccClient, cfConfig, nil, uaaClient, routingClient, clock.NewClock())

	// do some setup
	traceLogger := trace.NewLogger(os.Stdout, false, "", "")
	deps := commandregistry.NewDependency(os.Stdout, traceLogger, os.Getenv("CF_DIAL_TIMEOUT"))
	defer deps.Config.Close()

	// get access to all the repos we need
	orgsRepo := deps.RepoLocator.GetOrganizationRepository()
	spacesRepo := deps.RepoLocator.GetSpaceRepository()
	appsSummaryRepo := deps.RepoLocator.GetAppSummaryRepository()
	appsRepo := deps.RepoLocator.GetApplicationRepository()

	scannerOutput := NewScannerOutput("output.json")

	// find all orgs
	orgs, _ := orgsRepo.ListOrgs(100)
	for _, org := range orgs {

		// set the current org as target
		deps.Config.SetOrganizationFields(org.OrganizationFields)

		// iterate over all spaces in current org
		spacesRepo.ListSpacesFromOrg(org.GUID, func(space models.Space) bool {
			log.Printf("Found: Org=%s, Space=%s", org.Name, space.Name)

			// set the current space as target
			deps.Config.SetSpaceFields(space.SpaceFields)

			// find app-summaries in the current org/space
			apps, _ := appsSummaryRepo.GetSummariesInCurrentSpace()
			for _, app := range apps {

				// get more info about the app - like build-pack
				appObj, _ := appsRepo.GetApp(app.GUID)

				if appObj.DetectedBuildpack == "java" {
					//dropletBytes, _, _, _ := actor.DownloadCurrentDropletByAppName(app.Name, space.GUID)
					dropletBytes, _ := os.ReadFile("test-droplet.tgz")

					javaRuntimeVersion, bootVersion := findJavaRuntimeAndBootVersions(dropletBytes)
					javaCompilerVersion := findJavaCompilerVersion(dropletBytes)

					// print app info
					scannerOutput.AddJavaApp(org.Name, space.Name, app.Name, app.RunningInstances, appObj.State, javaRuntimeVersion, javaCompilerVersion, bootVersion)
					log.Printf(
						"*** App: Name=%s, BuildPack=%s, SpringBoot=%s, JavaCompiler=%s, JavaRuntime=%s",
						app.Name,
						appObj.DetectedBuildpack,
						bootVersion,
						javaCompilerVersion,
						javaRuntimeVersion,
					)
				} else {
					// print app info
					scannerOutput.AddOtherApp(org.Name, space.Name, app.Name, app.RunningInstances, appObj.State, appObj.DetectedBuildpack)
					log.Printf(
						"*** App: Name=%s, BuildPack=%s",
						app.Name,
						appObj.DetectedBuildpack,
					)
				}
			}
			return true
		})
	}

	scannerOutput.writeAsJSON()
}
