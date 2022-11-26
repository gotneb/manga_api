package web

import (
	"fmt"
	"strings"
)

type Manga struct {
	Title         string          `json:"title"`
	Author        string          `json:"author"`
	Thumbnail     string          `json:"thumbnail"`
	Genres        []string        `json:"genres"`
	Summary       string          `json:"summary"`
	Status        string          `json:"status"`
	TotalChapters int             `json:"total_chapters"`
	Chapters      []SimpleChapter `json:"chapters"`
}

/*
 * Each simple chapter has the following format:
 * value: release_date (year-month-day)
 * e.g => 4: 2022-09-03 ; 150: 2016-12-25
 */
type SimpleChapter struct {
	Value string `json:"value"`
	Date  string `json:"date"`
}

type Chapter struct {
	Title string `json:"title"`
	// It's suposed to be a float64, but rarely some chapters are found without a "floating point number"
	// eg.: Black Cover - 40v01 - o.o
	Value      string   `json:"value"`
	TotalPages int      `json:"total_pages"`
	Pages      []string `json:"pages"`
}

// Show details about the manga
func (m *Manga) Show() {
	fmt.Println("Title: " + m.Title)
	fmt.Printf("Author: %s\n", m.Author)
	fmt.Printf("Status: %s\n", m.Status)
	fmt.Printf("Total chapters: %d\n", m.TotalChapters)
	fmt.Println("Thumbnail: " + m.Thumbnail)
	fmt.Printf("Tags: ")
	for _, v := range m.Genres {
		fmt.Printf("%s, ", v)
	}
	fmt.Println()
	fmt.Println("Description: " + m.Summary)
	// TODO: Show all chapters
}

// Returns if any field is empty
func (m *Manga) IsEmpty() bool {
	return m.Summary == "" || m.Status == "" || m.Thumbnail == "" || m.Title == ""
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
