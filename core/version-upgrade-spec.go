package core

import (
	"fmt"
	"io"

	"github.com/seboste/sapper/ports"
)

type VersionUpgradeSpec struct {
	previous, target, latestAvailable, latestWorking string
}

func (vus VersionUpgradeSpec) UpgradeRequired() bool {
	return vus.previous != vus.latestAvailable
}
func (vus VersionUpgradeSpec) UpgradeToLatestSuccessful() bool {
	return vus.latestWorking == vus.latestAvailable
}
func (vus VersionUpgradeSpec) UpgradeToTargetSuccessful() bool {
	return vus.latestWorking == vus.target
}
func (vus VersionUpgradeSpec) UpgradePartiallyFailed() bool {
	return vus.UpgradeRequired() && !vus.UpgradeToLatestSuccessful() && !vus.UpgradeToTargetSuccessful() && vus.latestWorking != vus.previous
}
func (vus VersionUpgradeSpec) UpgradeCompletelyFailed() bool {
	return vus.UpgradeRequired() && vus.latestWorking == vus.previous
}

func (vus VersionUpgradeSpec) PrintStatus(w io.Writer, d ports.PackageDependency) {
	if !vus.UpgradeRequired() {
		fmt.Fprintf(w, "%s is already now up to date. No upgrade required.\n", d.Id)
	} else if vus.UpgradeToLatestSuccessful() {
		fmt.Fprintf(w, "upgrade from %s to %s succeeded. %s is now up to date.\n", vus.previous, vus.latestAvailable, d.Id)
	} else if vus.UpgradeToTargetSuccessful() {
		fmt.Fprintf(w, "upgrade from %s to %s succeeded. However, there is a newer version %s available.\n", vus.previous, vus.target, vus.latestAvailable)
	} else if vus.UpgradePartiallyFailed() {
		if vus.target == vus.latestAvailable {
			fmt.Fprintf(w, "upgrade from %s to %s failed => upgrade to latest working version %s instead\n", vus.previous, vus.target, vus.latestWorking)
		} else {
			fmt.Fprintf(w, "upgrade from %s to %s failed => upgrade to latest working version %s instead. Note that there is an even newer version %s available.\n", vus.previous, vus.target, vus.latestWorking, vus.latestAvailable)
		}
	} else if vus.UpgradeCompletelyFailed() {
		fmt.Fprintf(w, "upgrade from %s to %s failed => keeping version %s\n", vus.previous, vus.target, vus.previous)
	}
}
