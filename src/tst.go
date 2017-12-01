package main

import (
	"github.com/subchen/go-xmldom"
	"fmt"
)

const (
	ExampleXml = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE junit SYSTEM "junit-result.dtd">
<testsuites>
	<testsuite tests="2" failures="0" time="0.009" name="github.com/subchen/go-xmldom">
		<properties>
			<property name="go.version">go1.8.1</property>
		</properties>
		<testcase classname="go-xmldom" id="ExampleParseXML" time="0.004"></testcase>
		<testcase classname="go-xmldom" id="ExampleParse" time="0.005"></testcase>
	</testsuite>
</testsuites>`
)


func ExampleNode_Query() {
	node := xmldom.Must(xmldom.ParseXML(ExampleXml)).Root
	// xpath expr: https://github.com/antchfx/xpath

	// find all children
	fmt.Printf("children = %v\n", len(node.Query("//*")))

	// find node matched tag name
	nodeList := node.Query("testsuites/testcase")
	for _, c := range nodeList {
		fmt.Printf("%v: id = %v\n", c.Name, c.GetAttributeValue("id"))
	}
	// Output:
	// children = 5
	// testcase: id = ExampleParseXML
	// testcase: id = ExampleParse
}

func main() {
	ExampleNode_Query()
}
