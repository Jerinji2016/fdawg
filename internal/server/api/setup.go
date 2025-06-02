package api

import (
	"github.com/Jerinji2016/fdawg/pkg/flutter"
)

// SetupAPIRoutes sets up all API routes for the server
func SetupAPIRoutes(project *flutter.ValidationResult) {
	// Setup individual API route groups
	SetupAssetAPIRoutes(project)
	SetupEnvironmentAPIRoutes(project)
	SetupLocalizationAPIRoutes(project)
	SetupTranslationAPIRoutes(project)
	SetupNamerAPIRoutes(project)
	SetupBundlerAPIRoutes(project)
}
