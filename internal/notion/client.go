package notion

import (
	"bytes"
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

func getTodos(c *Config) (model.BlockResponse, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=1", c.ParentPageID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return model.BlockResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.NotionIntegrationSecret)
	req.Header.Set("Notion-Version", "2025-09-03")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return model.BlockResponse{}, err
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var monthBlockRes, dayBlockRes, todos model.BlockResponse
	if err := json.Unmarshal(body, &monthBlockRes); err != nil {
		return model.BlockResponse{}, err
	}

	url = fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=1", monthBlockRes.Results[0].Id)

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return model.BlockResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.NotionIntegrationSecret)
	req.Header.Set("Notion-Version", "2025-09-03")

	resp, err = client.Do(req)
	if err != nil {
		return model.BlockResponse{}, err
	}

	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()

	if err := json.Unmarshal(body, &dayBlockRes); err != nil {
		return model.BlockResponse{}, err
	}

	url = fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", dayBlockRes.Results[0].Id)

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return model.BlockResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.NotionIntegrationSecret)
	req.Header.Set("Notion-Version", "2025-09-03")

	resp, err = client.Do(req)
	if err != nil {
		return model.BlockResponse{}, err
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &todos); err != nil {
		return model.BlockResponse{}, err
	}

	return todos, nil
}

func (c *Config) GetTodos() error {
	todos, err := getTodos(c)
	if err != nil {
		return err
	}

	for i, block := range todos.Results {
		checkbox := "[ ]"
		if block.RichText.Checked {
			checkbox = "[x]"
		}
		fmt.Print(i, " - ")
		fmt.Print(block.RichText.Todo[0].Title + " ")
		fmt.Println(checkbox)
	}

	return nil
}

func markTodo(c *Config, id int, checked bool) error {
	todos, err := getTodos(c)
	if err != nil {
		return err
	}

	bodyMap := map[string]interface{}{
		"to_do": map[string]bool{
			"checked": checked,
		},
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		return err
	}

	var blockID string
	for i, block := range todos.Results {
		if i == id {
			blockID = block.Id
			break
		}
	}

	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s", blockID)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.NotionIntegrationSecret)
	req.Header.Set("Notion-Version", "2025-09-03")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	return err
}

func (c *Config) MarkTodo(id int) error {
	return markTodo(c, id, true)
}

func (c *Config) UnmarkTodo(id int) error {
	return markTodo(c, id, false)
}
