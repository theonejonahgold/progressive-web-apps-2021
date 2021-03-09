package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars"
	m "github.com/theonejonahgold/pwa/models"
	u "github.com/theonejonahgold/pwa/utils"
	s "github.com/theonejonahgold/pwa/utils/sort"
)

func main() {
	engine := handlebars.New("./views", ".hbs")
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("", "./static", fiber.Static{
		Compress: true,
	})
	app.Get("/", index)
	app.Get("/story/:id", story)
	app.Get("*", notFound)

	log.Fatal(app.Listen(":3000"))
}

func index(c *fiber.Ctx) error {
	stories, err := u.GetTopStories()
	if err != nil {
		fmt.Println(err)
		return err
	}
	sort.Sort(s.ByScore(*stories))

	return c.Render("pages/index", fiber.Map{
		"stories": stories,
		"title":   "hoom",
	}, "layouts/main")
}

func story(c *fiber.Ctx) error {
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

	var s m.Story
	err = json.Unmarshal(j, &s)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if s.Type == "comment" {
		c.Status(404).SendString("Kon die stoorie niet vinden. Probeer eens een andere!")
		return nil
	}

	comments, err := u.RetrieveCommentsForStory(&s)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.Render("pages/story", fiber.Map{
		"story":    s,
		"comments": comments,
		"title":    s.Title,
	}, "layouts/main")
}

func notFound(c *fiber.Ctx) error {
	return c.Status(404).SendString("Kon die stoorie niet vinden. Probeer eens een andere!")
}
