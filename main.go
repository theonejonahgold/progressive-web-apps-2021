package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars"
)

type story struct {
	ID          int    `json:"id"`
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	Deleted     bool   `json:"deleted"`
}

type storyIDArray []uint8

func main() {
	fmt.Println("Hi there, initialising app!")

	engine := handlebars.New("./views", ".hbs")
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", index)

	app.Listen(":3000")
}

func index(c *fiber.Ctx) error {
	stories, err := getTopStories()
	if err != nil {
		fmt.Println(err)
		return c.SendString("Henkernieuws is poepie")
	}

	ctx := fiber.Map{
		"stories": stories,
	}

	return c.Render("index", ctx)
}

func getTopStories() (*[]story, error) {
	res, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var storyIDs storyIDArray
	storyIDs, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var stories []story
	for i := 100; i > 90; i-- {
		res, err := http.Get("https://hacker-news.firebaseio.com/v0/item/" + string(storyIDs[i]) + ".json")
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		bytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var s story
		err = json.Unmarshal(bytes, &s)
		if err != nil {
			return nil, err
		}

		stories = append(stories, s)
	}
	return &stories, nil
}
