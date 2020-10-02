{{ if .PublicRest -}}
FORMAT: 1A

# Dliver Project Skeleton

GENERATED, DO NOT EDIT, to regenerate:
- modify [definitions file](apimd/main.go)
- run `go run apimd/main.go`

## Group Http

### Base API [/]

#### Healthcheck [GET /healthcheck]
Healthcheck route, used for liveness probe.

+ Response 200

        ok

#### Metrics [GET /metrics]
Metrics route, returns useful information about the service.

+ Response 200

        metrics

### Public Endpoints [/api/v1]

#### Dummy endpoint [POST /api/v1/dummy/{dummy_param_1}]
Dummy endpoint's description

+ Parameters
    + `dummy_param_1`: `dummy_p1` (string)

+ Request
    + Attributes
        + `dummy_data_1`: `dummy_d1` (string) - required parameter of the dummy endpoint
        + `dummy_data_2`: `dummy_d2` (string, optional)

+ Response 200

+ Response 400
    + Attributes
        + `error`
            + `code`: `ERR_VALIDATION` (string)
            + `details`(array)
                + (object)
                    + `field`: `dummy_data_1` (string)
                    + `validator`: `required` (string)

+ Response 500
    + Attributes
        + `error`
            + `code`: `ERR_INTERNAL` (string)
{{ end -}}
