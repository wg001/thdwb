package ketchup

import (
	"regexp"
	"strings"

	"github.com/danfragoso/thdwb/mayo"
	"github.com/danfragoso/thdwb/structs"
)

var xmlTag = regexp.MustCompile(`(\<.+?\>)|(\<//?\w+\>\\?)`)
var clTag = regexp.MustCompile(`\<\/\w+\>`)
var tagContent = regexp.MustCompile(`(.+?)\<\/`)
var tagName = regexp.MustCompile(`(\<\w+)`)
var attr = regexp.MustCompile(`\w+=".+?"`)

func extractAttributes(tag string) []*structs.Attribute {
	rawAttrArray := attr.FindAllString(tag, -1)
	elementAttrs := []*structs.Attribute{}

	for i := 0; i < len(rawAttrArray); i++ {
		attrStringSlice := strings.Split(rawAttrArray[i], "=")
		attr := &structs.Attribute{
			Name:  attrStringSlice[0],
			Value: strings.Trim(attrStringSlice[1], "\""),
		}

		elementAttrs = append(elementAttrs, attr)
	}

	return elementAttrs
}

func ParseHTML(document string) *structs.NodeDOM {
	DOM_Tree := &structs.NodeDOM{
		Element:  "root",
		Content:  "THDWB",
		Children: []*structs.NodeDOM{},
		Style:    nil,
		Parent:   nil,
	}

	lastNode := DOM_Tree
	parseDocument := xmlTag.MatchString(document)
	document = strings.ReplaceAll(document, "\n", "")

	for parseDocument == true {
		var currentNode *structs.NodeDOM

		currentTag := xmlTag.FindString(document)
		currentTagIndex := xmlTag.FindStringIndex(document)

		if clTag.MatchString(currentTag) {
			contentStringMatch := tagContent.FindStringSubmatch(document)
			contentString := ""

			if len(contentStringMatch) > 1 {
				contentString = contentStringMatch[1]
			}

			if clTag.MatchString(contentString) {
				lastNode.Content = ""
			} else {
				lastNode.Content = strings.TrimSpace(contentString)
			}

			lastNode = lastNode.Parent
		} else {
			currentTagName := tagName.FindString(currentTag)
			extractedAttributes := extractAttributes(currentTag)
			parsedStylesheet := mayo.ParseInlineStylesheet(extractedAttributes)

			currentNode = &structs.NodeDOM{
				Element:    strings.Trim(currentTagName, "<"),
				Content:    "",
				Children:   []*structs.NodeDOM{},
				Attributes: extractedAttributes,
				Style:      parsedStylesheet,
				Parent:     lastNode,
			}

			lastNode.Children = append(lastNode.Children, currentNode)
			lastNode = currentNode
		}

		document = document[currentTagIndex[1]:len(document)]

		if !xmlTag.MatchString(document) {
			parseDocument = false
		}
	}

	return DOM_Tree
}