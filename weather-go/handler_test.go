package function

import (
	"strings"
	"testing"
)

func TestNotHTMLResponse(t *testing.T) {

	result := Handle([]byte("london"))

	if strings.Contains(result,"<html>") {
		t.Errorf("Result is HTML expected plain text\n%s", result)
	}
}