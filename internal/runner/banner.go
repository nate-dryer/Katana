package runner

import (
	"fmt"

	"github.com/projectdiscovery/gologger"
	updateutils "github.com/projectdiscovery/utils/update"
)

var banner = fmt.Sprintf(`
   __        __                
  / /_____ _/ /____ ____  ___ _
 /  '_/ _  / __/ _  / _ \/ _  /
/_/\_\\_,_/\__/\_,_/_//_/\_,_/ %s							 
`, version)

var version = "v0.0.4-dev"

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s\n", banner)
	gologger.Print().Msgf("\t\tprojectdiscovery.io\n\n")
}

// GetUpdateCallback returns a callback function that updates katana
func GetUpdateCallback() func() {
	return func() {
		showBanner()
		updateutils.GetUpdateToolCallback("katana", version)()
	}
}
