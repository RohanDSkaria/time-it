package notion

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/RohanDSkaria/time-it/internal/model"
)

type Config struct {
	NotionIntegrationSecret string `json:"notion_integration_secret"`
	ParentPageID            string `json:"parent_page_id"`
}

func New() *Config {
	home, _ := os.UserHomeDir()
	path := home + "/.config/time-it/config.json"
	
	data, _ := os.ReadFile(path)
	var cfg Config
	json.Unmarshal(data, &cfg)

	return &cfg
}

func (c *Config) GetTodos() error {
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=1", c.ParentPageID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer " + c.NotionIntegrationSecret)
	req.Header.Set("Notion-Version", "2025-09-03")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var monthBlockRes, dayBlockRes, todos model.BlockResponse
	if err := json.Unmarshal(body, &monthBlockRes); err != nil {
		return err
	}

	url = fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=1", monthBlockRes.Results[0].Id)

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "Bearer " + c.NotionIntegrationSecret)
	req.Header.Set("Notion-Version", "2025-09-03")

	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()

	if err := json.Unmarshal(body, &dayBlockRes); err != nil {
		return err
	}

	url = fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", dayBlockRes.Results[0].Id)

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "Bearer " + c.NotionIntegrationSecret)
	req.Header.Set("Notion-Version", "2025-09-03")

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	body, _ = io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &todos); err != nil {
		return err
	}

	for _, block := range todos.Results {
		checkbox := "[ ]"
		if block.RichText.Checked {
			checkbox = "[x]"
		}
		fmt.Print(checkbox + " ")
		fmt.Println(block.RichText.Todo[0].Title)
	}

	return nil
}
