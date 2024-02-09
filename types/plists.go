package types

import (
	"os"
	"path"
	"strings"

	"howett.net/plist"
)

// AppPlist represents the structure of a macOS app's Info.plist file
type AppPlist struct {
	CFBundleGetInfoString         string `plist:"CFBundleGetInfoString"`
	CFBundleExecutable            string `plist:"CFBundleExecutable"`
	CFBundleIconFile              string `plist:"CFBundleIconFile"`
	CFBundleIdentifier            string `plist:"CFBundleIdentifier"`
	CFBundleInfoDictionaryVersion string `plist:"CFBundleInfoDictionaryVersion"`
	CFBundleName                  string `plist:"CFBundleName"`
	CFBundlePackageType           string `plist:"CFBundlePackageType"`
	CFBundleShortVersionString    string `plist:"CFBundleShortVersionString"`
	// CFBundleVersion               string `plist:"CFBundleVersion"`
	// LSApplicationCategoryType     string `plist:"LSApplicationCategoryType"`
	// LSMinimumSystemVersion        string `plist:"LSMinimumSystemVersion"`
	NSHighResolutionCapable  bool   `plist:"NSHighResolutionCapable"`
	NSHumanReadableCopyright string `plist:"NSHumanReadableCopyright"`
	// NSQuitAlwaysKeepsWindows             bool   `plist:"NSQuitAlwaysKeepsWindows"`
	NSSupportsAutomaticGraphicsSwitching bool `plist:"NSSupportsAutomaticGraphicsSwitching"`
	IFMajorVersion                       int  `plist:"IFMajorVersion"`
	IFMinorVersion                       int  `plist:"IFMinorVersion"`
}

func NewAppPlist(a App) *AppPlist {
	return &AppPlist{
		CFBundleGetInfoString: a.Name,
		CFBundleExecutable:    a.Name,
		CFBundleIdentifier:    correctReverseDomain(a.BundleID),
		CFBundleName:          a.Name,
		CFBundleIconFile:      "icon.icns",
		// CFBundleVersion:               removeNonNumbers(a.Version),
		CFBundleShortVersionString:    "0.01",
		CFBundleInfoDictionaryVersion: "6.0",
		CFBundlePackageType:           "APPL",
		// LSApplicationCategoryType:     "public.app-category.developer-tools",
		// LSMinimumSystemVersion:        "10.14.0",
		NSHumanReadableCopyright: a.Copyright,
		// NSQuitAlwaysKeepsWindows:             true,
		IFMajorVersion:                       0,
		IFMinorVersion:                       1,
		NSHighResolutionCapable:              true,
		NSSupportsAutomaticGraphicsSwitching: true,
	}
}

func (p *AppPlist) CreatePlist(a App) error {
	plistPath := path.Join(a.AppPath, "Contents", "Info.plist")
	file, err := os.Create(plistPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := plist.NewEncoder(file)
	encoder.Indent("\t")
	err = encoder.Encode(p)
	if err != nil {
		return err
	}

	return nil
}

// private functions
func splitCamelCase(s string) string {
	var result string
	for i, r := range s {
		if i > 0 && (r >= 'A' && r <= 'Z') {
			result += " "
		}
		result += string(r)
	}
	return result
}

func removeNonNumbers(s string) string {
	var result string
	for _, r := range s {
		if (r >= '0' && r <= '9') || r == '.' {
			result += string(r)
		}
	}
	return result
}

func correctReverseDomain(s string) string {
	if strings.HasPrefix(s, "com.") {
		return s
	} else if strings.HasSuffix(s, ".com") {
		split := strings.Split(s, ".")
		for i, j := 0, len(split)-1; i < j; i, j = i+1, j-1 {
			split[i], split[j] = split[j], split[i]
		}
		return strings.Join(split, ".")
	}
	return s
}
