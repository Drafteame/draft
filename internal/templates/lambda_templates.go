package templates

import "github.com/Drafteame/draft/internal/dtos"

type LambdaTemplates struct {
	Cron   LambdaCron
	HTTP   LambdaHTTP
	Plain  LambdaPlain
	SnsSqs LambdaSnsSqs
	Sqs    LambdaSqs
}

func (l *LambdaTemplates) SetSqs(sqs LambdaSqs) {
	l.Sqs = sqs
}

func (l *LambdaTemplates) SetSnsSqs(sqs LambdaSnsSqs) {
	l.SnsSqs = sqs
}

func (l *LambdaTemplates) SetPlain(plain LambdaPlain) {
	l.Plain = plain
}

func (l *LambdaTemplates) SetHTTP(http LambdaHTTP) {
	l.HTTP = http
}

func (l *LambdaTemplates) SetCron(cron LambdaCron) {
	l.Cron = cron
}

func NewLambdaTemplates(data dtos.LambdaInput) (*LambdaTemplates, error) {
	l := new(LambdaTemplates)

	if err := loadLambdaCron(l, data); err != nil {
		return nil, err
	}

	if err := loadLambdaHTTP(l, data); err != nil {
		return nil, err
	}

	if err := loadLambdaPlain(l, data); err != nil {
		return nil, err
	}

	if err := loadLambdaSnsSqs(l, data); err != nil {
		return nil, err
	}

	if err := loadLambdaSqs(l, data); err != nil {
		return nil, err
	}

	return l, nil
}
