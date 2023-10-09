package cf

import (
	"code.cloudfoundry.org/cli/actor/v7action"
	"code.cloudfoundry.org/cli/cf/commandregistry"
	"code.cloudfoundry.org/cli/cf/models"
	"code.cloudfoundry.org/cli/cf/trace"
	"code.cloudfoundry.org/cli/command/v7/shared"
	"code.cloudfoundry.org/cli/util/configv3"
	"code.cloudfoundry.org/cli/util/ui"
	"code.cloudfoundry.org/clock"
	"github.com/maliksalman/spring-boot-scanner/scan"
	"log"
	"os"
)

type CFApp struct {
	Name      string
	Org       string
	Space     string
	SpaceGUID string
	Buildpack string
	Instances int
	State     string
}

func ScanForApps() []CFApp {

	// do some setup
	traceLogger := trace.NewLogger(os.Stdout, false, "", "")
	deps := commandregistry.NewDependency(os.Stdout, traceLogger, os.Getenv("CF_DIAL_TIMEOUT"))
	defer deps.Config.Close()

	// get access to all the repos we need
	orgsRepo := deps.RepoLocator.GetOrganizationRepository()
	spacesRepo := deps.RepoLocator.GetSpaceRepository()
	appsSummaryRepo := deps.RepoLocator.GetAppSummaryRepository()
	appsRepo := deps.RepoLocator.GetApplicationRepository()

	// make the array to return
	cfApps := make([]CFApp, 0)

	// find all orgs
	orgs, _ := orgsRepo.ListOrgs(100)
	for _, org := range orgs {

		// set the current org as target
		deps.Config.SetOrganizationFields(org.OrganizationFields)

		// iterate over all spaces in current org
		spacesRepo.ListSpacesFromOrg(org.GUID, func(space models.Space) bool {

			// set the current space as target
			deps.Config.SetSpaceFields(space.SpaceFields)

			// find app-summaries in the current org/space
			apps, _ := appsSummaryRepo.GetSummariesInCurrentSpace()
			for _, app := range apps {

				// get more info about the app - like build-pack
				appObj, _ := appsRepo.GetApp(app.GUID)

				//
				cfApps = append(cfApps, CFApp{
					Name:      app.Name,
					Org:       org.Name,
					Space:     space.Name,
					SpaceGUID: space.GUID,
					Instances: app.RunningInstances,
					Buildpack: appObj.DetectedBuildpack,
					State:     appObj.State,
				})
			}
			return true
		})
	}

	return cfApps
}

func ScanAppContents(apps []CFApp) *ScannerOutput {

	// actor setup
	cfConfig, _ := configv3.GetCFConfig()
	commandUI, _ := ui.NewUI(cfConfig)
	ccClient, uaaClient, routingClient, _ := shared.GetNewClientsAndConnectToCF(cfConfig, commandUI, "")
	actor := v7action.NewActor(ccClient, cfConfig, nil, uaaClient, routingClient, clock.NewClock())

	// create the object we want to return
	scannerOutput := NewScannerOutput()

	// scan each app
	for _, app := range apps {

		if app.Buildpack == "java" {
			dropletBytes, _, _, _ := actor.DownloadCurrentDropletByAppName(app.Name, app.SpaceGUID)
			runtimeInfo := scan.FindRuntimeInfoFromContent(dropletBytes)
			_, javaCompilerVersion := scan.FindJavaCompilerVersionFromContent(dropletBytes, runtimeInfo.BootContentPrefix)

			// print app info
			scannerOutput.AddJavaApp(app.Org, app.Space, app.Name, app.Instances, app.State, runtimeInfo.JavaRuntimeVersion, javaCompilerVersion, runtimeInfo.BootVersion)
			log.Printf(
				"*** App: Org=%s, Space=%s, Name=%s, BuildPack=%s, SpringBoot=%s, JavaCompiler=%s, JavaRuntime=%s",
				app.Org,
				app.Space,
				app.Name,
				app.Buildpack,
				runtimeInfo.BootVersion,
				javaCompilerVersion,
				runtimeInfo.JavaRuntimeVersion,
			)
		} else {
			// print app info
			scannerOutput.AddOtherApp(app.Org, app.Space, app.Name, app.Instances, app.State, app.Buildpack)
			log.Printf(
				"*** App: Org=%s, Space=%s, Name=%s, BuildPack=%s",
				app.Org,
				app.Space,
				app.Name,
				app.Buildpack,
			)
		}
	}

	return scannerOutput
}
