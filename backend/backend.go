package backend

import (
	"cmp"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"slices"

	"github.com/adrg/xdg"
)

type Index struct {
	Count int         `json:"count"`
	Cards []CardEntry `json:"cards"`
}

type CardEntry struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Subject string `json:"subject"`
	Path    string `json:"path"`
}

type CardData struct {
	Title    string    `json:"title"`
	Subject  string    `json:"subject"`
	Sections []Section `json:"sections"`
}

type Section struct {
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
}

type Question struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Card struct {
	ID   string   `json:"id"`
	Data CardData `json:"data"`
}

type CardErrors struct {
	ID     string          `json:"id"`
	Title  string          `json:"title"`
	Errors []SectionErrors `json:"errors"`
}

type SectionErrors struct {
	Section   int   `json:"section"`
	Questions []int `json:"questions"`
}

func StartupChecks() error {
	log.Println("------- Starting checks... -------")

	if _, err := os.Stat(xdg.DataHome + "/librecards"); os.IsNotExist(err) {
		log.Println("Librecards data directory doesn't exist. Creating...")
		err := os.Mkdir(xdg.DataHome+"/librecards", 0777)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(xdg.DataHome + "/librecards/index.json"); os.IsNotExist(err) {
		log.Println("Librecards index file doesn't exist. Creating...")
		f, err := os.Create(xdg.DataHome + "/librecards/index.json")
		if err != nil {
			return err
		}
		defer f.Close()

		f.Write([]byte(`{"count":0,"cards":[]}`))
	}

	if _, err := os.Stat(xdg.DataHome + "/librecards/errors.json"); os.IsNotExist(err) {
		log.Println("Librecards errors file doesn't exist. Creating...")
		f, err := os.Create(xdg.DataHome + "/librecards/errors.json")
		if err != nil {
			return err
		}
		defer f.Close()

		f.Write([]byte(`[]`))
	}

	log.Println("------- Checks finished! -------")

	return nil
}

var DataPath = xdg.DataHome + "/librecards/"

var IndexPath = xdg.DataHome + "/librecards/index.json"
var ErrorsPath = xdg.DataHome + "/librecards/errors.json"

var ErrorCardEntry = CardEntry{ID: "0", Title: "Error", Subject: "Try again"}
var ErrorCard = Card{ID: "0", Data: CardData{Sections: []Section{{Questions: []Question{{Question: "Error", Answer: "Try again"}}}}}}
var ErrorCardErrors = CardErrors{ID: "0", Title: "Error", Errors: []SectionErrors{{Section: 0, Questions: []int{0}}}}

func ReadIndex() Index {
	log.Println("--- Reading index... ---")

	f, err := os.ReadFile(IndexPath)
	if err != nil {
		log.Printf("[!] Error while reading index file: %s\n", err.Error())
		return Index{Count: 0, Cards: []CardEntry{ErrorCardEntry}}
	}

	var index Index
	err = json.Unmarshal(f, &index)
	if err != nil {
		log.Printf("[!] Error while unmarshalling index: %s\n", err.Error())
		return Index{Count: 0, Cards: []CardEntry{ErrorCardEntry}}
	}

	log.Printf("Total entries: %d", index.Count)

	log.Println("--- Index read! ---")

	return index
}

func SaveIndex(index Index) error {
	data, err := json.Marshal(index)
	if err != nil {
		log.Printf("[!] Error while marshalling index: %s\n", err.Error())
		return err
	}

	if err = os.WriteFile(IndexPath, data, 0777); err != nil {
		log.Printf("[!] Error while saving index: %s\n", err.Error())
		return err
	}

	return nil
}

func GenerateNewID() string {
	charset := "0123456789abcdefghijklmnopqrstuvwxyz"
	length := 8
	id := ""
	for range length {
		char := string(charset[rand.Intn(len(charset)-1)])
		id += char
	}
	return id
}

func CreateNewCard(card CardData) int {
	log.Println("--- Creating new card... ---")

	index := ReadIndex()

	id := GenerateNewID()
	path := DataPath + id + ".json"

	index.Cards = append(index.Cards, CardEntry{ID: id, Title: card.Title, Subject: card.Subject, Path: path})
	index.Count += 1

	log.Println(card)

	data, err := json.Marshal(Card{ID: id, Data: card})
	if err != nil {
		return 1
	}
	if err = os.WriteFile(path, data, 0777); err != nil {
		return 1
	}

	if err = SaveIndex(index); err != nil {
		return 1
	}

	log.Printf("--- Card created! (#%s) ---\n", id)

	return 0
}

func GetCard(id string) Card {
	log.Printf("--- Reading card... (#%s) ---\n", id)

	if id == "0" {
		return ErrorCard
	}

	index := ReadIndex()

	i := slices.IndexFunc(index.Cards, func(c CardEntry) bool { return c.ID == id })

	if i == -1 {
		log.Printf("[!] Unable to get card: the card does not exist\n")
		return ErrorCard
	}

	f, err := os.ReadFile(index.Cards[i].Path)
	if err != nil {
		log.Printf("[!] Unable to read card: %s\n", err.Error())
		return ErrorCard
	}

	var card Card
	err = json.Unmarshal(f, &card)
	if err != nil {
		log.Printf("[!] Unable to unmarshal card: %s\n", err.Error())
		return ErrorCard
	}

	log.Printf("--- Card read! (#%s) ---\n", id)

	return card
}

func GetErrors(id string) CardErrors {
	log.Printf("--- Obtaining card errors... (#%s) ---\n", id)

	f, err := os.ReadFile(ErrorsPath)
	if err != nil {
		log.Printf("[!] Unable to read errors: %s\n", err.Error())
		return ErrorCardErrors
	}

	var errors []CardErrors
	err = json.Unmarshal(f, &errors)
	if err != nil {
		log.Printf("[!] Unable to unmarshal errors: %s\n", err.Error())
		return ErrorCardErrors
	}

	i := slices.IndexFunc(errors, func(c CardErrors) bool { return c.ID == id })

	if i == -1 {
		return CardErrors{ID: id, Title: GetCard(id).Data.Title, Errors: []SectionErrors{}}
	}

	slices.SortFunc(errors[i].Errors, func(a, b SectionErrors) int {
		return cmp.Compare(a.Section, b.Section)
	})

	for _, s := range errors[i].Errors {
		slices.SortFunc(s.Questions, func(a, b int) int {
			return cmp.Compare(a, b)
		})
	}

	log.Printf("--- Errors obtained! (#%s) ---\n", id)

	return errors[i]
}

func SaveError(id string, section int, question int) int {
	log.Printf("--- Saving errors... (#%s) ---\n", id)

	f, err := os.ReadFile(ErrorsPath)
	if err != nil {
		log.Printf("[!] Unable to read errors while saving: %s\n", err.Error())
		return 1
	}

	var errors []CardErrors
	err = json.Unmarshal(f, &errors)
	if err != nil {
		log.Printf("[!] Unable to unmarshal errors: %s\n", err.Error())
		return 1
	}

	i := slices.IndexFunc(errors, func(c CardErrors) bool { return c.ID == id })
	if i == -1 {
		errors = append(errors, CardErrors{ID: id, Title: GetCard(id).Data.Title, Errors: []SectionErrors{{Section: section, Questions: []int{question}}}})
		data, err := json.Marshal(errors)
		if err != nil {
			log.Printf("[!] Unable to marshal errors: %s\n", err.Error())
			return 1
		}

		if err = os.WriteFile(ErrorsPath, data, 0777); err != nil {
			log.Printf("[!] Unable to save errors: %s\n", err.Error())
			return 1
		}

		log.Printf("--- Errors saved! (#%s) ---\n", id)

		return 0
	}

	i2 := slices.IndexFunc(errors[i].Errors, func(s SectionErrors) bool { return s.Section == section })
	if i2 == -1 {
		errors[i].Errors = append(errors[i].Errors, SectionErrors{Section: section, Questions: []int{question}})
		data, err := json.Marshal(errors)
		if err != nil {
			log.Printf("[!] Unable to marshal errors: %s\n", err.Error())
			return 1
		}

		if err = os.WriteFile(ErrorsPath, data, 0777); err != nil {
			log.Printf("[!] Unable to save errors: %s\n", err.Error())
			return 1
		}

		log.Printf("--- Errors saved! (#%s) ---\n", id)

		return 0
	}

	if !slices.Contains(errors[i].Errors[i2].Questions, question) {
		errors[i].Errors[i2].Questions = append(errors[i].Errors[i2].Questions, question)
	}

	data, err := json.Marshal(errors)
	if err != nil {
		log.Printf("[!] Unable to marshal errors: %s\n", err.Error())
		return 1
	}

	if err = os.WriteFile(ErrorsPath, data, 0777); err != nil {
		log.Printf("[!] Unable to save errors: %s\n", err.Error())
		return 1
	}

	log.Printf("--- Errors saved! (#%s) ---\n", id)

	return 0
}

func SaveCorrect(id string, section int, question int) int {
	log.Printf("--- Saving error correction... (#%s) ---\n", id)

	f, err := os.ReadFile(ErrorsPath)
	if err != nil {
		log.Printf("[!] Unable to read errors while saving: %s\n", err.Error())
		return 1
	}

	var errors []CardErrors
	err = json.Unmarshal(f, &errors)
	if err != nil {
		log.Printf("[!] Unable to unmarshal errors: %s\n", err.Error())
		return 1
	}

	i := slices.IndexFunc(errors, func(c CardErrors) bool { return c.ID == id })
	if i == -1 {
		log.Printf("--- Error correction saved! (#%s) ---\n", id)
		return 0
	}

	i2 := slices.IndexFunc(errors[i].Errors, func(s SectionErrors) bool { return s.Section == section })
	if i2 == -1 {
		log.Printf("--- Error correction saved! (#%s) ---\n", id)
		return 0
	}

	if slices.Contains(errors[i].Errors[i2].Questions, question) {
		errI := slices.IndexFunc(errors[i].Errors[i2].Questions, func(e int) bool { return e == question })
		errors[i].Errors[i2].Questions = slices.Delete(errors[i].Errors[i2].Questions, errI, errI+1)
	}

	if len(errors[i].Errors[i2].Questions) == 0 {
		errors[i].Errors = slices.Delete(errors[i].Errors, i2, i2+1)
	}
	if len(errors[i].Errors) == 0 {
		errors = slices.Delete(errors, i, i+1)
	}

	data, err := json.Marshal(errors)
	if err != nil {
		log.Printf("[!] Unable to marshal errors: %s\n", err.Error())
		return 1
	}

	if err = os.WriteFile(ErrorsPath, data, 0777); err != nil {
		log.Printf("[!] Unable to save errors: %s\n", err.Error())
		return 1
	}

	log.Printf("--- Error correction saved! (#%s) ---\n", id)

	return 0
}

func DeleteCard(id string) int {
	log.Printf("--- Deleting card... (#%s) ---\n", id)

	index := ReadIndex()

	i := slices.IndexFunc(index.Cards, func(c CardEntry) bool { return c.ID == id })
	if i == -1 {
		log.Printf("[!] Unable to delete card: card wasn't found\n")
		return 1
	}

	if err := os.Remove(index.Cards[i].Path); err != nil {
		log.Printf("[!] Unable to delete card: %s\n", err.Error())
		return 1
	}

	index.Cards = slices.Delete(index.Cards, i, i+1)
	index.Count -= 1

	if err := SaveIndex(index); err != nil {
		return 1
	}

	log.Printf("--- Card deleted! (#%s) ---\n", id)

	return 0
}

func UpdateCard(id string, card CardData) int {
	log.Printf("--- Updating card... (#%s) ---\n", id)

	index := ReadIndex()

	i := slices.IndexFunc(index.Cards, func(c CardEntry) bool { return c.ID == id })
	if i == -1 {
		log.Printf("[!] Unable to update card: the card wasn't found\n")
		return 1
	}

	index.Cards[i].Title = card.Title
	index.Cards[i].Subject = card.Subject

	f, err := os.Create(index.Cards[i].Path)
	if err != nil {
		log.Printf("[!] Unable to update/create card file: %s\n", err.Error())
		return 1
	}
	defer f.Close()

	data, err := json.Marshal(Card{ID: id, Data: card})
	if err != nil {
		log.Printf("[!] Unable to marshal card while updating: %s\n", err.Error())
		return 1
	}

	_, err = f.Write(data)
	if err != nil {
		log.Printf("[!] Unable to save card while updating: %s\n", err.Error())
		return 1
	}

	if err = SaveIndex(index); err != nil {
		return 1
	}

	log.Printf("--- Card updated! (#%s) ---\n", id)

	return 0
}
