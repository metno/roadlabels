package exttools

import (
	"math"
	"os"
	"strings"
	"testing"
	"time"
)

func init() {

}

var Epsilon = math.Nextafter(1.0, 2.0) - 1.0

func TestTempNotfound(t *testing.T) {
	os.Setenv("PYTHONPATH", os.Getenv("PYTHONPATH")+":"+os.Getenv("PWD")+"/../../exttools")
	theTime := time.Date(1998, 8, 15, 16, 0, 0, 0, time.UTC)
	_, err := GetTemp(theTime, 59.959496162, 10.783663532)
	want := "No suitable analysis file found"
	if !strings.Contains(err.Error(), want) {
		t.Errorf(`Want: %v Got:  %v`, want, err.Error())
	}
}

func TestGetTemp(t *testing.T) {
	os.Setenv("PYTHONPATH", os.Getenv("PYTHONPATH")+":"+os.Getenv("PWD")+"/../../exttools")
	AnalysisDir = os.Getenv("PWD") + "/../../test-fixtures"
	theTime := time.Date(2023, 3, 2, 9, 0, 0, 0, time.UTC)
	temp, err := GetTemp(theTime, 59.959496162, 10.783663532)

	if err != nil {
		t.Fatalf(`Want: err==nil Got:  %v`, err.Error())
	}

	if -1.232825-temp > Epsilon {
		t.Errorf("Want: temp==-1.23 Got: %.8f", temp)
	}

}
