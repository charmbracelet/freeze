package main

import "github.com/alecthomas/chroma"

var style = chroma.MustNewStyle("charm", chroma.StyleEntries{
	chroma.Text:                "#C4C4C4",
	chroma.Error:               "#F1F1F1 bg:#F05B5B",
	chroma.Comment:             "#676767",
	chroma.CommentPreproc:      "#FF875F",
	chroma.Keyword:             "#00AAFF",
	chroma.KeywordReserved:     "#FF5FD2",
	chroma.KeywordNamespace:    "#FF5F87",
	chroma.KeywordType:         "#6E6ED8",
	chroma.Operator:            "#EF8080",
	chroma.Punctuation:         "#E8E8A8",
	chroma.Name:                "#C4C4C4",
	chroma.NameBuiltin:         "#FF8EC7",
	chroma.NameTag:             "#B083EA",
	chroma.NameAttribute:       "#7A7AE6",
	chroma.NameClass:           "#F1F1F1 underline bold",
	chroma.NameDecorator:       "#FFFF87",
	chroma.NameFunction:        "#00D787",
	chroma.LiteralNumber:       "#6EEFC0",
	chroma.LiteralString:       "#C69669",
	chroma.LiteralStringEscape: "#AFFFD7",
	chroma.GenericDeleted:      "#FD5B5B",
	chroma.GenericEmph:         "italic",
	chroma.GenericInserted:     "#00D787",
	chroma.GenericStrong:       "bold",
	chroma.GenericSubheading:   "#777777",
	chroma.Background:          "bg:#373737",
})
