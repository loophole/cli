package urlmaker

import "testing"

func TestReturnsCorrectHttpsUrl(t *testing.T) {
	expectedSiteID := "https://some-site.loophole.site"
	result := GetSiteURL("https", "some-site", "loophole.site")

	if result != expectedSiteID {
		t.Fatalf("Site URL '%s' is different than expected: %s", result, expectedSiteID)
	}
}

func TestReturnsCorrectHttpsUrlDifferentDomain(t *testing.T) {
	expectedSiteID := "https://some-site.loophole.host"
	result := GetSiteURL("https", "some-site", "loophole.host")

	if result != expectedSiteID {
		t.Fatalf("Site URL '%s' is different than expected: %s", result, expectedSiteID)
	}
}

func TestReturnsCorrectWebdavUrl(t *testing.T) {
	expectedSiteID := "webdav://some-site.loophole.site"
	result := GetSiteURL("webdav", "some-site", "loophole.site")

	if result != expectedSiteID {
		t.Fatalf("Site URL '%s' is different than expected: %s", result, expectedSiteID)
	}
}

func TestReturnsCorrectWebdavUrlDifferentDomain(t *testing.T) {
	expectedSiteID := "webdav://some-site.loophole.host"
	result := GetSiteURL("webdav", "some-site", "loophole.host")

	if result != expectedSiteID {
		t.Fatalf("Site URL '%s' is different than expected: %s", result, expectedSiteID)
	}
}

func TestReturnsCorrectDavsUrl(t *testing.T) {
	expectedSiteID := "davs://some-site.loophole.site"
	result := GetSiteURL("davs", "some-site", "loophole.site")

	if result != expectedSiteID {
		t.Fatalf("Site URL '%s' is different than expected: %s", result, expectedSiteID)
	}
}

func TestReturnsCorrectDavsUrlDifferentDomain(t *testing.T) {
	expectedSiteID := "davs://some-site.loophole.host"
	result := GetSiteURL("davs", "some-site", "loophole.host")

	if result != expectedSiteID {
		t.Fatalf("Site URL '%s' is different than expected: %s", result, expectedSiteID)
	}
}

func TestReturnsCorrectFQDN(t *testing.T) {
	expectedSiteID := "some-site.loophole.site"
	result := GetSiteFQDN("some-site", "loophole.site")

	if result != expectedSiteID {
		t.Fatalf("Site FQDN '%s' is different than expected: %s", result, expectedSiteID)
	}
}

func TestReturnsCorrectFQDNDifferentDomain(t *testing.T) {
	expectedSiteID := "some-site.loophole.host"
	result := GetSiteFQDN("some-site", "loophole.host")

	if result != expectedSiteID {
		t.Fatalf("Site FQDN '%s' is different than expected: %s", result, expectedSiteID)
	}
}
