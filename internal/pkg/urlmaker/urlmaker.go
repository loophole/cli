package urlmaker

import "fmt"

const (
	// HostURL is the top-level domain used for hosting loophole sites
	HostURL = "loophole.site"
)

// GetSiteURL produces URL for the site (with the protocol)
func GetSiteURL(protocol string, siteID string) string {
	return fmt.Sprintf("%s://%s.%s", protocol, siteID, HostURL)
}

// GetSiteFQDN produces fully qualified domain name for the site (without the protocol)
func GetSiteFQDN(siteID string) string {
	return fmt.Sprintf("%s.%s", siteID, HostURL)

}
