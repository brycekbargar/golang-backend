package articledomain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.CustomTypeTagMap.Set("positive",
		govalidator.CustomTypeValidator(func(i interface{}, o interface{}) bool {
			switch num := i.(type) {
			case int:
				return num >= 0
			default:
				return false
			}
		}))
}

// Comment is an individual comment associated with a single Article.
type Comment struct {
	ID           int    `valid:"positive"`
	Body         string `valid:"required"`
	CreatedAtUTC time.Time
	AuthorEmail  string `valid:"required,email"`
}

// NewComment creates a new comment with the provided information and defaults for the rest
func NewComment(body string, author string) (*Comment, error) {
	return (&Comment{
		Body:        body,
		AuthorEmail: author,
	}).Validate()
}

// Validate returns the provided Article if it is valid, otherwise error will contain validation errors.
func (c *Comment) Validate() (*Comment, error) {
	if v, err := govalidator.ValidateStruct(c); !v {
		return nil, err
	}

	return c, nil
}
