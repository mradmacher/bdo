package bdo

import (
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
func assertEqualStr(t *testing.T, got, want, name string) {
	if got != want {
		t.Errorf("%s = %q, want %q", name, got, want)
	}
}
func assertEqual[K comparable](t *testing.T, got, want K, name string) {
	if got != want {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

func assertHasCapability(t *testing.T, inst Installation, want Capability) {
	found := false
	for _, c := range inst.Capabilities {
		if c.WasteCode == want.WasteCode &&
			c.Dangerous == want.Dangerous &&
			c.ProcessCode == want.ProcessCode &&
			c.ActivityCode == want.ActivityCode &&
			c.Quantity == want.Quantity {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected %q to has capability %v", inst.Name, want)
	}
}

func isSameMaterials(col1, col2 []Material) bool {
	if len(col1) != len(col2) {
		return false
	}
	for _, c1 := range col1 {
		found := false
		for _, c2 := range col2 {
			if c1.Code == c2.Code {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func assertIncludeCapability(t *testing.T, collection []Capability, want Capability) {
	found := false
	for _, c := range collection {
		if c.WasteCode == want.WasteCode &&
			c.Dangerous == want.Dangerous &&
			c.ProcessCode == want.ProcessCode &&
			c.ActivityCode == want.ActivityCode &&
			c.Quantity == want.Quantity &&
			isSameMaterials(c.Materials, want.Materials) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected %v to has capability %v", collection, want)
	}
}

type repositoryTest func(*testing.T, *Repository)

func setupSuite(t *testing.T) (func(*testing.T), *Repository) {
	if err := godotenv.Load("../.test.env"); err != nil {
		t.Fatalf("No .env file found")
	}
	dbUri := os.Getenv("BDO_DB_URI")

	r := Repository{DBUri: dbUri}
	err := r.Connect()
	if err != nil {
		t.Fatalf("Problem with DB connection: %v", err)
	}
	return func(t *testing.T) {
		r.Disconnect()
	}, &r
}

func setupTest(t *testing.T, r *Repository) func(*testing.T) {
	return func(t *testing.T) {
		r.Purge()
	}
}

func TestRepo(t *testing.T) {
	teardownSuite, r := setupSuite(t)
	defer teardownSuite(t)

	tests := map[string]repositoryTest{
		"Summarize":           testSummarize,
		"FindInstallation":    testFindInstallation,
		"SearchCapabilities":  testSearchCapabilities,
		"SearchInstallations": testSearchInstallations,
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			teardownTest := setupTest(t, r)
			defer teardownTest(t)

			test(t, r)
		})
	}
}
