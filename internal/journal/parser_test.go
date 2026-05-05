package journal

import (
	"strings"
	"testing"
)

func TestParseMarkdown_Empty(t *testing.T) {
	sections := ParseMarkdown("")
	if len(sections) != 0 {
		t.Errorf("expected 0 sections, got %d", len(sections))
	}

	sections = ParseMarkdown("   \n\n  ")
	if len(sections) != 0 {
		t.Errorf("expected 0 sections for whitespace, got %d", len(sections))
	}
}

func TestParseMarkdown_SingleSection(t *testing.T) {
	src := "# 09:55\nhello world\n"
	sections := ParseMarkdown(src)

	if len(sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(sections))
	}
	if sections[0].Heading != "# 09:55" {
		t.Errorf("expected heading '# 09:55', got %q", sections[0].Heading)
	}
}

func TestParseMarkdown_MultipleSections(t *testing.T) {
	src := "# 09:00\nmorning entry\n\n# 10:00\nevening entry\n"
	sections := ParseMarkdown(src)

	if len(sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(sections))
	}
	if sections[0].Heading != "# 09:00" {
		t.Errorf("expected heading '# 09:00', got %q", sections[0].Heading)
	}
	if sections[1].Heading != "# 10:00" {
		t.Errorf("expected heading '# 10:00', got %q", sections[1].Heading)
	}
}

func TestParseMarkdown_CodeBlockWithHeading(t *testing.T) {
	src := "# 09:00\n```go\n# this is NOT a heading\n```\nbody text\n"
	sections := ParseMarkdown(src)

	if len(sections) != 1 {
		t.Fatalf("expected 1 section (code block # should be ignored), got %d", len(sections))
	}
	if sections[0].Heading != "# 09:00" {
		t.Errorf("expected heading '# 09:00', got %q", sections[0].Heading)
	}
	if !strings.Contains(sections[0].Raw, "```go") {
		t.Error("expected code block to be preserved in Raw")
	}
}

func TestRenderHeader_Default(t *testing.T) {
	result := RenderHeader("# {time}\n", "09:55")
	if result != "# 09:55" {
		t.Errorf("expected '# 09:55', got %q", result)
	}
}

func TestRenderHeader_WithSeparator(t *testing.T) {
	result := RenderHeader("# {time}\n---\n", "09:55")
	if result != "# 09:55\n---" {
		t.Errorf("expected '# 09:55\\n---', got %q", result)
	}
}

func TestAppendToSections_NewSection(t *testing.T) {
	var sections []Section
	result := AppendToSections(sections, "09:55", "# {time}\n", false, "hello")

	if !strings.Contains(result, "# 09:55") {
		t.Errorf("expected heading '# 09:55', got %q", result)
	}
	if !strings.Contains(result, "hello") {
		t.Errorf("expected content 'hello', got %q", result)
	}
}

func TestAppendToSections_SameTime(t *testing.T) {
	sections := ParseMarkdown("# 09:55\nfirst entry\n")
	result := AppendToSections(sections, "09:55", "# {time}\n", true, "second entry")

	count := strings.Count(result, "# 09:55")
	if count != 1 {
		t.Errorf("expected 1 heading, got %d", count)
	}
	if !strings.Contains(result, "first entry") {
		t.Errorf("expected 'first entry', got %q", result)
	}
	if !strings.Contains(result, "second entry") {
		t.Errorf("expected 'second entry', got %q", result)
	}
}

func TestAppendToSections_DifferentTime(t *testing.T) {
	sections := ParseMarkdown("# 09:00\npast entry\n")
	result := AppendToSections(sections, "10:00", "# {time}\n", false, "current entry")

	count := strings.Count(result, "#")
	if count != 2 {
		t.Errorf("expected 2 headings, got %d", count)
	}
	if !strings.Contains(result, "# 09:00") {
		t.Errorf("expected old heading, got %q", result)
	}
	if !strings.Contains(result, "# 10:00") {
		t.Errorf("expected new heading, got %q", result)
	}
}

func TestEntryFile_Format(t *testing.T) {
	file, timeValue := EntryFile("/tmp/journal", "24", "02.md")
	if file == "" {
		t.Error("expected non-empty file path")
	}
	if timeValue == "" {
		t.Error("expected non-empty time value")
	}
}
