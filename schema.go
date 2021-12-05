package main

import "time"

//+migu
type Users struct {
	Id        int64     `migu:"pk,autoincrement" db:"id"`
	Uniq      string    `migu:"unique" db:"uniq"`
	CreatedAt time.Time `migu:"extra:DEFAULT CURRENT_TIMESTAMP" db:"created_at"`
	UpdatedAt time.Time `migu:"extra:DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" db:"updated_at"`
}
