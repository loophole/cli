package updatecheck

import (
	"fmt"
	"runtime"
	"time"

	"github.com/blang/semver/v4"
	"github.com/loophole/cli/config"
	"github.com/loophole/cli/internal/pkg/apiclient"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/mitchellh/go-homedir"
	"github.com/ncruces/zenity"
	"github.com/spf13/viper"
)

func CheckVersion() {
	availableVersion, err := apiclient.GetLatestAvailableVersion()
	if err != nil {
		communication.Debug("There was a problem obtaining info response, skipping further checking")
		return
	}
	currentVersionParsed, err := semver.Make(config.Config.Version)
	if err != nil {
		communication.Debug(fmt.Sprintf("Cannot parse current version '%s' as semver version, skipping further checking", config.Config.Version))
		return
	}
	availableVersionParsed, err := semver.Make(availableVersion.Version)
	if err != nil {
		communication.Debug(fmt.Sprintf("Cannot parse available version '%s' as semver version, skipping further checking", availableVersion))
		return
	}
	if currentVersionParsed.LT(availableVersionParsed) {
		if config.Config.ClientMode == "cli" {
			communication.NewVersionAvailable(availableVersion.Version)
		} else {
			remind, err := remindAgainCheck()
			if err != nil {
				communication.Error(err.Error()) //errors in retrieving a download link should be noted, but not interrupt the program
			}
			if !remind {
				return
			}
			downloadlink := getDownloadLink(availableVersion.Version)
			err = zenity.Notify(fmt.Sprintf("A new version is available for you at \n%s \n", downloadlink), zenity.Title("New version available!"))
			if err != nil {
				communication.Debug(err.Error()) //errors in showing a download link should be noted, but not interrupt the program
			}
		}
	}
}

func getDownloadLink(availableVersion string) string {
	archiveType := ".tar.gz"
	operatingSystem := runtime.GOOS
	architecture := runtime.GOARCH
	if operatingSystem == "windows" {
		archiveType = ".zip"
	} else if operatingSystem == "darwin" {
		operatingSystem = "macos" //rename for use in downloadlink
	}
	if architecture == "amd64" {
		architecture = "64bit"
	} else if architecture == "386" {
		architecture = "32bit"
	} else {
		communication.Error("There was an error detecting your system architecture.") //if arch is unexpected, only link to the release page
		return fmt.Sprintf("https://github.com/loophole/cli/releases/tag/%s", availableVersion)
	}
	link := fmt.Sprintf("https://github.com/loophole/cli/releases/download/%s/loophole-desktop_%s_%s_%s%s", availableVersion, availableVersion, operatingSystem, operatingSystem, archiveType)
	fmt.Println(link)
	return link
}

func remindAgainCheck() (bool, error) {
	home, err := homedir.Dir()
	if err != nil {
		return true, err
	}

	layout := "2006-02-01"                                        //golangs arcane time format string
	viper.SetDefault("last-reminder", time.Time{}.Format(layout)) //zero value for time
	viper.SetConfigName("config")                                 // name of config file (without extension)
	viper.SetConfigType("json")
	viper.AddConfigPath(fmt.Sprintf("%s/.loophole/", home))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { //create a config if none exist yet
			viper.WriteConfigAs(fmt.Sprintf("%s/.loophole/config.json", home))
		} else {
			return true, err
		}
	}

	t, err := time.Parse(layout, fmt.Sprintf("%v", viper.Get("last-reminder")))
	if err != nil {
		return true, err
	}
	if (t.Year() < time.Now().Year()) || (t.YearDay() < time.Now().YearDay()) { //check if reminder has been done today
		viper.Set("last-reminder", time.Now().Format(layout))
		viper.WriteConfigAs("/home/work/.loophole/config.json")
		return true, nil
	}

	return false, nil
}
