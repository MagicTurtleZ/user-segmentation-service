package audit

import (
	csvgen "woonbeaj/segments/internal/CSV"
	"woonbeaj/segments/internal/halpers/response"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	Year 	int 	`json:"year"`
	Month 	int		`json:"month"`
}

func ReportCSV(ag csvgen.AuditGetter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(Request)
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
		}

		report, err := csvgen.GenReport(req.Month, req.Year, ag)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
		}

		c.Set("Content-Disposition", "attachment; filename=report.csv")
		c.Set("Content-Type", "text/csv")

		err = c.SendFile(report, false)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
		}

		return nil
	}
}