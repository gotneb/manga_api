package web

import "fmt"

type Manga struct {
	Title       string
	Author      string
	Tags        []string
	Chapters    int
	Description string
	Thumbnail   string
	Situation   string
}

// Show details about the manga
func (m *Manga) Show() {
	fmt.Println("Title: " + m.Title)
	fmt.Printf("Author: %s\n", m.Author)
	fmt.Printf("Situation: %s\n", m.Situation)
	fmt.Printf("Chapters: %d\n", m.Chapters)
	fmt.Println("Thumbnail: " + m.Thumbnail)
	fmt.Printf("Tags: ")
	for _, v := range m.Tags {
		fmt.Printf("%s, ", v)
	}
	fmt.Println()
	fmt.Println("Description: " + m.Description)
}
