package create

import (
	"woonbeaj/segments/internal/halpers/response"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	Slug string `json:"slug"`
}

type SegmentCreater interface {
	CreateSegment(segName string) error
}

func NewSegment(db SegmentCreater) fiber.Handler { 
	return func(c *fiber.Ctx) error {
		req := new(Request)
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
		}

		err = db.CreateSegment(req.Slug)
		if err != nil {
			return c.Status(fiber.StatusGatewayTimeout).JSON(response.Error(err.Error()))
		}

		return c.Status(fiber.StatusOK).JSON(response.OK())
	}
}