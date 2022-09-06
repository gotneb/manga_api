package web

import (
	"fmt"
	"strings"
)

type Manga struct {
	Title         string   `json:"title"`
	Author        string   `json:"author"`
	Thumbnail     string   `json:"thumbnail"`
	Tags          []string `json:"tags"`
	Description   string   `json:"description"`
	Situation     string   `json:"situation"`
	TotalChapters int      `json:"total_chapters"`
	Chapters      []string `json:"chapters"`
}

type Chapter struct {
	Title string `json:"title"`
	// It's supossed to be a float64, but rarely some chapters are found without a "floating point number"
	// eg.: Black Cover - 40v01 - o.o
	Value      string   `json:"value"`
	TotalPages int      `json:"total_pages"`
	Pages      []string `json:"pages"`
}

// Show details about the manga
func (m *Manga) Show() {
	fmt.Println("Title: " + m.Title)
	fmt.Printf("Author: %s\n", m.Author)
	fmt.Printf("Situation: %s\n", m.Situation)
	fmt.Printf("Total chapters: %d\n", m.TotalChapters)
	fmt.Println("Thumbnail: " + m.Thumbnail)
	fmt.Printf("Tags: ")
	for _, v := range m.Tags {
		fmt.Printf("%s, ", v)
	}
	fmt.Println()
	fmt.Println("Description: " + m.Description)
	// TODO: Show all chapters
}

// Returns if any field is empty
func (m *Manga) IsEmpty() bool {
	return m.Description == "" || m.Situation == "" || m.Thumbnail == "" || m.Title == ""
}

/*
Formats manga name to another one whose website is able to reach
Example: From: "Huge: Stupid Large    NAME" to "huge-stupid-large-name"
*/
func FormatedTitle(title string) string {
	const notAllowedSymbols = "’.:!@'\"`~#$%&()[]{}"

	nameFormated := strings.Clone(title)
	nameFormated = strings.ReplaceAll(strings.ToLower(nameFormated), " ", "-")
	for _, symbol := range notAllowedSymbols {
		if strings.Contains(nameFormated, string(symbol)) {
			nameFormated = strings.ReplaceAll(nameFormated, string(symbol), "")
		}
	}

	return nameFormated
}
