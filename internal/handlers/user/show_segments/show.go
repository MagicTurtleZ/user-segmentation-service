package showsegments

import (
	"woonbeaj/segments/internal/halpers/response"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	UserId int 			`json:"userId"`
}

type Response struct {
	response.Response
	UserRule []string	`json:"userSlug,omitempty"`
}

type SegmentsGeter interface {
	GetAllSegments(userId int) ([]string, error)
}

func AddSegmentsToUser(db SegmentsGeter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(Request)
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
		}

		res, err := db.GetAllSegments(req.UserId)
		if err != nil {
			return c.Status(fiber.StatusGatewayTimeout).JSON(response.Error(err.Error()))
		}

		resp := Response{
			Response: response.OK(),
			UserRule: res,
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}
}