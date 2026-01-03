package model

import (
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

type Time time.Time

func (t Time) MarshalGQL(w io.Writer) {
	graphql.MarshalTime(time.Time(t)).MarshalGQL(w)
}

func (t *Time) UnmarshalGQL(v interface{}) error {
	parsed, err := graphql.UnmarshalTime(v)
	if err != nil {
		return err
	}
	*t = Time(parsed)
	return nil
}
