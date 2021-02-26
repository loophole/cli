package urlmaker

import "fmt"

// GetSiteURL produces URL for the site (with the protocol)
func GetSiteURL(protocol string, siteID string, domain string) string {
	return fmt.Sprintf("%s://%s.%s", protocol, siteID, domain)
}

// GetSiteFQDN produces fully qualified domain name for the site (without the protocol)
func GetSiteFQDN(siteID string, domain string) string {
	return fmt.Sprintf("%s.%s", siteID, domain)
}
