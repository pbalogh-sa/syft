package cyclonedx13json

import (
	"flag"
	"regexp"
	"testing"

	"github.com/anchore/syft/internal/formats/common/testutils"
)

var updateCycloneDx = flag.Bool("update-cyclonedx", false, "update the *.golden files for cyclone-dx presenters")

func TestCycloneDxDirectoryPresenter(t *testing.T) {
	testutils.AssertPresenterAgainstGoldenSnapshot(t,
		Format().Presenter(testutils.DirectoryInput(t)),
		*updateCycloneDx,
		cycloneDxRedactor,
	)
}

func TestCycloneDxImagePresenter(t *testing.T) {
	testImage := "image-simple"
	testutils.AssertPresenterAgainstGoldenImageSnapshot(t,
		Format().Presenter(testutils.ImageInput(t, testImage)),
		testImage,
		*updateCycloneDx,
		cycloneDxRedactor,
	)
}

func cycloneDxRedactor(s []byte) []byte {
	serialPattern := regexp.MustCompile(`urn:uuid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	rfc3339Pattern := regexp.MustCompile(`([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\.[0-9]+)?(([Zz])|([\+|\-]([01][0-9]|2[0-3]):[0-5][0-9]))`)
	sha256Pattern := regexp.MustCompile(`sha256:[A-Fa-f0-9]{64}`)

	for _, pattern := range []*regexp.Regexp{serialPattern, rfc3339Pattern, sha256Pattern} {
		s = pattern.ReplaceAll(s, []byte("redacted"))
	}
	return s
}
