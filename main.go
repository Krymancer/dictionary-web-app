package main

import (
	"encoding/json"
	"flag"
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
	port := flag.Int("port", -1, "specify a port")
	flag.Parse()

	portStr := "8080"

	if *port != -1 {
		portStr = fmt.Sprintf(":%d", *port)
	}

	index := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		template := template.Must(template.ParseFiles("pages/index.html"))
		template.Execute(w, nil)
	}

	word := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
			return
		}

		word := r.Form.Get("word")

		baseUrl := "https://api.dictionaryapi.dev/api/v2/entries/en/"
		url := baseUrl + word
		fmt.Println(url)
		resp, err := http.Get(url)

		if err != nil {
			http.Error(w, "Failed to get word from dictionary API", http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			http.Error(w, "Failed to get the API response body", http.StatusInternalServerError)
			return
		}

		var response APIResponse

		err = json.Unmarshal(body, &response)

		if err != nil {
			http.Error(w, "Failed to parse the body data", http.StatusInternalServerError)
			return
		}

		template := template.Must(template.ParseFiles("components/word.html"))
		template.Execute(w, response[0])
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/word", word)

	log.Fatal(http.ListenAndServe(portStr, nil))
}
