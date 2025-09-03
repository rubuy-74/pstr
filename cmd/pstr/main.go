package main

import (
	"fmt"
	"os"

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

	app.Post("/check", func(c *fiber.Ctx) error {
		regexRequest := new(RegexRequest)
		if err := c.BodyParser(regexRequest); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		regex := regexRequest.Regex
		matchString := regexRequest.MatchString
		fmt.Println(regex)
		fmt.Println(matchString)

		parsedRegex, err := parser.Parse(regex)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "failed to parse regex"})
		}

		nfa := state_machine.ToNFA(parsedRegex)
		valid := nfa.Check(matchString, -1)
		return c.JSON(fiber.Map{
			"valid": valid,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(":" + port)
}
