package main

import "time"

//+migu
type Users struct {
	Id        int64     `migu:"pk,autoincrement"`
	Uniq      string    `migu:"unique"`
	CreatedAt time.Time `migu:"extra:DEFAULT CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `migu:"extra:DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}
