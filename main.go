// rethink naming so I get rid of types in names

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

func handler(story map[string]ChapterType) http.HandlerFunc {
	templ := template.Must(template.ParseFiles("layout.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		if chapter, ok := story[strings.Trim(r.URL.Path, "/")]; ok {
			templ.Execute(w, chapter)
		} else {
			templ.Execute(w, story["intro"])
		}
	}
}

func readStory(storyPath string) (map[string]ChapterType, error) {
	storyBytes, error := ioutil.ReadFile(storyPath)
	if error != nil {
		return nil, error
	}
	var storyMap map[string]ChapterType
	_ = json.Unmarshal(storyBytes, &storyMap)
	return storyMap, nil
}

// ChapterType defines the structure of an adventure chapter
type ChapterType struct {
	Title   string        `json:"title,omitempty"`
	Story   []string      `json:"story,omitempty"`
	Options []OptionsType `json:"options,omitempty"`
}

// OptionsType defines the structure of the options to choose the next chapter
type OptionsType struct {
	Chapter string `json:"arc,omitempty"`
	Text    string `json:"text,omitempty"`
}

func main() {
	story, error := readStory("gopher.json")
	if error != nil {
		fmt.Println(error)
		return
	}
	h := handler(story)
	http.ListenAndServe(":8080", h)
}
