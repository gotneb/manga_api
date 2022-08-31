package web

import "fmt"

type Manga struct {
	Title         string             `json:"title"`
	Author        string             `json:"author"`
	Tags          []string           `json:"tags"`
	TotalChapters int                `json:"totalChapters"`
	Description   string             `json:"description"`
	Thumbnail     string             `json:"thumbnail"`
	Situation     string             `json:"situation"`
	Chapters      map[float64]string `json:"chapters"`
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
