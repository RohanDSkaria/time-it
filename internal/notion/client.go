package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/RohanDSkaria/time-it/internal/model"
)

type cache struct {
	MonthKey    string `json:"month_key"`
	MonthPageID string `json:"month_page_id"`
	DayKey      string `json:"day_key"`
	DayPageID   string `json:"day_page_id"`
}

func cachePath() string {
	home, _ := os.UserHomeDir()
	return home + "/.config/time-it/cache.json"
}

func getPageID(format string, isMonth bool) string {
	key := time.Now().Format(format)

	data, err := os.ReadFile(cachePath())
	if err != nil {
		return ""
	}

	var c cache
	if err := json.Unmarshal(data, &c); err != nil {
		return ""
	}

	if isMonth {
		if c.MonthKey == key {
			return c.MonthPageID
		}
	} else {
		if c.DayKey == key {
			return c.DayPageID
		}
	}

	return ""
}

func saveMonthPageID(pageID string) {
	var c cache

	data, err := os.ReadFile(cachePath())
	if err != nil || len(data) == 0 {
		c = cache{
			MonthKey:    time.Now().Format("2006-01"),
			MonthPageID: pageID,
		}
	} else {
		json.Unmarshal(data, &c)

		if c.MonthPageID == pageID {
			return
		}

		c.MonthKey = time.Now().Format("2006-01")
		c.MonthPageID = pageID
	}

	data, err = json.Marshal(c)
	if err != nil {
		return
	}

	os.WriteFile(cachePath(), data, 0644)
}

func saveDayPageID(pageID string) {
	var c cache

	data, _ := os.ReadFile(cachePath())

	json.Unmarshal(data, &c)

	if c.DayPageID == pageID {
		return
	}

	c.DayKey = time.Now().Format("2006-01-02")
	c.DayPageID = pageID

	data, err := json.Marshal(c)
	if err != nil {
		return
	}

	os.WriteFile(cachePath(), data, 0644)
}

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

func getBlockID(parentPageID, notionIntegrationSecret string) (string, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=1", parentPageID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+notionIntegrationSecret)
	req.Header.Set("Notion-Version", "2025-09-03")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var blockRes model.BlockResponse
	if err := json.Unmarshal(body, &blockRes); err != nil {
		return "", err
	}

	return blockRes.Results[0].Id, nil
}

func getTodos(c *Config) (model.BlockResponse, error) {
	dayPageID := getPageID("2006-01-02", false)

	if dayPageID == "" {
		monthPageID := getPageID("2006-01", true)

		if monthPageID == "" {
			id, err := getBlockID(c.ParentPageID, c.NotionIntegrationSecret)
			if err != nil {
				return model.BlockResponse{}, err
			}

			monthPageID = id
			saveMonthPageID(monthPageID)
		}

		id, err := getBlockID(monthPageID, c.NotionIntegrationSecret)
		if err != nil {
			return model.BlockResponse{}, err
		}

		dayPageID = id
		saveDayPageID(dayPageID)
	}

	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children", dayPageID)

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
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var todos model.BlockResponse
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
