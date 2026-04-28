package main

import (
	"strings"
	"testing"
)

func TestReadEventsAndToXML_ConvertsGoTestOutputToJUnit(t *testing.T) {
	input := strings.Join([]string{
		`{"Action":"run","Package":"example/service","Test":"TestHappyPath"}`,
		`{"Action":"output","Package":"example/service","Test":"TestHappyPath","Output":"=== RUN   TestHappyPath\n"}`,
		`{"Action":"pass","Package":"example/service","Test":"TestHappyPath","Elapsed":0.012}`,
		`{"Action":"run","Package":"example/service","Test":"TestFailurePath"}`,
		`{"Action":"output","Package":"example/service","Test":"TestFailurePath","Output":"expected failure details\n"}`,
		`{"Action":"fail","Package":"example/service","Test":"TestFailurePath","Elapsed":0.015}`,
	}, "\n")

	suites, err := readEvents(strings.NewReader(input))
	if err != nil {
		t.Fatalf("readEvents() error = %v", err)
	}

	out, err := toXML(suites)
	if err != nil {
		t.Fatalf("toXML() error = %v", err)
	}

	xml := string(out)
	if !strings.Contains(xml, `testsuite name="example/service"`) {
		t.Fatalf("xml = %s, want testsuite", xml)
	}
	if !strings.Contains(xml, `testcase name="TestFailurePath"`) {
		t.Fatalf("xml = %s, want failing testcase", xml)
	}
	if !strings.Contains(xml, `<failure`) {
		t.Fatalf("xml = %s, want failure element", xml)
	}
}
