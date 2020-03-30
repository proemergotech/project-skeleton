# Release Notes

## v1.18.0 / 2020-03-30
- create public details for validation.Error

## v1.17.0 / 2020-03-26
- add http method to routeNotFoundError and methodNotAllowedError

## v1.16.0 / 2020-03-23
- update log-go to v3.0.1 (truncate long zaplog messages)

## v1.15.0 / 2020-03-23
- create lib from errors
- update yafuds-client-go to v1.2.0

## v1.14.0 / 2020-03-20
- remove .env.example (remove config example too if it exists)
- move initConfig to command Run method so we can create separate config for every command

## v1.13.0 / 2020-03-18
- update retry middleware (not backward compatible)

## v1.12.0 / 2020-03-05
- remove API.md reference from README.md
- update log-go to v3.0.0

## v1.11.0 / 2020-03-02
- remove apimd generator
- update to go 1.14
- change environment url to `https://camplace.dev` in .gitlab-ci.yml

## v1.10.0 / 2020-02-27
- refactor clientHTTPError to return better error messages

## v1.9.0 / 2020-02-24
- added centrifugeJSON so centrifuge messages are marshalled using the `centrifuge` tag

## v1.8.0 / 2020-02-14
- fix gebQueue.Start() usage in event controller
- move stage to the first place in gitlab ci definitions
- update README.md to mention go mod instead of dep

## v1.7.0 / 2020-02-13
- update geb/trace/log/yafuds-client library

## v1.6.0 / 2020-02-06
- added -mod-readonly to .gitlab-ci.yml test job
- updated centrifuge-client-go to v2.2.2

## v1.5.0 / 2020-01-31
- removed -mod-readonly from build.sh

## v1.4.0 / 2020-01-27
- added notblank validation rule

## v1.3.1 / 2020-01-20
- use built-in echo methods

## v1.3.0 / 2020-01-07
- update proemergotech libraries to v1.0.0
- update verify job in ci script to run only on non master branches and ignore everything else (tag pipelines etc)

## v1.2.0 / 2019-12-18
- update dependencies
- fix verify.sh
- remove `newHTTPServer` from container

## v1.1.0 / 2019-12-13
- use echo directly in rest/server

## v1.0.0 / 2019-12-13
- update to echo v4
- add CHANGELOG.md update checking for pipeline
- restructure gitlab-ci.yml to run test/build and verify parallel
