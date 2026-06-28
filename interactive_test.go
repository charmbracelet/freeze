package main

import "testing"

func TestEmbeddedFontFamilySelection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		ligatures bool
		want      string
	}{
		{
			name:      "ligatures enabled",
			ligatures: true,
			want:      embeddedFontJetBrainsMono,
		},
		{
			name:      "ligatures disabled",
			ligatures: false,
			want:      embeddedFontJetBrainsMonoNL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := &Config{
				Font: Font{Ligatures: tt.ligatures},
			}
			if got := embeddedFontFamilySelection(config); got != tt.want {
				t.Fatalf("embeddedFontFamilySelection() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestApplyEmbeddedFontFamily(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		selection string
		want      Font
	}{
		{
			name:      "jetbrains mono",
			selection: embeddedFontJetBrainsMono,
			want: Font{
				Family:    embeddedFontJetBrainsMono,
				Ligatures: true,
			},
		},
		{
			name:      "jetbrains mono nl",
			selection: embeddedFontJetBrainsMonoNL,
			want: Font{
				Family:    embeddedFontJetBrainsMono,
				Ligatures: false,
			},
		},
		{
			name:      "unknown selection defaults to mono",
			selection: "Fira Code",
			want: Font{
				Family:    embeddedFontJetBrainsMono,
				Ligatures: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := &Config{
				Font: Font{
					Family:    "custom",
					Ligatures: false,
				},
			}
			applyEmbeddedFontFamily(tt.selection, config)
			if config.Font != tt.want {
				t.Fatalf("applyEmbeddedFontFamily() = %+v, want %+v", config.Font, tt.want)
			}
		})
	}
}
