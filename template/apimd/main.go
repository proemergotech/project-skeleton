//%: {{ if .PublicRest }}
package main

import (
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/proemergotech/apimd-generator-go/generator"
	"gitlab.com/proemergotech/microtime-go"

	//%:{{ `
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/di"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/skeleton"
	//%: ` | replace "dliver-project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
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
			Name:        "Base API",
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
				{
					Name: "Metrics",
					Description: []string{
						"Metrics route, returns useful information about the service.",
					},
					Path:   "/metrics",
					Method: http.MethodGet,
					Responses: map[int]interface{}{
						http.StatusOK: d.body("metrics").String(),
					},
				},
			},
		},
		&generator.HTTPGroup{
			Name:        "Public Endpoints",
			RoutePrefix: "/api/v1",
			Routes: []*generator.HTTPRoute{
				{
					Name: "Dummy endpoint",
					Description: []string{
						"Dummy endpoint's description",
					},
					Path:   "/dummy/:dummy_param_1",
					Method: http.MethodPost,
					Request: &struct {
						DummyParam1 string `param:"dummy_param_1" validate:"required"`
						DummyData1  string `json:"dummy_data_1" validate:"required"`
						DummyData2  string `json:"dummy_data_2"`
					}{
						DummyParam1: d.param("dummy_p1").String(),
						DummyData1:  d.body("dummy_d1").desc("required parameter of the dummy endpoint").String(),
						DummyData2:  d.body("dummy_d2").opt().String(),
					},
					Responses: map[int]interface{}{
						http.StatusOK: nil,
						http.StatusBadRequest: d.validationError(&struct {
							DummyData1 string `json:"dummy_data_1" validate:"required"`
							DummyData2 string `json:"dummy_data_2"`
						}{}),
						//%:{{ `
						http.StatusInternalServerError: d.publicHTTPError(skeleton.SemanticError{Msg: "Caused by internal problem, should be solved on server side"}.E()),
						//%: ` | replace "skeleton" .SchemaPackage }}
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
		err := t.UnmarshalJSON([]byte("\"" + ind + "\""))
		if err == nil {
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

func (d *definitions) validationError(i interface{}) schema.HTTPError {
	v, err := di.NewValidator()
	if err != nil {
		panic(err)
	}
	return d.publicHTTPError(v.Validate(i))
}
func (d *definitions) publicHTTPError(err error) schema.HTTPError {
	res, _ := schema.ToPublicHTTPError(err)
	return *res
}

//%: {{ end }}
