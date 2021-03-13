package ssr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars"
	m "github.com/theonejonahgold/pwa/models"
	u "github.com/theonejonahgold/pwa/utils"
)

func SSR() error {
	engine := handlebars.New("./views", ".hbs")
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./dist/static", fiber.Static{
		Compress: true,
	})
	app.Get("/", index)
	app.Get("/story/:id", story)
	app.Get("*", notFound)

	return app.Listen(":3000")
}

func index(c *fiber.Ctx) error {
	stories, err := u.GetTopStories()
	if err != nil {
		fmt.Println(err)
		return err
	}
	sort.Sort(m.StoriesByScore(stories))

	return c.Render("index", fiber.Map{
		"stories": stories,
	})
}

func story(c *fiber.Ctx) error {
	fmt.Println("hoi")
	id := c.Params("id")
	res, err := http.Get("https://hacker-news.firebaseio.com/v0/item/" + id + ".json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	j, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	story := m.NewStory()
	err = json.Unmarshal(j, story)
	if err != nil {
		fmt.Println(story)
		return err
	}

	if story.Type == "comment" {
		c.Status(404).SendString("Kon die stoorie niet vinden. Probeer eens een andere!")
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(1)
	u.FetchComments(story, &wg)
	wg.Wait()

	return c.Render("story", fiber.Map{
		"story": story,
	})
}

func notFound(c *fiber.Ctx) error {
	return c.Status(404).SendString("Kon die stoorie niet vinden. Probeer eens een andere!")
}
