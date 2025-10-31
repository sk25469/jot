package styles

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestGetTagStyle(t *testing.T) {
	// Test that same tag always gets same style
	style1 := GetTagStyle("golang")
	style2 := GetTagStyle("golang")

	// We can't directly compare styles, but we can test consistency
	// by checking that the same input produces the same output
	rendered1 := style1.Render("golang")
	rendered2 := style2.Render("golang")

	if rendered1 != rendered2 {
		t.Errorf("Same tag should produce same styled output")
	}

	// Test different tags get different styles
	styleA := GetTagStyle("tagA")
	styleB := GetTagStyle("tagB")

	renderedA := styleA.Render("tagA")
	renderedB := styleB.Render("tagB")

	// They should be different (very unlikely to be the same with different inputs)
	if renderedA == renderedB {
		t.Logf("Different tags produced same output (rare but possible): %s", renderedA)
	}
}

func TestGetModeStyle(t *testing.T) {
	testCases := []struct {
		mode     string
		expected string // We'll check if it contains the mode text
	}{
		{"dev", "dev"},
		{"journal", "journal"},
		{"meeting", "meeting"},
		{"unknown", "unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.mode, func(t *testing.T) {
			style := GetModeStyle(tc.mode)
			rendered := style.Render(tc.mode)

			// The rendered output should contain the mode text
			if !strings.Contains(rendered, tc.expected) {
				t.Errorf("Mode style for %s should contain %s, got: %s", tc.mode, tc.expected, rendered)
			}

			// Should not be empty
			if rendered == "" {
				t.Errorf("Mode style for %s should not be empty", tc.mode)
			}
		})
	}
}

func TestRenderTags(t *testing.T) {
	testCases := []struct {
		name     string
		tags     []string
		expected string
	}{
		{
			name:     "empty tags",
			tags:     []string{},
			expected: "",
		},
		{
			name:     "single tag",
			tags:     []string{"golang"},
			expected: "golang", // Should contain the tag name
		},
		{
			name:     "multiple tags",
			tags:     []string{"golang", "test", "cli"},
			expected: "", // We'll check that all tags are present
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := RenderTags(tc.tags)

			if tc.name == "empty tags" {
				if result != tc.expected {
					t.Errorf("Expected empty string for empty tags, got: %s", result)
				}
				return
			}

			if tc.name == "single tag" {
				if !strings.Contains(result, tc.expected) {
					t.Errorf("Result should contain %s, got: %s", tc.expected, result)
				}
				return
			}

			if tc.name == "multiple tags" {
				// Check that all tags are present in the output
				for _, tag := range tc.tags {
					if !strings.Contains(result, tag) {
						t.Errorf("Result should contain tag %s, got: %s", tag, result)
					}
				}
			}
		})
	}
}

func TestRenderHeader(t *testing.T) {
	testCases := []struct {
		title    string
		expected string
	}{
		{"Test Title", "Test Title"},
		{"", ""},
		{"Long Title With Spaces", "Long Title With Spaces"},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			result := RenderHeader(tc.title)

			// Should contain the title text and the emoji
			if !strings.Contains(result, tc.expected) {
				t.Errorf("Header should contain %s, got: %s", tc.expected, result)
			}

			if !strings.Contains(result, "📝") {
				t.Errorf("Header should contain emoji 📝, got: %s", result)
			}
		})
	}
}

func TestRenderSeparator(t *testing.T) {
	result := RenderSeparator()

	// Should not be empty
	if result == "" {
		t.Errorf("Separator should not be empty")
	}

	// Should contain dashes or similar separator characters
	if !strings.Contains(result, "─") {
		t.Errorf("Separator should contain separator characters, got: %s", result)
	}
}

func TestRenderProgress(t *testing.T) {
	testCases := []struct {
		name     string
		current  int
		total    int
		expected string
	}{
		{
			name:     "zero total",
			current:  5,
			total:    0,
			expected: "",
		},
		{
			name:     "50 percent",
			current:  5,
			total:    10,
			expected: "5/10",
		},
		{
			name:     "100 percent",
			current:  10,
			total:    10,
			expected: "10/10",
		},
		{
			name:     "zero current",
			current:  0,
			total:    10,
			expected: "0/10",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := RenderProgress(tc.current, tc.total)

			if tc.name == "zero total" {
				if result != tc.expected {
					t.Errorf("Expected empty string for zero total, got: %s", result)
				}
				return
			}

			// Should contain the expected progress text
			if !strings.Contains(result, tc.expected) {
				t.Errorf("Progress should contain %s, got: %s", tc.expected, result)
			}

			// Should contain progress bar characters
			if !strings.Contains(result, "█") && !strings.Contains(result, "░") {
				t.Errorf("Progress should contain progress bar characters, got: %s", result)
			}
		})
	}
}

func TestRenderBox(t *testing.T) {
	title := "Test Box"
	content := "This is test content"

	result := RenderBox(title, content)

	// Should not be empty
	if result == "" {
		t.Errorf("Box should not be empty")
	}

	// Should contain title and content
	if !strings.Contains(result, title) {
		t.Errorf("Box should contain title %s, got: %s", title, result)
	}

	if !strings.Contains(result, content) {
		t.Errorf("Box should contain content %s, got: %s", content, result)
	}

	// Should contain box border characters
	if !strings.Contains(result, "╭") && !strings.Contains(result, "│") {
		t.Errorf("Box should contain border characters, got: %s", result)
	}
}

func TestColorDefinitions(t *testing.T) {
	// Test that colors are defined and not empty
	colors := []struct {
		name  string
		color lipgloss.Color
	}{
		{"Primary", Primary},
		{"Secondary", Secondary},
		{"Accent", Accent},
		{"Success", Success},
		{"Warning", Warning},
		{"Error", Error},
		{"Subtle", Subtle},
		{"Muted", Muted},
		{"Text", Text},
	}

	for _, c := range colors {
		t.Run(c.name, func(t *testing.T) {
			if string(c.color) == "" {
				t.Errorf("Color %s should not be empty", c.name)
			}

			// Colors should start with # (hex colors)
			if !strings.HasPrefix(string(c.color), "#") {
				t.Errorf("Color %s should be a hex color, got: %s", c.name, string(c.color))
			}
		})
	}
}

func TestStyleDefinitions(t *testing.T) {
	// Test that key styles are defined
	styles := []struct {
		name  string
		style lipgloss.Style
	}{
		{"TitleStyle", TitleStyle},
		{"SubtitleStyle", SubtitleStyle},
		{"IDStyle", IDStyle},
		{"DateStyle", DateStyle},
		{"ContentStyle", ContentStyle},
		{"ErrorStyle", ErrorStyle},
		{"WarningStyle", WarningStyle},
		{"SuccessStyle", SuccessStyle},
	}

	for _, s := range styles {
		t.Run(s.name, func(t *testing.T) {
			// Test that we can render with the style without panicking
			testText := "test"
			result := s.style.Render(testText)

			// Should contain the test text
			if !strings.Contains(result, testText) {
				t.Errorf("Style %s should render text correctly, got: %s", s.name, result)
			}
		})
	}
}
