package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rubuy-74/pstr/internal/parser"
	"github.com/rubuy-74/pstr/internal/state_machine"
)

type RegexRequest struct {
	Regex       string `json:"regex"`
	MatchString string `json:"string"`
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Post("/check", func(c *fiber.Ctx) error {
		regexRequest := new(RegexRequest)
		if err := c.BodyParser(regexRequest); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		regex := regexRequest.Regex
		matchString := regexRequest.MatchString

		parsedRegex, err := parser.Parse(regex)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "failed to parse regex",
				"message": err,
			})
		}

		nfa, err := state_machine.ToNFA(parsedRegex)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   "failed to create NFA",
				"message": err,
			})
		}
		valid := nfa.Check(matchString, -1)
		return c.JSON(fiber.Map{
			"valid": valid,
		})
	})

	app.Listen(":3000")
}
