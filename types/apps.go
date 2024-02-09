package types

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path"
)

type App struct {
	Name           string
	Description    string
	Copyright      string
	BundleID       string
	ExecutablePath string
	Icon           Icon
	IconPath       string
	AppPath        string
	Version        string
	FullAppName    string
}

func NewApp(name, desc, copyright, vers, bundleID, executable, iconPath string) *App {
	return &App{
		Name:           name,
		Description:    desc,
		Copyright:      copyright,
		Version:        vers,
		BundleID:       bundleID,
		ExecutablePath: executable,
		IconPath:       iconPath,
		FullAppName:    name + ".app",
	}
}

func (a *App) CreateDirectory(outputFolderPath string) error {
	// Create the app's directory structure
	a.AppPath = path.Join(outputFolderPath, a.Name+".app")
	err := os.MkdirAll(a.AppPath, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Join(a.AppPath, "Contents", "MacOS"), os.ModePerm)
	if err != nil {
		return err
	} else {
		// Copy the executable to the app's MacOS directory
		srcFile, err := os.Open(a.ExecutablePath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(path.Join(a.AppPath, "Contents", "MacOS", a.Name))
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}
	}
	err = os.MkdirAll(path.Join(a.AppPath, "Contents", "Resources"), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) CreateIcons() error {
	if a.IconPath != "" {
		a.Icon = Icon{
			FilePath: a.IconPath,
		}
	} else {
		return errors.New("no icon provided")
	}
	if a.Icon.Exists() {
		err := a.Icon.CreateMacIcons(path.Join(a.AppPath, "Contents", "Resources", "icon.iconset"))
		if err != nil {
			return err
		} else {
			err = iconUtil(path.Join(a.AppPath, "Contents", "Resources"))
			if err != nil {
				return err
			}
			// else {
			// 	return os.RemoveAll(path.Join(a.AppPath, "Contents", "Resources", "icon.iconset"))
			// }
		}
	}

	return nil
}

// private function
func iconUtil(resourcePath string) error {
	_, err := os.Stat(resourcePath)
	if err != nil {
		return err
	}
	// Change working directory
	err = os.Chdir(resourcePath)
	if err != nil {
		return err
	}
	// Continue with the rest of your code
	cmd := exec.Command("iconutil", "-c", "icns", "icon.iconset")

	return cmd.Run()
}
