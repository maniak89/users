//go:generate reform
package storage

import (
	"time"
)

//reform:users
type User struct {
	ID        string     `reform:"id,pk"`
	Login     string     `reform:"login"`
	Password  string     `reform:"password"`
	CreatedAt time.Time  `reform:"created_at"`
	UpdatedAt time.Time  `reform:"updated_at"`
	LastLogin *time.Time `reform:"last_login"`
}

func (s *User) BeforeUpdate() error {
	s.UpdatedAt = time.Now()
	return nil
}
