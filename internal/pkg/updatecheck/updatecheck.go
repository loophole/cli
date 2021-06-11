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
	"github.com/pkg/browser"
	"github.com/spf13/viper"
)

func CheckForUpdates() {
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
			remind, usePrompt, err := remindAgainCheck(availableVersionParsed)
			if err != nil {
				communication.Error(err.Error()) //errors in determining the type reminder should be noted, but not interrupt the program
			}
			if !remind {
				return
			}
			downloadlink := getDownloadLink(availableVersion.Version)
			if usePrompt { //either use a notification that the user needs to click away, or use a notification they can ignore
				openLink := false //needs to be declared here instead of below with := so we can still have access to err outside of this scope
				openLink, err = zenity.Question(fmt.Sprintf("A new version is available for you at \n%s \n Do you want to open this link in your browser?", downloadlink), zenity.NoWrap(), zenity.Title("New version available!"))
				if openLink {
					browser.OpenURL(downloadlink)
				}
			} else {
				err = zenity.Notify(fmt.Sprintf("A new version is available for you at \n%s \n", downloadlink), zenity.Title("New version available!"))
			}
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
		operatingSystem = "macos" //rename for use in download url
	}
	if architecture == "amd64" {
		architecture = "64bit"
	} else if architecture == "386" {
		architecture = "32bit"
	} else {
		communication.Error("There was an error detecting your system architecture.") //if arch is unexpected, only link to the release page
		return fmt.Sprintf("https://github.com/loophole/cli/releases/tag/%s", availableVersion)
	}
	link := fmt.Sprintf("https://github.com/loophole/cli/releases/download/%s/loophole-desktop_%s_%s_%s%s", availableVersion, availableVersion, operatingSystem, architecture, archiveType)
	fmt.Println(link)
	return link
}

func remindAgainCheck(availableVersionParsed semver.Version) (bool, bool, error) {
	home, err := homedir.Dir()
	if err != nil {
		return true, false, err
	}

	viper.SetDefault("lastreminder", time.Time{})         //date of last reminder, default zero value for time
	viper.SetDefault("availableversion", "1.0.0-beta.14") //last seen latest version
	viper.SetDefault("remindercount", 3)                  //counts to zero, then switches from prompt to notification reminder
	viper.SetConfigName("config")                         //name of config file (without extension)
	viper.SetConfigType("json")
	viper.AddConfigPath(fmt.Sprintf("%s/.loophole/", home))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok { //create a config if none exist yet
			viper.WriteConfigAs(fmt.Sprintf("%s/.loophole/config.json", home))
		} else {
			return true, false, err
		}
	}

	lastSeenLatestVersion, err := semver.Make(viper.GetString("availableversion"))
	if availableVersionParsed.GT(lastSeenLatestVersion) { //reset reminder count if new version is out
		viper.Set("availableversion", availableVersionParsed.String())
		viper.Set("remindercount", 3)
	}
	if err != nil {
		return true, false, err
	}
	lastReminder := viper.GetTime("lastreminder")
	if (lastReminder.Year() < time.Now().Year()) || (lastReminder.YearDay() < time.Now().YearDay()) { //check if reminder has been done today
		viper.Set("lastreminder", time.Now())
		viper.WriteConfigAs("/home/work/.loophole/config.json")
		if viper.GetInt("remindercount") < 1 {
			return true, false, nil
		} else {
			viper.Set("remindercount", viper.GetInt("remindercount")-1)
			viper.WriteConfigAs("/home/work/.loophole/config.json")
			return true, true, nil
		}
	}

	return false, false, nil
}
