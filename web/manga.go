package web

import "fmt"

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
