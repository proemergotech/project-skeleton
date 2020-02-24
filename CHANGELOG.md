# Release Notes

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
