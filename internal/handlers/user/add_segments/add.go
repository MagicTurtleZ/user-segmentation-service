package addsegments

import (
	"time"
	"woonbeaj/segments/internal/halpers/response"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	UserId 	int 		`json:"userId"`
	AddSlug []string 	`json:"addSlug,omitempty"`
	DelSlug []string 	`json:"deleteSlug,omitempty"`
	TTL		string		`json:"ttl,omitempty"`
}

type SegmentAdder interface {
	AddSegment(userId int, segName string, ttl *time.Time) error
	RemoveSegment(userId int, segName string) error
}

func (r *Request) GetTimeStamp() (*time.Time, error) {
	curDate := time.Now().Unix()
	ttlDur, err := time.ParseDuration(r.TTL)
	if err != nil {
		return nil, err
	} 
	newDate := time.Unix(curDate, 0).Add(ttlDur)

	return &newDate, nil
}

func AddSegmentsToUser(db SegmentAdder) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(Request)
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.Error(err.Error()))
		}

		for _, val := range req.DelSlug {
			err = db.RemoveSegment(req.UserId, val)
			if err != nil {
				return c.Status(fiber.StatusGatewayTimeout).JSON(response.Error(err.Error()))
			} 
		}

		var ttl *time.Time 
		if len(req.TTL) > 0 {
			ttl, err = req.GetTimeStamp()
			if err != nil {
				return c.Status(fiber.StatusGatewayTimeout).JSON(response.Error(err.Error()))
			}
		}

		for _, val := range req.AddSlug {
			err = db.AddSegment(req.UserId, val, ttl)
			if err != nil {
				return c.Status(fiber.StatusGatewayTimeout).JSON(response.Error(err.Error()))
			} 
		}
		
		return c.Status(fiber.StatusOK).JSON(response.OK())
	}
}