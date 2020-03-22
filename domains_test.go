package abuse

import (
	"os"
	"testing"
)

func TestCheckAbusiveDomain(t *testing.T) {
	if !IsAbusiveDomain("mailinator.com") {
		t.Fail()
	}
}

func TestCheckValidDomain(t *testing.T) {
	if IsAbusiveDomain("gmail.com") {
		t.Fail()
	}
}
func TestCheckTempEmail(t *testing.T) {
	if abusive, err := IsTempEmail("user@mailinator.com"); err != nil {
		t.Log(err.Error())
		t.Fail()
	} else if !abusive {
		t.Fail()
	}
}

func TestCheckValidEmail(t *testing.T) {
	if abusive, err := IsTempEmail("user@google.com"); err != nil {
		t.Log(err.Error())
		t.Fail()
	} else if abusive {
		t.Fail()
	}
}

func TestNoInitAddition(t *testing.T) {
	os.Remove(additionsPath)

	domain := "an-abusive-domain.com"
	if err := Add(domain); err != nil {
		// Adding domain without prior call to Init("")
		t.Log(err.Error())
		t.Fail()
	}

	if _, err := os.Stat(additionsPath); err != nil {
		// File not created
		t.Log(err.Error())
		t.Fail()
	}

	if !IsAbusiveDomain(domain) {
		t.Log("newly added domain not in tempDomain map")
		t.Fail()
	}

	if err := Close(); err != nil {
		t.Log(err.Error())
		t.Fail()
	}
}
