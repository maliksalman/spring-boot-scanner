package main

import "github.com/maliksalman/spring-boot-scanner/cf"

func main() {
	apps := cf.ScanForApps()
	scannerOutput := cf.ScanAppContents(apps)
	scannerOutput.WriteAsJSON("output.json")
}
