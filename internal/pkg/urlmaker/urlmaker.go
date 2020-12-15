package urlmaker

import "fmt"

func GetSiteUrl(protocol string, siteID string) string {
	return fmt.Sprintf("%s://%s.loophole.host", protocol, siteID)
}
