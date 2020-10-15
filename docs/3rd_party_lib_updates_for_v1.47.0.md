# Go lib updates

## JSON iterator from v1.1.7 to v1.1.10

### Notable changes
- Nil maps are now correctly encoded into 'null' values
- Removed quotation check for key when decoding map
- Limited nesting depth and added config option for it

### Quotation check removal for key when decoding map
It does not need to be checked whether the key is surrounded by quotation.
The key might not be a string if an extension is registered to customize the
map key encoder/decoder. It may be an integer, float, or even a struct.

## Validator from v10.0.1 to v10.4.0

### Notable changes
- country_code validation
- required_if validation
- required_unless validation

### Country code validation

#### ISO3166-1 alpha-2
This validates that a string value is a valid country code based on iso3166-1 alpha-2 standard.

see: https://www.iso.org/iso-3166-country-codes.html

```
iso3166_1_alpha2
```

#### ISO3166-1 alpha-3
This validates that a string value is a valid country code based on ISO3166-1 alpha-3 standard.

see: https://www.iso.org/iso-3166-country-codes.html
```
iso3166_1_alpha3
```

#### ISO3166-1 alphanumeric
This validates that a string value is a valid country code based on ISO3166-1 alphanumeric standard.

see: https://www.iso.org/iso-3166-country-codes.html

```
iso3166_1_alpha_numeric
```

### Alias
alias is "iso3166_1_alpha2|iso3166_1_alpha3|iso3166_1_alpha_numeric"

```
country_code
```

### Required if validation
The field under validation must be present and not empty only if all
the other specified fields are equal to the value following the specified
field. For strings ensures value is not "". For slices, maps, pointers,
interfaces, channels and functions ensures the value is not nil.

#### Examples
```
// require the field if Field1 is equal to the parameter given
required_if=Field1 foobar

// require the field if Field1 and Field2 is equal to the values respectively
required_if=Field1 foo Field2 bar
```

### Required unless validation
The field under validation must be present and not empty unless all
the other specified fields are equal to the value following the specified
field. For strings ensures value is not "". For slices, maps, pointers,
interfaces, channels and functions ensures the value is not nil.

#### Examples
```
// require the field unless Field1 is equal to the parameter given
required_unless=Field1 foobar

// require the field unless Field1 and Field2 is equal to the values respectively
required_unless=Field1 foo Field2 bar
```

## Echo from v4.1.11 to v4.1.17

### Notable changes
- HTTP/2 Cleartext mode (H2C)
- Proxy middleware can now modify response
- gzipResponseWriter now implements http.Pusher interface

### HTTP/2 Cleartext mode
It is essentially HTTP/2 but without TLS. The standard golang code supports
HTTP2 but does not directly support H2C. H2C support only existed in the
golang.org/x/net/http2/h2c package until now.
May be beneficial during Yafuds rewrite to REST instead of GRPC.

### gzipResponseWriter http.Pusher interface implementation
Pusher is the interface implemented by ResponseWriters that support HTTP/2
 server push. see https://tools.ietf.org/html/rfc7540#section-8.2.

## Elastic from v6.2.25 to v6.2.35

### Notable changes
- Added support to other_bucket for Filters aggregation
- Added support for span queries

### Other bucket for Filters aggregation
The other_bucket parameter can be set to add a bucket to the response which
will contain all documents that do not match any of the given filters.

see: https://www.elastic.co/guide/en/elasticsearch/reference/6.8/search-aggregations-bucket-filters-aggregation.html#other-bucket

### Span queries
Span queries are low-level positional queries which provide expert control over
the order and proximity of the specified terms. These are typically used to
implement very specific queries on legal documents or patents.

see: https://www.elastic.co/guide/en/elasticsearch/reference/6.8/span-queries.html

## Opentracing from v1.1.0 to v1.2.0

### Notable changes
- Added an extension to Tracer interface for custom go context creation
- Added log/fields helpers for keys from specification
- Go modules support

### Tracer extension interface
The new TracerContextWithSpanExtension interface provides a way to
hook into the ContextWithSpan function, so the implementation can put
some extra information to the context.

The opentracing to opentelemetry bridge needs the context to set the
current opentelemetry span, so the opentelemetry API in the layer
below the one using opentracing can still get the right parent span.

### Log/fields helpers for keys
Opentracing spec specifies preferred names for log messages [1]
Added declaration for fields which are meaningful for golang

- log: add helpers function for spec keys
- ext: add LogError() helper for error

LogError() helper can not be declarated in log package because it depends
opentrace.Span which result in cyclic depencency

Footnotes:
[1] https://github.com/opentracing/specification/blob/master/semantic_conventions.md#log-fields-table

## Prometheus client from v1.1.0 to v.1.7.1

### Notable changes
- Add exemplars to counter and histogram
- Added promlint

### Exemplars
Exemplars allow linking certain metrics to example traces:

```
# TYPE foo histogram
foo_bucket{le="0.01"} 0
foo_bucket{le="0.1"} 8 # {id="abc"} 0.043
foo_bucket{le="1"} 10 # {id="def"} 0.29
foo_bucket{le="10"} 17 # {id="ghi"} 7.73
foo_bucket{le="+Inf"} 18
foo_count 18
foo_sum 324789.3
foo_created 1520430000.123
```

### Promlint
Provides a linter for Prometheus metrics.

## Cobra from v0.0.5 to v1.0.0

### Notable changes
- Added support for context.Context to commands

### Context adding to commands

```
ctx := context.TODO()

rootCmd := &Command{Use: "root", Run: ctxRun, PreRun: ctxRun}

_, err := executeCommandWithContext(ctx, rootCmd, "")
```

## Viper from v1.4.0 to v1.7.1

### Notable changes
- Added support for dot env files
- Added support for int slice flags
- Implemented ability to unmarshal keys containing dots to structs
- Added support for ini files with sections
- Added string replacer interface and env key replacer option
- Added support for config files without extensions

## Jaeger client from v2.20.1 to v.2.25.0

### Notable changes
- Added support for custom HTTP headers when reporting spans over HTTP

## Zap from v1.10.0 to v1.16.0

### Notable changes
- Switched to go modules
- Function Log Entry

### Function log entry
Introduces function name in log entry.

see: https://github.com/uber-go/zap/issues/716