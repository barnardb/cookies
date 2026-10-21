package main

import (
	"regexp"
	"testing"
)

func TestProfileMatcher(t *testing.T) {
	tests := []struct {
		name    string
		matcher profileMatcher
		profile string
		want    bool
	}{
		{"no selection matches everything", profileMatcher{}, "Default", true},
		{"exact name", profileMatcher{names: []string{"Profile 1"}}, "Profile 1", true},
		{"exact name is not a prefix match", profileMatcher{names: []string{"Profile 1"}}, "Profile 10", false},
		{"exact name is not a regex", profileMatcher{names: []string{"Profile ."}}, "Profile 1", false},
		{"any of several names", profileMatcher{names: []string{"Work", "Default"}}, "Default", true},
		{"pattern is unanchored", profileMatcher{patterns: []*regexp.Regexp{regexp.MustCompile("ork")}}, "Work", true},
		{"anchored pattern", profileMatcher{patterns: []*regexp.Regexp{regexp.MustCompile("^Profile [0-9]$")}}, "Profile 10", false},
		{"name or pattern", profileMatcher{names: []string{"Work"}, patterns: []*regexp.Regexp{regexp.MustCompile("^Profile")}}, "Profile 2", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.matcher.matches(tt.profile); got != tt.want {
				t.Errorf("matches(%q) = %v, want %v", tt.profile, got, tt.want)
			}
		})
	}
}
