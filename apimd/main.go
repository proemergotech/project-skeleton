package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.com/proemergotech/apimd-generator-go/generator"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	microtime "gitlab.com/proemergotech/microtime-go"
	uuid "gitlab.com/proemergotech/uuid-go"
)

func main() {
	g := generator.NewGenerator()
	d := &definitions{}
	g.Generate(d)
}

type definitions struct {
	factory *generator.Factory
}

type value struct {
	*generator.Value
}

func (d *definitions) Name() string {
	return strings.Title(strings.Replace(config.AppName, "-", " ", -1))
}

func (d *definitions) OutputPath() string {
	return "./API.md"
}

func (d *definitions) Usage() []string {
	return []string{
		"modify [definitions file](apimd/main.go)",
		"run `go run apimd/main.go`",
	}
}

func (d *definitions) Groups(factory *generator.Factory) []generator.Group {
	d.factory = factory
	defer func() {
		d.factory = nil
	}()

	return []generator.Group{
		&generator.HTTPGroup{
			Name:        "Health API",
			RoutePrefix: "",
			Routes: []*generator.HTTPRoute{
				{
					Name: "Healthcheck",
					Description: []string{
						"Healthcheck route, used for liveness probe.",
					},
					Path:   "/healthcheck",
					Method: http.MethodGet,
					Responses: map[int]interface{}{
						http.StatusOK: d.body("ok").String(),
					},
				},
			},
		},
	}
}

func (*definitions) ParseIndex(index interface{}) (int, error) {
	switch ind := index.(type) {
	case float64:
		return int(ind), nil

	case string:
		t := &microtime.Time{}
		if err := t.UnmarshalJSON([]byte("\"" + ind + "\"")); err == nil {
			return int(t.Unix()), nil
		}

		indInt, err := strconv.Atoi(ind)
		if err != nil {
			// use as-is
			return 0, nil
		}
		return indInt, nil

	default:
		return 0, nil
	}
}

func (d *definitions) param(val string) *value {
	return &value{Value: d.factory.Param(val)}
}

func (d *definitions) query(val string) *value {
	return &value{Value: d.factory.Query(val)}
}

func (d *definitions) body(val string) *value {
	return &value{Value: d.factory.Body(val)}
}

func (v *value) desc(d string) *value {
	v.Description(d)
	return v
}

func (v *value) opt() *value {
	v.Optional()
	return v
}

func (v *value) microtime() microtime.Time {
	return microtime.Time{Time: time.Unix(v.Int64(), 0)}
}

func (v *value) uuid() uuid.UUID {
	return uuid.UUID(v.String())
}

func (d *definitions) httpError(err error, details ...map[string]interface{}) schema.HTTPError {
	return schema.HTTPError{
		Error: schema.Error{
			Code:    d.body(schema.ErrorCode(err)).String(),
			Message: d.body(err.Error()).String(),
			Details: details,
		},
	}
}
