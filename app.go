package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/0xF1dev/librecards/backend"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Dialog struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Buttons []string `json:"buttons"`
	Default string   `json:"default"`
	Cancel  string   `json:"cancel"`
}

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	backend.StartupChecks()
}

func (a *App) ReadSettings() backend.Settings {
	return backend.ReadSettings()
}

func (a *App) SaveSettings(settings backend.Settings) int {
	return backend.SaveSettings(settings)
}

func (a *App) GetCards() backend.Index {
	return backend.ReadIndex()
}

func (a *App) NewCard(data backend.CardData) int {
	return backend.CreateNewCard(data)
}

func (a *App) GetCard(id string) backend.Card {
	return backend.GetCard(id)
}

func (a *App) SaveError(id string, section int, question int) int {
	return backend.SaveError(id, section, question)
}

func (a *App) SaveCorrect(id string, section int, question int) int {
	return backend.SaveCorrect(id, section, question)
}

func (a *App) GetErrors(id string) backend.CardErrors {
	return backend.GetErrors(id)
}

func (a *App) DeleteCard(id string, dialog Dialog) int {
	choice, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         dialog.Title,
		Message:       dialog.Content,
		Buttons:       dialog.Buttons,
		DefaultButton: dialog.Default,
		CancelButton:  dialog.Cancel,
	})
	if err != nil {
		log.Printf("[!] Error while showing dialog: %s\n", err.Error())
		return 1
	}
	if choice == dialog.Buttons[0] {
		return backend.DeleteCard(id)
	} else {
		return 1
	}
}

func (a *App) UpdateCard(id string, card backend.CardData) int {
	return backend.UpdateCard(id, card)
}

func (a *App) ExportCards(ids []string, dialogTitle string) int {
	now := time.Now()
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: fmt.Sprintf("%d-%d-%d_%d:%d:%d.lcf", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second()),
		Filters: []runtime.FileFilter{
			{DisplayName: "Librecards Flashcard Collection (*.lcf)", Pattern: "*.lcf"},
		},
		Title: dialogTitle,
	})
	if err != nil {
		log.Printf("[!] Error while showing export dialog: %s\n", err.Error())
		return 1
	}

	if path == "" {
		return 0
	} else {
		return backend.ExportCards(ids, path)
	}
}
