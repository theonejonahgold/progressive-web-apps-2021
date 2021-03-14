package ssr

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars"
	m "github.com/theonejonahgold/pwa/models"
	u "github.com/theonejonahgold/pwa/utils"
)

func SSR() (context.Context, error) {
	ctx := context.Background()
	runSnowpackDevBuilds(ctx)

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

	return ctx, app.Listen(":3000")
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
	}, "layouts/main")
}

func story(c *fiber.Ctx) error {
	id := c.Params("id")
	j, err := u.Fetch("https://hacker-news.firebaseio.com/v0/item/" + id + ".json")
	if err != nil {
		return err
	}

	story, err := u.ParseStory(j)
	if err != nil {
		return err
	}

	if story.Type == "comment" {
		c.Status(404).SendString("Kon die stoorie niet vinden. Probeer eens een andere!")
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(1)
	u.FetchComments(story, &wg, 2)
	wg.Wait()

	return c.Render("story", fiber.Map{
		"story": story,
	}, "layouts/main")
}

func notFound(c *fiber.Ctx) error {
	return c.Status(404).SendString("Kon die stoorie niet vinden. Probeer eens een andere!")
}
