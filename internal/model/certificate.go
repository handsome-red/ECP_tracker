package model

import (
	"fmt"
	"time"
)

const format = "2006-01-02"

type Certificate struct {
	SecondName string
	FirstName  string
	Patronymic string
	District   string
	Expiration time.Time
}

func (c *Certificate) FullName() string {
	return fmt.Sprintf("%s%s%s", c.SecondName, c.FirstName, c.Patronymic)
}

func (c *Certificate) Expired() string {
	return c.Expiration.Format(format)
}
