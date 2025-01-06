package deleteproject

import (
	"github.com/Drafteame/draft/internal/dtos"
)

func GetForm(input *dtos.DeleteProjectInput) error {
	return baseForm(input)
}
