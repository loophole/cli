package urlmaker

import "testing"

func TestReturnsCorrectHttpsUrl(t *testing.T) {
	expectedSiteID := "https://some-site.loophole.host"
	result := GetSiteUrl("https", "some-site")

	if result != expectedSiteID {
		t.Fatalf("Site ID '%s' is different than expected: %s", result, expectedSiteID)
	}
}

func TestReturnsCorrectWebdavUrl(t *testing.T) {
	expectedSiteID := "webdav://some-site.loophole.host"
	result := GetSiteUrl("webdav", "some-site")

	if result != expectedSiteID {
		t.Fatalf("Site ID '%s' is different than expected: %s", result, expectedSiteID)
	}
}

func TestReturnsCorrectDavsUrl(t *testing.T) {
	expectedSiteID := "davs://some-site.loophole.host"
	result := GetSiteUrl("davs", "some-site")

	if result != expectedSiteID {
		t.Fatalf("Site ID '%s' is different than expected: %s", result, expectedSiteID)
	}
}
