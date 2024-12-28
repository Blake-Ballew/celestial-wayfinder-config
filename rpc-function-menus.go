package main

import (
	"fmt"
	"strconv"
)

func GenerateDeleteMessagesMenu() *MenuPage {
	newPage := NewMenuPage("Delete Message")

	newPage.OnDisplay = func(thisPage *MenuPage) {
		thisPage.ClearMenuSelections()

		messagePacket, err := GetSavedMessages()
		if err != nil {
			fmt.Println("Error getting saved messages:", err)
			thisPage.AssignMenuSelection("back", "Back", func(key string) (int, error) {
				return WINDOW_BACK, nil
			})
		}

		if messages, ok := messagePacket["messages"]; ok {
			if messagesSlice, ok := messages.([]interface{}); ok {
				for idx, message := range messagesSlice {
					if messageStr, ok := message.(string); ok {
						thisPage.AssignMenuSelection(fmt.Sprint(idx), messageStr, func(key string) (int, error) {
							idx, err := strconv.Atoi(key)
							if err != nil {
								return 0, err
							}

							_, err = DeleteSavedMessage(idx)

							return 0, err
						})
					}
				}
			}
		}

		thisPage.AssignMenuSelection("back", "Back", func(key string) (int, error) {
			return WINDOW_BACK, nil
		})
	}

	return newPage
}
