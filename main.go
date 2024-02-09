// Thanks to https://gist.github.com/mattholt/2e4f9d729d3a5c8c2e8d3d7e7f7d6d6d
// https://medium.com/@mattholt/packaging-a-go-application-for-macos-f7084b00f6b5
//

package main

import (
	"flag"
	"fmt"
	"log"
	"path"

	"macbuilder/types"
)

var (
	name, desc, copyright, vers, bundleID, executable, iconPath, outputPath string
)

func init() {
	flag.StringVar(&name, "nm", "", "Name of the application (NO-SPACES or EXTENSION)")
	flag.StringVar(&executable, "ex", "", "Path to the executable you are wanting to bundle")
	flag.StringVar(&bundleID, "id", "", "reverse domain name style identifier")
	flag.StringVar(&iconPath, "icon", "", "path to the icon file must be at least 1024px square and PNG only")
	flag.StringVar(&vers, "v", "", "Version of the application")
	flag.StringVar(&desc, "d", "", "Description of the Application")
	flag.StringVar(&copyright, "c", "", "Copyright")
	flag.StringVar(&outputPath, "o", "", "Path to the output directory")
}

func main() {
	flag.Parse()
	if name == "" || executable == "" || iconPath == "" || outputPath == "" || bundleID == "" || vers == "" || desc == "" || copyright == "" {
		log.Println("All flags are required")
		flag.PrintDefaults()
		return
	}
	app := types.NewApp(name, desc, copyright, vers, bundleID, executable, iconPath)
	// Create the app's directory structure
	err := app.CreateDirectory(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	// Create the app's Info.plist file
	appPlist := types.NewAppPlist(*app)
	err = appPlist.CreatePlist(*app)
	if err != nil {
		log.Fatal(err)
	}
	// Create the app's icon file
	icon, err := types.NewIcon(iconPath, 1024, 1024)
	if err != nil {
		log.Fatal(err)
	}
	app.Icon = *icon
	tmpIcon := path.Join(app.AppPath, "Contents", "Resources", "icon.iconset")
	err = app.CreateIcons()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(tmpIcon)
	fmt.Println()
	fmt.Println("Application created at:", app.AppPath)
}
