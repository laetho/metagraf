package metagraf

import "testing"

func TestSanitizeLabelValue(t *testing.T) {
	input := "Zaphod Bebelbrox (HHGTTG)"
	expected := "Zaphod_Bebelbrox"

	if sanitizeLabelValue(input) != expected {
		t.Errorf("Sanitizing label value with () failed, expected %v got %v", expected, sanitizeLabelValue(input))
	}

}