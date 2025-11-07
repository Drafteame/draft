package createproject

import (
	"github.com/Drafteame/draft/internal/dtos"
)

func GetForm(input *dtos.CreateProjectInput) error {
	return baseForm(input)
}
