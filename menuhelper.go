package main

import (
	"container/list"

	"github.com/nexidian/gocliselect"
)

const (
	WINDOW_BACK   = 1
	WINDOW_SELECT = 2
)

type MenuPage struct {
	MenuObject   *gocliselect.Menu
	SelectionMap map[string]MenuSelection
	AdjacentMenu map[string]func() *MenuPage
	OnDisplay    func(thisPage *MenuPage)
}

func NewMenuPage(prompt string) *MenuPage {
	return &MenuPage{
		MenuObject:   gocliselect.NewMenu(prompt),
		SelectionMap: make(map[string]MenuSelection),
		AdjacentMenu: make(map[string]func() *MenuPage),
		OnDisplay:    nil,
	}
}

func (mp *MenuPage) AssignMenuSelection(key string, text string, execute func(key string) (int, error)) *MenuPage {
	mp.SelectionMap[key] = MenuSelection{
		SelectionKey:         key,
		SelectionText:        text,
		ExecuteMenuSelection: execute,
	}

	mp.MenuObject.AddItem(text, key)

	return mp
}

func (mp *MenuPage) ClearMenuSelections() *MenuPage {
	mp.SelectionMap = make(map[string]MenuSelection)
	mp.MenuObject.MenuItems = make([]*gocliselect.MenuItem, 0)
	return mp
}

func (mp *MenuPage) AssignAdjacentMenu(key string, adjacentMenuFactory func() *MenuPage) *MenuPage {
	mp.AdjacentMenu[key] = adjacentMenuFactory

	return mp
}

type MenuSelection struct {
	SelectionKey  string
	SelectionText string
	// A function that takes in a string and returns a MenuPage or an int
	ExecuteMenuSelection func(key string) (int, error)
}

type MenuStack struct {
	MenuPageStack *list.List
}

func NewMenuStack() *MenuStack {
	return &MenuStack{
		MenuPageStack: list.New(),
	}
}

func (ms *MenuStack) Pop() error {
	if ms.MenuPageStack.Len() == 0 {
		return nil
	}
	ms.MenuPageStack.Remove(ms.MenuPageStack.Front())
	return nil
}

func (ms *MenuStack) Push(menu *MenuPage) {
	ms.MenuPageStack.PushFront(menu)
}
