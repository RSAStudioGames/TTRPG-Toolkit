package formula

import "testing"

func TestParseValid(t *testing.T) {
	valid := []string{
		"2d6 + {strength_mod}",
		"max(1, 2, 3)",
		"clamp(5, 0, 10)",
		"2^3",
		"1 == 1",
		"4d6kh3",
		"1d20r1",
		"(1d20 + 5) >= 15",
		"1d20 == 20 and 2d6 > 0",
		"not 1d20 < 10",
		"{blessed} or 1d20 >= 15",
		"if 1d20 == 20 then 2d6 else 1d6",
		"if 1d20 >= 15 then 10 else 5",
		"floor(2)",
		"abs(3)",
	}
	for _, input := range valid {
		t.Run(input, func(t *testing.T) {
			_, err := Parse(input)
			if err != nil {
				t.Fatalf("Parse(%q): %v; errors: %v", input, err, FormatErrors(err))
			}
		})
	}
}

func TestParseInvalid(t *testing.T) {
	invalid := []string{
		"2d",
		"{unclosed",
		"unknown(1)",
		"clamp(1)",
		"max(1)",
		"if 1d20 then 2d6",
		"if 1d20 == 20 then 2d6",
		"1d20 and",
		"not",
		"",
	}
	for _, input := range invalid {
		t.Run(input, func(t *testing.T) {
			_, err := Parse(input)
			if err == nil {
				t.Fatalf("Parse(%q) expected error", input)
			}
		})
	}
}

func TestFormatErrors(t *testing.T) {
	_, err := Parse("2d + 1")
	if err == nil {
		t.Fatal("expected parse error")
	}
	msgs := FormatErrors(err)
	if len(msgs) == 0 {
		t.Fatal("expected formatted errors")
	}
}
