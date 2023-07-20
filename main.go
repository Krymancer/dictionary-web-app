package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type APIResponse []struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic"`
	Phonetics []struct {
		Text      string `json:"text"`
		Audio     string `json:"audio"`
		SourceURL string `json:"sourceUrl,omitempty"`
		License   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"license,omitempty"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
			Synonyms   []any  `json:"synonyms"`
			Antonyms   []any  `json:"antonyms"`
			Example    string `json:"example,omitempty"`
		} `json:"definitions"`
		Synonyms []string `json:"synonyms"`
		Antonyms []string `json:"antonyms"`
	} `json:"meanings"`
	License struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"license"`
	SourceUrls []string `json:"sourceUrls"`
}

func main() {
	resp, err := http.Get("https://api.dictionaryapi.dev/api/v2/entries/en/beauty")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(body))

	index := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("pages/index.html"))
		tmpl.Execute(w, nil)
	}

	word := func(w http.ResponseWriter, r *http.Request) {
		baseUrl := "https://api.dictionaryapi.dev/api/v2/entries/en/"
		resp, err := http.Get(baseUrl + r.URL.Query().Get("word"))

		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}

		var response APIResponse

		err = json.Unmarshal(body, &response)

		if err != nil {
			panic(err)
		}

		fmt.Print(response)
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/word", word)

	log.Fatal(http.ListenAndServe(":42069", nil))
}
