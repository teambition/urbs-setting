# Change Log

All notable changes to this project will be documented in this file starting from version **v1.0.0**.
This project adheres to [Semantic Versioning](http://semver.org/).

-----
## [1.8.0] - 2020-09-16

**Change:**
- Support parent-child relationship label.

## [1.7.1] - 2020-08-13

**Change:**
- Don't apply label rule in a product when someone exists.

## [1.7.0] - 2020-07-08

**Change:**
- Support primary and secondary mysql connections.

## [1.6.0] - 2020-07-07

**Change:**
- `GET /users/:uid/labels:cache` and `GET /v1/users/:uid/settings:unionAll` support anonymous user, uid should with prefix `anon-`.
- `GET /v1/users/:uid/settings` and `GET /v1/groups/:uid/settings` support `channel` and `client` query.

## [1.5.0] - 2020-06-08

**Change:**

- Add API `DELETE /v1/products/{product}/modules/{module}/settings/{setting}:cleanup` that cleanup all rules, users and groups on the setting.
- Add API `DELETE /v1/products/{product}/labels/{label}:cleanup` that cleanup all rules, users and groups on the label.

## [1.4.0] - 2020-05-27

**Change:**

- Change user's setting and label API.
- Change group's setting and label API.
- Use [goqu](github.com/doug-martin/goqu/v9) instead of gorm.
- Support more query parameters for settings API.

## [1.3.3] - 2020-05-18

**Fixed:**

- Fix API's totalSize count.
- Fix tracing middleware.

## [1.3.2] - 2020-05-13

**Change:**

- Create setting with more params.

**Fixed:**

- Fix `name` field for `urbs_statistic` table and `urbs_lock` table.

## [1.3.1] - 2020-05-11

**Change:**

- Create label with more params.

## [1.3.0] - 2020-05-08

**Change:**

- Support label rule and setting rule.
- Support search for list APIs.
- Change APIs to camelCase, see https://github.com/json-api/json-api/issues/1255.

## [1.2.3] - 2020-04-03

**Change:**

- Improve cached labels API.
- Add module, setting, label documents.

## [1.2.2] - 2020-04-01

**Change:**

- Improve swagger document.
- Add user and group documents.

**Fixed:**

- Fix settings API.

## [1.2.1] - 2020-03-29

**Change:**

- Update Gear version.

## [1.2.0] - 2020-03-25

**Change:**

- Add test cases for all APIs.

## [1.1.2] - 2020-03-19

**Fixed:**

- API should not response `id` field.
- Fixed request body template.

## [1.1.0] - 2020-03-14

**Changed:**

- Support kind for group.
- Support pagination for List API.
- Improve SQL schemas.
- Improve code.
