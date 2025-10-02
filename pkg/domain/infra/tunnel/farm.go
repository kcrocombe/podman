package tunnel

import (
	"context"
	"fmt"

	"github.com/containers/podman/v5/pkg/bindings/system"
	"github.com/containers/podman/v5/pkg/domain/entities"
)

const (
	remoteFarmImageBuilderDriver = "podman-remote"
)

// FarmNodeName returns the remote engine's name.
func (ir *ImageEngine) FarmNodeName(ctx context.Context) string {
	return ir.NodeName
}

// FarmNodeDriver returns a description of the image builder driver
func (ir *ImageEngine) FarmNodeDriver(ctx context.Context) string {
	return remoteFarmImageBuilderDriver
}

// Issue #26822 : Function now returns an additional array of emulatedPlatforms (describes the platforms the engine can emulate)
func (ir *ImageEngine) fetchInfo(_ context.Context) (os, arch, variant string, nativePlatforms []string, emulatedPlatforms []string, err error) {
	engineInfo, err := system.Info(ir.ClientCtx, &system.InfoOptions{})
	if err != nil {
		return "", "", "", nil, nil, fmt.Errorf("retrieving host info from %q: %w", ir.NodeName, err)
	}
	nativePlatform := engineInfo.Host.OS + "/" + engineInfo.Host.Arch
	if engineInfo.Host.Variant != "" {
		nativePlatform = nativePlatform + "/" + engineInfo.Host.Variant
	}
	return engineInfo.Host.OS, engineInfo.Host.Arch, engineInfo.Host.Variant, []string{nativePlatform}, engineInfo.Host.EmulatedArchitectures, nil
}

// FarmNodeInspect returns information about the remote engines in the farm
// Issue #26822 : FarmInspectReport struct now also includes info regarding the platforms the farm can emulate
func (ir *ImageEngine) FarmNodeInspect(ctx context.Context) (*entities.FarmInspectReport, error) {
	ir.platforms.Do(func() {
		ir.os, ir.arch, ir.variant, ir.nativePlatforms, ir.emulatedPlatforms, ir.platformsErr = ir.fetchInfo(ctx)
	})
	return &entities.FarmInspectReport{NativePlatforms: ir.nativePlatforms, EmulatedPlatforms: ir.emulatedPlatforms,
		OS:      ir.os,
		Arch:    ir.arch,
		Variant: ir.variant}, ir.platformsErr
}
