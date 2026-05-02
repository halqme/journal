package journal

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Section struct {
	Heading string
	Raw     string
}

func headingText(n *ast.Heading, source []byte) string {
	var buf bytes.Buffer
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if t, ok := c.(*ast.Text); ok {
			buf.Write(t.Segment.Value(source))
		}
	}
	return buf.String()
}

func ParseMarkdown(src string) []Section {
	if strings.TrimSpace(src) == "" {
		return nil
	}
	md := goldmark.New()
	source := []byte(src)
	doc := md.Parser().Parse(text.NewReader(source))

	var headingPositions []int
	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		if h, ok := n.(*ast.Heading); ok && h.Level == 2 {
			headingPositions = append(headingPositions, n.Pos())
		}
	}
	if len(headingPositions) == 0 {
		return nil
	}

	var sections []Section
	for i := 0; i < len(headingPositions); i++ {
		start := headingPositions[i]
		end := len(source)
		if i+1 < len(headingPositions) {
			end = headingPositions[i+1]
		}
		line := strings.Split(string(source[start:end]), "\n")[0]
		sections = append(sections, Section{
			Heading: line,
			Raw:     string(source[start:end]),
		})
	}
	return sections
}

func RenderHeader(template, timeValue string) string {
	h := strings.ReplaceAll(template, "{time}", timeValue)
	return strings.TrimRight(h, "\n")
}

func AppendToSections(sections []Section, timeValue, template string, sameTime bool, entryText string) string {
	rendered := RenderHeader(template, timeValue)
	if sameTime && len(sections) > 0 {
		lastIdx := len(sections) - 1
		sections[lastIdx].Raw = strings.TrimRight(sections[lastIdx].Raw, "\n") + "\n" + entryText
	} else {
		sections = append(sections, Section{
			Heading: rendered,
			Raw:     rendered + "\n" + entryText,
		})
	}

	var sb strings.Builder
	for i, s := range sections {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(strings.TrimRight(s.Raw, "\n"))
	}
	sb.WriteString("\n")
	return sb.String()
}
