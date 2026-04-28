package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type goTestEvent struct {
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Output  string  `json:"Output"`
	Elapsed float64 `json:"Elapsed"`
}

type testCase struct {
	Name     string
	Duration float64
	Failure  string
	SystemOut strings.Builder
}

type testSuite struct {
	Name      string
	TestCases map[string]*testCase
	Failures  []string
}

type xmlFailure struct {
	Message string `xml:"message,attr"`
	Body    string `xml:",chardata"`
}

type xmlCase struct {
	XMLName   xml.Name    `xml:"testcase"`
	Name      string      `xml:"name,attr"`
	ClassName string      `xml:"classname,attr"`
	Time      string      `xml:"time,attr"`
	Failure   *xmlFailure `xml:"failure,omitempty"`
	SystemOut string      `xml:"system-out,omitempty"`
}

type xmlSuite struct {
	XMLName  xml.Name  `xml:"testsuite"`
	Name     string    `xml:"name,attr"`
	Tests    int       `xml:"tests,attr"`
	Failures int       `xml:"failures,attr"`
	Cases    []xmlCase `xml:"testcase"`
}

type xmlSuites struct {
	XMLName xml.Name   `xml:"testsuites"`
	Suites  []xmlSuite `xml:"testsuite"`
}

func main() {
	suites, err := readEvents(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read events: %v\n", err)
		os.Exit(1)
	}

	out, err := toXML(suites)
	if err != nil {
		fmt.Fprintf(os.Stderr, "encode junit: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stdout.Write(out); err != nil {
		fmt.Fprintf(os.Stderr, "write junit: %v\n", err)
		os.Exit(1)
	}
}

func readEvents(r io.Reader) (map[string]*testSuite, error) {
	scanner := bufio.NewScanner(r)
	suites := map[string]*testSuite{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var event goTestEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return nil, err
		}

		if event.Package == "" {
			continue
		}

		suite := suites[event.Package]
		if suite == nil {
			suite = &testSuite{
				Name:      event.Package,
				TestCases: map[string]*testCase{},
			}
			suites[event.Package] = suite
		}

		if event.Test == "" {
			if event.Action == "fail" && strings.TrimSpace(event.Output) != "" {
				suite.Failures = append(suite.Failures, event.Output)
			}
			continue
		}

		tc := suite.TestCases[event.Test]
		if tc == nil {
			tc = &testCase{Name: event.Test}
			suite.TestCases[event.Test] = tc
		}

		if event.Output != "" {
			tc.SystemOut.WriteString(event.Output)
		}

		switch event.Action {
		case "fail":
			tc.Duration = event.Elapsed
			tc.Failure = strings.TrimSpace(tc.SystemOut.String())
		case "pass":
			tc.Duration = event.Elapsed
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return suites, nil
}

func toXML(suites map[string]*testSuite) ([]byte, error) {
	names := make([]string, 0, len(suites))
	for name := range suites {
		names = append(names, name)
	}
	sort.Strings(names)

	doc := xmlSuites{
		Suites: make([]xmlSuite, 0, len(names)),
	}

	for _, name := range names {
		suite := suites[name]
		caseNames := make([]string, 0, len(suite.TestCases))
		for testName := range suite.TestCases {
			caseNames = append(caseNames, testName)
		}
		sort.Strings(caseNames)

		xs := xmlSuite{
			Name:  suite.Name,
			Tests: len(caseNames),
		}

		for _, testName := range caseNames {
			tc := suite.TestCases[testName]
			xc := xmlCase{
				Name:      tc.Name,
				ClassName: suite.Name,
				Time:      fmt.Sprintf("%.3f", tc.Duration),
				SystemOut: strings.TrimSpace(tc.SystemOut.String()),
			}
			if tc.Failure != "" {
				xs.Failures++
				xc.Failure = &xmlFailure{
					Message: "test failed",
					Body:    tc.Failure,
				}
			}
			xs.Cases = append(xs.Cases, xc)
		}

		doc.Suites = append(doc.Suites, xs)
	}

	raw, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), raw...), nil
}
