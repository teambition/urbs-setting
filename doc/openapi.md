---
title: urbs-setting
language_tabs:
  - shell: Shell
  - http: HTTP
toc_footers: []
includes: []
search: true
highlight_theme: darkula
headingLevel: 2

---

<h1 id="urbs-setting">urbs-setting v1.3.0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Urbs 灰度平台灰度策略服务

Base URLs:

* <a href="https://urbs-setting:8443">https://urbs-setting:8443</a>

# Authentication

- HTTP Authentication, scheme: bearer 

<h1 id="urbs-setting-version">Version</h1>

获取 urbs-setting 服务版本信息

## 获取版本信息

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/version \
  -H 'Accept: application/json'

```

```http
GET https://urbs-setting:8443/version HTTP/1.1
Host: urbs-setting:8443
Accept: application/json

```

`GET /version`

> Example responses

> 200 Response

```json
{
  "name": "urbs-setting",
  "version": "v1.2.0",
  "gitSHA1": "cd7e82a",
  "buildTime": "2020-05-08T05:38:45.231Z"
}
```

<h3 id="获取版本信息-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|version 返回结果|[Version](#schemaversion)|

<aside class="success">
This operation does not require authentication
</aside>

## 服务健康检查接口

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/healthz \
  -H 'Accept: application/json'

```

```http
GET https://urbs-setting:8443/healthz HTTP/1.1
Host: urbs-setting:8443
Accept: application/json

```

`GET /healthz`

> Example responses

> 200 Response

```json
{
  "dbConnect": true
}
```

<h3 id="服务健康检查接口-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Healthz 返回结果|Inline|

<h3 id="服务健康检查接口-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» dbConnect|boolean|false|none|是否连接了数据库|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="urbs-setting-user">User</h1>

User 用户相关接口

## 该接口为灰度网关提供用户的灰度信息，用于服务端灰度。获取指定 uid 用户在指定 product 产品下的所有（未分页，最多 400 条）灰度标签，包括从 group 群组继承的灰度标签，按照 label 指派时间反序。网关只会取匹配 client 和 channel 的第一条。标签列表不是实时数据，会被服务缓存，缓存时间在 config.cache_label_expire 配置，默认为 1 分钟，建议生产配置为 5 分钟。当 uid 对应用户不存在或 product 对应产品不存在时，该接口会返回空灰度标签列表

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/users/{uid}/labels:cache?product=string \
  -H 'Accept: application/json'

```

```http
GET https://urbs-setting:8443/users/{uid}/labels:cache?product=string HTTP/1.1
Host: urbs-setting:8443
Accept: application/json

```

`GET /users/{uid}/labels:cache`

<h3 id="该接口为灰度网关提供用户的灰度信息，用于服务端灰度。获取指定-uid-用户在指定-product-产品下的所有（未分页，最多-400-条）灰度标签，包括从-group-群组继承的灰度标签，按照-label-指派时间反序。网关只会取匹配-client-和-channel-的第一条。标签列表不是实时数据，会被服务缓存，缓存时间在-config.cache_label_expire-配置，默认为-1-分钟，建议生产配置为-5-分钟。当-uid-对应用户不存在或-product-对应产品不存在时，该接口会返回空灰度标签列表-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|uid|path|string|true|用户/群组 uid|
|product|query|string|true|产品名称|

> Example responses

> 200 Response

```json
{
  "timestamp": 1585129360,
  "result": [
    {
      "l": "beta",
      "cls": [
        "ios",
        "android",
        "web"
      ],
      "chs": [
        "stable",
        "beta",
        "dev"
      ]
    }
  ]
}
```

<h3 id="该接口为灰度网关提供用户的灰度信息，用于服务端灰度。获取指定-uid-用户在指定-product-产品下的所有（未分页，最多-400-条）灰度标签，包括从-group-群组继承的灰度标签，按照-label-指派时间反序。网关只会取匹配-client-和-channel-的第一条。标签列表不是实时数据，会被服务缓存，缓存时间在-config.cache_label_expire-配置，默认为-1-分钟，建议生产配置为-5-分钟。当-uid-对应用户不存在或-product-对应产品不存在时，该接口会返回空灰度标签列表-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|用于网关的用户灰度标签列表返回结果|Inline|

<h3 id="该接口为灰度网关提供用户的灰度信息，用于服务端灰度。获取指定-uid-用户在指定-product-产品下的所有（未分页，最多-400-条）灰度标签，包括从-group-群组继承的灰度标签，按照-label-指派时间反序。网关只会取匹配-client-和-channel-的第一条。标签列表不是实时数据，会被服务缓存，缓存时间在-config.cache_label_expire-配置，默认为-1-分钟，建议生产配置为-5-分钟。当-uid-对应用户不存在或-product-对应产品不存在时，该接口会返回空灰度标签列表-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» timestamp|integer(int64)|false|none|灰度标签列表缓存生成时间，1970 以来的秒数|
|» result|[[CacheLabelInfo](#schemacachelabelinfo)]|false|none|灰度标签列表|
|»» l|string|false|none|灰度标签名称|
|»» cls|[string]|false|none|灰度标签适用的 Clients 客户端类型列表，当列表为空时表示全适用|
|»» chs|[string]|false|none|灰度标签适用的 Channels 版本通道列表，当列表为空时表示全适用|

<aside class="success">
This operation does not require authentication
</aside>

## 获取用户列表，支持分页，按照标签指派时间倒序

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/users \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/users HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/users`

<h3 id="获取用户列表，支持分页，按照标签指派时间倒序-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "result": [
    {
      "uid": "5e82d747fe02a50021d339f3",
      "labels": "{\"teambition\":[{\"l\":\"dev\"}]}",
      "activeAt": 1585636012,
      "createdAt": "2020-05-08T05:38:45.237Z"
    }
  ]
}
```

<h3 id="获取用户列表，支持分页，按照标签指派时间倒序-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|用户列表|Inline|

<h3 id="获取用户列表，支持分页，按照标签指派时间倒序-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[[User](#schemauser)]|false|none|none|
|»» uid|string|false|none|用户的 uid|
|»» labels|string|false|none|缓存的用户当前被设置的 labels|
|»» activeAt|integer(int64)|false|none|用户活跃时间戳，1970 以来的秒数，非精确值|
|»» createdAt|string(date-time)|false|none|用户创建时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 该接口为客户端提供用户的产品功能模块配置项信息，用于客户端功能灰度。获取指定 uid 用户在指定 product 产品下的功能模块配置项信息列表，包括从 group 群组继承的配置项信息列表，按照 setting 值更新时间 updatedAt 反序。该 API 支持分页，默认获取最新更新的前 10 条，分页参数 nextPageToken 为更新时间 updatedAt 值（进行了 encodeURI 转义）。如果客户端本地缓存了 setting 列表，可以判断 nextPageToken 的值，如果 **为空** 或者其值小于本地缓存的最大 updatedAt 值，就不用读取下一页了。该 API 还支持 channel 和 client 参数，让客户端只读取匹配 client 和 channel 的 setting 列表。当 uid 对应用户不存在时，该接口会返回空配置项列表

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/users/{uid}/settings:unionAll?product=string \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/users/{uid}/settings:unionAll?product=string HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/users/{uid}/settings:unionAll`

<h3 id="该接口为客户端提供用户的产品功能模块配置项信息，用于客户端功能灰度。获取指定-uid-用户在指定-product-产品下的功能模块配置项信息列表，包括从-group-群组继承的配置项信息列表，按照-setting-值更新时间-updatedat-反序。该-api-支持分页，默认获取最新更新的前-10-条，分页参数-nextpagetoken-为更新时间-updatedat-值（进行了-encodeuri-转义）。如果客户端本地缓存了-setting-列表，可以判断-nextpagetoken-的值，如果-**为空**-或者其值小于本地缓存的最大-updatedat-值，就不用读取下一页了。该-api-还支持-channel-和-client-参数，让客户端只读取匹配-client-和-channel-的-setting-列表。当-uid-对应用户不存在时，该接口会返回空配置项列表-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|
|product|query|string|true|产品名称|
|channel|query|string|false|版本通道，必须为服务端配置的有效值，只返回匹配 channel 的 setting 列表|
|client|query|string|false|客户端类型，必须为服务端配置的有效值，只返回匹配 client 的 setting 列表|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|

> Example responses

> 200 Response

```json
{
  "nextPageToken": "",
  "result": [
    {
      "hid": "AwAAAAAAAAB25V_QnbhCuRwF",
      "product": "teambition",
      "module": "task",
      "name": "task-share",
      "desc": "string",
      "value": "disable",
      "lastValue": "",
      "release": 1,
      "assignedAt": "2020-05-08T05:38:45.237Z"
    }
  ]
}
```

<h3 id="该接口为客户端提供用户的产品功能模块配置项信息，用于客户端功能灰度。获取指定-uid-用户在指定-product-产品下的功能模块配置项信息列表，包括从-group-群组继承的配置项信息列表，按照-setting-值更新时间-updatedat-反序。该-api-支持分页，默认获取最新更新的前-10-条，分页参数-nextpagetoken-为更新时间-updatedat-值（进行了-encodeuri-转义）。如果客户端本地缓存了-setting-列表，可以判断-nextpagetoken-的值，如果-**为空**-或者其值小于本地缓存的最大-updatedat-值，就不用读取下一页了。该-api-还支持-channel-和-client-参数，让客户端只读取匹配-client-和-channel-的-setting-列表。当-uid-对应用户不存在时，该接口会返回空配置项列表-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|用户或群组被指派的配置项列表返回结果|Inline|

<h3 id="该接口为客户端提供用户的产品功能模块配置项信息，用于客户端功能灰度。获取指定-uid-用户在指定-product-产品下的功能模块配置项信息列表，包括从-group-群组继承的配置项信息列表，按照-setting-值更新时间-updatedat-反序。该-api-支持分页，默认获取最新更新的前-10-条，分页参数-nextpagetoken-为更新时间-updatedat-值（进行了-encodeuri-转义）。如果客户端本地缓存了-setting-列表，可以判断-nextpagetoken-的值，如果-**为空**-或者其值小于本地缓存的最大-updatedat-值，就不用读取下一页了。该-api-还支持-channel-和-client-参数，让客户端只读取匹配-client-和-channel-的-setting-列表。当-uid-对应用户不存在时，该接口会返回空配置项列表-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[MySetting](#schemamysetting)]|false|none|none|
|»» hid|string|false|none|配置项的 hid|
|»» product|string|false|none|配置项所属的产品名称|
|»» module|string|false|none|配置项所属的功能模块名称|
|»» name|string|false|none|配置项名称|
|»» desc|string|false|none|配置项描述，|
|»» value|string|false|none|配置项值|
|»» lastValue|string|false|none|配置项值|
|»» release|integer(int64)|false|none|被设置批次|
|»» assignedAt|string(date-time)|false|none|被设置时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 获取指定 uid 用户灰度标签列表，不包括从群组继承的灰度标签，支持分页，按照标签指派时间倒序

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/users/{uid}/labels \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/users/{uid}/labels HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/users/{uid}/labels`

<h3 id="获取指定-uid-用户灰度标签列表，不包括从群组继承的灰度标签，支持分页，按照标签指派时间倒序-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "nextPageToken": "",
  "result": [
    {
      "hid": "AwAAAAAAAAB25V_QnbhCuRwF",
      "product": "teambition",
      "name": "task-share",
      "desc": "string",
      "release": 1,
      "assignedAt": "2020-05-08T05:38:45.238Z"
    }
  ]
}
```

<h3 id="获取指定-uid-用户灰度标签列表，不包括从群组继承的灰度标签，支持分页，按照标签指派时间倒序-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|用户或群组被指派的灰度标签列表返回结果|Inline|

<h3 id="获取指定-uid-用户灰度标签列表，不包括从群组继承的灰度标签，支持分页，按照标签指派时间倒序-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[MyLabel](#schemamylabel)]|false|none|none|
|»» hid|string|false|none|灰度标签的 hid|
|»» product|string|false|none|灰度标签所属的产品名称|
|»» name|string|false|none|灰度标签名称|
|»» desc|string|false|none|灰度标签描述|
|»» release|integer(int64)|false|none|被设置批次|
|»» assignedAt|string(date-time)|false|none|被设置时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 强制刷新指定用户的灰度标签列表缓存

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/users/{uid}/labels:cache \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/users/{uid}/labels:cache HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`PUT /v1/users/{uid}/labels:cache`

<h3 id="强制刷新指定用户的灰度标签列表缓存-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|

> Example responses

> 200 Response

```json
{
  "result": {
    "uid": "5e82d747fe02a50021d339f3",
    "labels": "{\"teambition\":[{\"l\":\"dev\"}]}",
    "activeAt": 1585636012,
    "createdAt": "2020-05-08T05:38:45.238Z"
  }
}
```

<h3 id="强制刷新指定用户的灰度标签列表缓存-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|单个用户结果|Inline|

<h3 id="强制刷新指定用户的灰度标签列表缓存-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[User](#schemauser)|false|none|none|
|»» uid|string|false|none|用户的 uid|
|»» labels|string|false|none|缓存的用户当前被设置的 labels|
|»» activeAt|integer(int64)|false|none|用户活跃时间戳，1970 以来的秒数，非精确值|
|»» createdAt|string(date-time)|false|none|用户创建时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 获取指定 uid 用户在指定产品线下的功能模块配置项列表，不包括从群组继承的配置项，支持分页，按照配置项指派时间倒序

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/users/{uid}/settings \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/users/{uid}/settings HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/users/{uid}/settings`

<h3 id="获取指定-uid-用户在指定产品线下的功能模块配置项列表，不包括从群组继承的配置项，支持分页，按照配置项指派时间倒序-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "nextPageToken": "",
  "result": [
    {
      "hid": "AwAAAAAAAAB25V_QnbhCuRwF",
      "product": "teambition",
      "module": "task",
      "name": "task-share",
      "desc": "string",
      "value": "disable",
      "lastValue": "",
      "release": 1,
      "assignedAt": "2020-05-08T05:38:45.239Z"
    }
  ]
}
```

<h3 id="获取指定-uid-用户在指定产品线下的功能模块配置项列表，不包括从群组继承的配置项，支持分页，按照配置项指派时间倒序-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|用户或群组被指派的配置项列表返回结果|Inline|

<h3 id="获取指定-uid-用户在指定产品线下的功能模块配置项列表，不包括从群组继承的配置项，支持分页，按照配置项指派时间倒序-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[MySetting](#schemamysetting)]|false|none|none|
|»» hid|string|false|none|配置项的 hid|
|»» product|string|false|none|配置项所属的产品名称|
|»» module|string|false|none|配置项所属的功能模块名称|
|»» name|string|false|none|配置项名称|
|»» desc|string|false|none|配置项描述，|
|»» value|string|false|none|配置项值|
|»» lastValue|string|false|none|配置项值|
|»» release|integer(int64)|false|none|被设置批次|
|»» assignedAt|string(date-time)|false|none|被设置时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 判断指定 uid 用户是否存在

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/users/{uid}/exists \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/users/{uid}/exists HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/users/{uid}/exists`

<h3 id="判断指定-uid-用户是否存在-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="判断指定-uid-用户是否存在-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="判断指定-uid-用户是否存在-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 批量添加用户，忽略已存在的用户

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/users:batch \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/users:batch HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`POST /v1/users:batch`

> Body parameter

```json
{
  "users": [
    "50c32afae8cf1439d35a87e6",
    "5e69a9bd6ac3cd00213ea969"
  ]
}
```

<h3 id="批量添加用户，忽略已存在的用户-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|body|body|[UsersBody](#schemausersbody)|true|批量添加用户请求数据|
|» users|body|[string]|false|用户 uid 数组，必须符合正则 /^[0-9A-Za-z._=-]{3,63}$/|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="批量添加用户，忽略已存在的用户-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="批量添加用户，忽略已存在的用户-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 移除指定 uid 用户的指定 hid 灰度标签

> Code samples

```shell
# You can also use wget
curl -X DELETE https://urbs-setting:8443/v1/users/{uid}/labels/{hid} \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
DELETE https://urbs-setting:8443/v1/users/{uid}/labels/{hid} HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`DELETE /v1/users/{uid}/labels/{hid}`

<h3 id="移除指定-uid-用户的指定-hid-灰度标签-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="移除指定-uid-用户的指定-hid-灰度标签-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="移除指定-uid-用户的指定-hid-灰度标签-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 回滚指定 uid 用户的指定 hid 配置项值到上一个，只能回退到上一个值，不能到上上个值

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/users/{uid}/settings/{hid}:rollback \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/users/{uid}/settings/{hid}:rollback HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`PUT /v1/users/{uid}/settings/{hid}:rollback`

<h3 id="回滚指定-uid-用户的指定-hid-配置项值到上一个，只能回退到上一个值，不能到上上个值-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="回滚指定-uid-用户的指定-hid-配置项值到上一个，只能回退到上一个值，不能到上上个值-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="回滚指定-uid-用户的指定-hid-配置项值到上一个，只能回退到上一个值，不能到上上个值-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 移除指定 uid 用户的指定 hid 配置项

> Code samples

```shell
# You can also use wget
curl -X DELETE https://urbs-setting:8443/v1/users/{uid}/settings/{hid} \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
DELETE https://urbs-setting:8443/v1/users/{uid}/settings/{hid} HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`DELETE /v1/users/{uid}/settings/{hid}`

<h3 id="移除指定-uid-用户的指定-hid-配置项-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="移除指定-uid-用户的指定-hid-配置项-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="移除指定-uid-用户的指定-hid-配置项-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

<h1 id="urbs-setting-group">Group</h1>

Group 群组相关接口

## 获取指定 uid 群组灰度标签列表，支持分页，按照标签指派时间倒序。

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/groups/{uid}/labels \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/groups/{uid}/labels HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/groups/{uid}/labels`

<h3 id="获取指定-uid-群组灰度标签列表，支持分页，按照标签指派时间倒序。-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "nextPageToken": "",
  "result": [
    {
      "hid": "AwAAAAAAAAB25V_QnbhCuRwF",
      "product": "teambition",
      "name": "task-share",
      "desc": "string",
      "release": 1,
      "assignedAt": "2020-05-08T05:38:45.240Z"
    }
  ]
}
```

<h3 id="获取指定-uid-群组灰度标签列表，支持分页，按照标签指派时间倒序。-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|用户或群组被指派的灰度标签列表返回结果|Inline|

<h3 id="获取指定-uid-群组灰度标签列表，支持分页，按照标签指派时间倒序。-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[MyLabel](#schemamylabel)]|false|none|none|
|»» hid|string|false|none|灰度标签的 hid|
|»» product|string|false|none|灰度标签所属的产品名称|
|»» name|string|false|none|灰度标签名称|
|»» desc|string|false|none|灰度标签描述|
|»» release|integer(int64)|false|none|被设置批次|
|»» assignedAt|string(date-time)|false|none|被设置时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 获取指定 uid 群组的功能模块配置项列表，支持分页，按照配置项指派时间倒序

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/groups/{uid}/settings \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/groups/{uid}/settings HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/groups/{uid}/settings`

<h3 id="获取指定-uid-群组的功能模块配置项列表，支持分页，按照配置项指派时间倒序-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "nextPageToken": "",
  "result": [
    {
      "hid": "AwAAAAAAAAB25V_QnbhCuRwF",
      "product": "teambition",
      "module": "task",
      "name": "task-share",
      "desc": "string",
      "value": "disable",
      "lastValue": "",
      "release": 1,
      "assignedAt": "2020-05-08T05:38:45.241Z"
    }
  ]
}
```

<h3 id="获取指定-uid-群组的功能模块配置项列表，支持分页，按照配置项指派时间倒序-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|用户或群组被指派的配置项列表返回结果|Inline|

<h3 id="获取指定-uid-群组的功能模块配置项列表，支持分页，按照配置项指派时间倒序-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[MySetting](#schemamysetting)]|false|none|none|
|»» hid|string|false|none|配置项的 hid|
|»» product|string|false|none|配置项所属的产品名称|
|»» module|string|false|none|配置项所属的功能模块名称|
|»» name|string|false|none|配置项名称|
|»» desc|string|false|none|配置项描述，|
|»» value|string|false|none|配置项值|
|»» lastValue|string|false|none|配置项值|
|»» release|integer(int64)|false|none|被设置批次|
|»» assignedAt|string(date-time)|false|none|被设置时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 读取群组列表，支持分页，按照创建时间倒序

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/groups \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/groups HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/groups`

<h3 id="读取群组列表，支持分页，按照创建时间倒序-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|kind|query|string|false|查询指定 kind 类型的群组，未提供则查询所有类型|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "totalSize": 1,
  "nextPageToken": "",
  "result": [
    {
      "uid": "5e82d747fe02a50021d339f3",
      "kind": "organization",
      "desc": "string",
      "syncAt": 1585636012,
      "status": 99,
      "createdAt": "2020-05-08T05:38:45.241Z",
      "updatedAt": "2020-05-08T05:38:45.241Z"
    }
  ]
}
```

<h3 id="读取群组列表，支持分页，按照创建时间倒序-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|群组列表返回结果|Inline|

<h3 id="读取群组列表，支持分页，按照创建时间倒序-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» totalSize|[TotalSize](#schematotalsize)(int64)|false|none|当前分页查询的总数据量|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[Group](#schemagroup)]|false|none|none|
|»» uid|string|false|none|群组的 uid|
|»» kind|string|false|none|群组类型|
|»» desc|string|false|none|群组的描述|
|»» syncAt|integer(int64)|false|none|群组成员同步时间点，1970 以来的秒数|
|»» status|integer(int64)|false|none|群组成员数量，非精确值|
|»» createdAt|string(date-time)|false|none|灰度标签创建时间|
|»» updatedAt|string(date-time)|false|none|灰度标签更新时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 判断指定 uid 群组是否存在

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/groups/{uid}/exists \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/groups/{uid}/exists HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/groups/{uid}/exists`

<h3 id="判断指定-uid-群组是否存在-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="判断指定-uid-群组是否存在-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="判断指定-uid-群组是否存在-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 批量添加群组，忽略已存在的群组，群组 uid 必须唯一

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/groups:batch \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/groups:batch HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`POST /v1/groups:batch`

> Body parameter

```json
{
  "groups": [
    {
      "uid": "50c32afae8cf1439d35a87e6",
      "kind": "organization"
    }
  ]
}
```

<h3 id="批量添加群组，忽略已存在的群组，群组-uid-必须唯一-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|body|body|[GroupsBody](#schemagroupsbody)|true|批量添加群组请求数据|
|» groups|body|[object]|false|群组信息数组|
|»» uid|body|string|false|群组 uid，必须符合正则 /^[0-9A-Za-z._=-]{3,63}$/|
|»» kind|body|string|false|群组类型|
|»» desc|body|string|false|群组描述|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="批量添加群组，忽略已存在的群组，群组-uid-必须唯一-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="批量添加群组，忽略已存在的群组，群组-uid-必须唯一-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 更新指定 uid 群组

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/groups/{uid} \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/groups/{uid} HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`PUT /v1/groups/{uid}`

> Body parameter

```json
{
  "syncAt": 1585638012
}
```

<h3 id="更新指定-uid-群组-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|
|body|body|[GroupUpdateBody](#schemagroupupdatebody)|true|更新群组请求数据|
|» syncAt|body|integer(int64)|false|群组成员同步时间点，1970 以来的秒数|
|» desc|body|string|false|群组描述|

> Example responses

> 200 Response

```json
{
  "result": {
    "uid": "5e82d747fe02a50021d339f3",
    "kind": "organization",
    "desc": "string",
    "syncAt": 1585636012,
    "status": 99,
    "createdAt": "2020-05-08T05:38:45.243Z",
    "updatedAt": "2020-05-08T05:38:45.243Z"
  }
}
```

<h3 id="更新指定-uid-群组-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|单个群组返回结果|Inline|

<h3 id="更新指定-uid-群组-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[Group](#schemagroup)|false|none|none|
|»» uid|string|false|none|群组的 uid|
|»» kind|string|false|none|群组类型|
|»» desc|string|false|none|群组的描述|
|»» syncAt|integer(int64)|false|none|群组成员同步时间点，1970 以来的秒数|
|»» status|integer(int64)|false|none|群组成员数量，非精确值|
|»» createdAt|string(date-time)|false|none|灰度标签创建时间|
|»» updatedAt|string(date-time)|false|none|灰度标签更新时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 删除指定 uid 群组

> Code samples

```shell
# You can also use wget
curl -X DELETE https://urbs-setting:8443/v1/groups/{uid} \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
DELETE https://urbs-setting:8443/v1/groups/{uid} HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`DELETE /v1/groups/{uid}`

<h3 id="删除指定-uid-群组-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="删除指定-uid-群组-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="删除指定-uid-群组-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 移除指定 uid 群组的指定 hid 灰度标签

> Code samples

```shell
# You can also use wget
curl -X DELETE https://urbs-setting:8443/v1/groups/{uid}/labels/{hid} \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
DELETE https://urbs-setting:8443/v1/groups/{uid}/labels/{hid} HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`DELETE /v1/groups/{uid}/labels/{hid}`

<h3 id="移除指定-uid-群组的指定-hid-灰度标签-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="移除指定-uid-群组的指定-hid-灰度标签-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="移除指定-uid-群组的指定-hid-灰度标签-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 回滚指定 uid 群组的指定 hid 配置项值到上一个，只能回退到上一个值，不能到上上个值

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/groups/{uid}/settings/{hid}:rollback \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/groups/{uid}/settings/{hid}:rollback HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`PUT /v1/groups/{uid}/settings/{hid}:rollback`

<h3 id="回滚指定-uid-群组的指定-hid-配置项值到上一个，只能回退到上一个值，不能到上上个值-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="回滚指定-uid-群组的指定-hid-配置项值到上一个，只能回退到上一个值，不能到上上个值-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="回滚指定-uid-群组的指定-hid-配置项值到上一个，只能回退到上一个值，不能到上上个值-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 移除指定 uid 群组的指定 hid 配置项

> Code samples

```shell
# You can also use wget
curl -X DELETE https://urbs-setting:8443/v1/groups/{uid}/settings/{hid} \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
DELETE https://urbs-setting:8443/v1/groups/{uid}/settings/{hid} HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`DELETE /v1/groups/{uid}/settings/{hid}`

<h3 id="移除指定-uid-群组的指定-hid-配置项-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="移除指定-uid-群组的指定-hid-配置项-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="移除指定-uid-群组的指定-hid-配置项-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 批量添加群组成员，如果群组成员已存在，则会更新成员的 syncAt 值为 group 的 syncAt 值

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/groups/{uid}/members:batch \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/groups/{uid}/members:batch HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`POST /v1/groups/{uid}/members:batch`

> Body parameter

```json
{
  "users": [
    "50c32afae8cf1439d35a87e6",
    "5e69a9bd6ac3cd00213ea969"
  ]
}
```

<h3 id="批量添加群组成员，如果群组成员已存在，则会更新成员的-syncat-值为-group-的-syncat-值-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|body|body|[UsersBody](#schemausersbody)|true|批量添加用户请求数据|
|» users|body|[string]|false|用户 uid 数组，必须符合正则 /^[0-9A-Za-z._=-]{3,63}$/|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="批量添加群组成员，如果群组成员已存在，则会更新成员的-syncat-值为-group-的-syncat-值-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="批量添加群组成员，如果群组成员已存在，则会更新成员的-syncat-值为-group-的-syncat-值-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 获取指定 uid 群组的成员列表，支持分页，按照成员添加时间倒序

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/groups/{uid}/members \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/groups/{uid}/members HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/groups/{uid}/members`

<h3 id="获取指定-uid-群组的成员列表，支持分页，按照成员添加时间倒序-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "totalSize": 1,
  "nextPageToken": "",
  "result": [
    {
      "user": "5e82d747fe02a50021d339f3",
      "syncAt": 1585636012,
      "createdAt": "2020-05-08T05:38:45.245Z"
    }
  ]
}
```

<h3 id="获取指定-uid-群组的成员列表，支持分页，按照成员添加时间倒序-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|群组成员列表返回结果|Inline|

<h3 id="获取指定-uid-群组的成员列表，支持分页，按照成员添加时间倒序-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» totalSize|[TotalSize](#schematotalsize)(int64)|false|none|当前分页查询的总数据量|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[GroupMember](#schemagroupmember)]|false|none|none|
|»» user|string|false|none|群组成员的用户 uid|
|»» syncAt|integer(int64)|false|none|该群组成员同步时间，1970 以来的秒数|
|»» createdAt|string(date-time)|false|none|该群组成员添加时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 移除群组指定 user 的成员或批量移除同步时间点小于 syncLt 的成员

> Code samples

```shell
# You can also use wget
curl -X DELETE https://urbs-setting:8443/v1/groups/{uid}/members \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
DELETE https://urbs-setting:8443/v1/groups/{uid}/members HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`DELETE /v1/groups/{uid}/members`

<h3 id="移除群组指定-user-的成员或批量移除同步时间点小于-synclt-的成员-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户/群组 uid|
|user|query|string|false|移除群组指定 user 的成员|
|syncLt|query|string(date-time)|false|批量移除同步时间点小于 syncLt 的成员|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="移除群组指定-user-的成员或批量移除同步时间点小于-synclt-的成员-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="移除群组指定-user-的成员或批量移除同步时间点小于-synclt-的成员-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

<h1 id="urbs-setting-product">Product</h1>

Product 产品相关接口

## 读取产品列表，支持分页，按照创建时间倒序

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/products \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/products HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/products`

<h3 id="读取产品列表，支持分页，按照创建时间倒序-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "totalSize": 1,
  "nextPageToken": "",
  "result": [
    {
      "name": "urbs",
      "desc": "string",
      "status": 0,
      "createdAt": "2020-05-08T05:38:45.248Z",
      "updatedAt": "2020-05-08T05:38:45.248Z",
      "deletedAt": null,
      "offlineAt": null
    }
  ]
}
```

<h3 id="读取产品列表，支持分页，按照创建时间倒序-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|产品列表返回结果|Inline|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|标准错误返回结果|Inline|

<h3 id="读取产品列表，支持分页，按照创建时间倒序-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» totalSize|[TotalSize](#schematotalsize)(int64)|false|none|当前分页查询的总数据量|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[Product](#schemaproduct)]|false|none|none|
|»» name|string|false|none|产品名称|
|»» desc|string|false|none|产品的描述|
|»» status|integer(int64)|false|none|产品状态值|
|»» createdAt|string(date-time)|false|none|产品创建时间|
|»» updatedAt|string(date-time)|false|none|产品更新时间|
|»» deletedAt|string(date-time)|false|none|产品删除时间|
|»» offlineAt|string(date-time)|false|none|产品下线时间|

Status Code **401**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» error|string|false|none|错误代号|
|» message|string|false|none|错误详情|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 添加产品，产品 name 必须唯一

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/products \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/products HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`POST /v1/products`

> Body parameter

```json
{
  "name": "product or module or setting"
}
```

<h3 id="添加产品，产品-name-必须唯一-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|body|body|[NameDescBody](#schemanamedescbody)|true|创建产品、功能模块、配置项的请求数据|
|» name|body|string|false|名称，必须符合正则 /^[0-9a-z][0-9a-z.-]{0,61}[0-9a-z]$/|
|» desc|body|string|false|描述|

> Example responses

> 200 Response

```json
{
  "result": {
    "name": "urbs",
    "desc": "string",
    "status": 0,
    "createdAt": "2020-05-08T05:38:45.248Z",
    "updatedAt": "2020-05-08T05:38:45.248Z",
    "deletedAt": null,
    "offlineAt": null
  }
}
```

<h3 id="添加产品，产品-name-必须唯一-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|单个产品返回结果|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|标准错误返回结果|Inline|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|标准错误返回结果|Inline|

<h3 id="添加产品，产品-name-必须唯一-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[Product](#schemaproduct)|false|none|none|
|»» name|string|false|none|产品名称|
|»» desc|string|false|none|产品的描述|
|»» status|integer(int64)|false|none|产品状态值|
|»» createdAt|string(date-time)|false|none|产品创建时间|
|»» updatedAt|string(date-time)|false|none|产品更新时间|
|»» deletedAt|string(date-time)|false|none|产品删除时间|
|»» offlineAt|string(date-time)|false|none|产品下线时间|

Status Code **400**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» error|string|false|none|错误代号|
|» message|string|false|none|错误详情|

Status Code **401**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» error|string|false|none|错误代号|
|» message|string|false|none|错误详情|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 更新指定 product name 的产品

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/products/{product} \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/products/{product} HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`PUT /v1/products/{product}`

> Body parameter

```json
{
  "desc": "Urbs 产品线，负责人：XXX"
}
```

<h3 id="更新指定-product-name-的产品-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|body|body|[ProductUpdateBody](#schemaproductupdatebody)|true|更新产品请求数据|
|» desc|body|string|false|产品描述|

> Example responses

> 200 Response

```json
{
  "result": {
    "uid": "5e82d747fe02a50021d339f3",
    "kind": "organization",
    "desc": "string",
    "syncAt": 1585636012,
    "status": 99,
    "createdAt": "2020-05-08T05:38:45.248Z",
    "updatedAt": "2020-05-08T05:38:45.248Z"
  }
}
```

<h3 id="更新指定-product-name-的产品-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|单个群组返回结果|Inline|

<h3 id="更新指定-product-name-的产品-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[Group](#schemagroup)|false|none|none|
|»» uid|string|false|none|群组的 uid|
|»» kind|string|false|none|群组类型|
|»» desc|string|false|none|群组的描述|
|»» syncAt|integer(int64)|false|none|群组成员同步时间点，1970 以来的秒数|
|»» status|integer(int64)|false|none|群组成员数量，非精确值|
|»» createdAt|string(date-time)|false|none|灰度标签创建时间|
|»» updatedAt|string(date-time)|false|none|灰度标签更新时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 删除指定 product name 的产品，产品必须下线后才能被删除

> Code samples

```shell
# You can also use wget
curl -X DELETE https://urbs-setting:8443/v1/products/{product} \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
DELETE https://urbs-setting:8443/v1/products/{product} HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`DELETE /v1/products/{product}`

<h3 id="删除指定-product-name-的产品，产品必须下线后才能被删除-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="删除指定-product-name-的产品，产品必须下线后才能被删除-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="删除指定-product-name-的产品，产品必须下线后才能被删除-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 将指定 product name 的产品下线，此操作会将产品名下的所有功能模块和配置项都下线，所有设置给用户或群组的对应配置项和灰度标签也会被移除！

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/products/{product}:offline \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/products/{product}:offline HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`PUT /v1/products/{product}:offline`

<h3 id="将指定-product-name-的产品下线，此操作会将产品名下的所有功能模块和配置项都下线，所有设置给用户或群组的对应配置项和灰度标签也会被移除！-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="将指定-product-name-的产品下线，此操作会将产品名下的所有功能模块和配置项都下线，所有设置给用户或群组的对应配置项和灰度标签也会被移除！-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="将指定-product-name-的产品下线，此操作会将产品名下的所有功能模块和配置项都下线，所有设置给用户或群组的对应配置项和灰度标签也会被移除！-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 将指定 product name 的产品的统计数据

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/products/{product}/statistics \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/products/{product}/statistics HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`PUT /v1/products/{product}/statistics`

<h3 id="将指定-product-name-的产品的统计数据-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|

> Example responses

> 200 Response

```json
{
  "result": {
    "labels": 0,
    "modules": 0,
    "settings": 0,
    "release": 0,
    "status": 0
  }
}
```

<h3 id="将指定-product-name-的产品的统计数据-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|产品统计数据结果|Inline|

<h3 id="将指定-product-name-的产品的统计数据-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[ProductStatistics](#schemaproductstatistics)|false|none|none|
|»» labels|integer(int64)|false|none|产品下灰度标签数量|
|»» modules|integer(int64)|false|none|产品下功能模块数量|
|»» settings|integer(int64)|false|none|产品下配置项数量|
|»» release|integer(int64)|false|none|产品下灰度标签和配置项总发布次数|
|»» status|integer(int64)|false|none|产品下灰度标签和配置项总作用人数（非精确值）|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

<h1 id="urbs-setting-label">Label</h1>

Label 灰度标签相关接口

## 读取产品下灰度标签列表，支持分页，按照创建时间倒序

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/products/{product}/labels \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/products/{product}/labels HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/products/{product}/labels`

<h3 id="读取产品下灰度标签列表，支持分页，按照创建时间倒序-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "totalSize": 1,
  "nextPageToken": "",
  "result": [
    {
      "hid": "AwAAAAAAAAB25V_QnbhCuRwF",
      "product": "teambition",
      "name": "beta",
      "desc": "string",
      "channels": [
        "beta"
      ],
      "clients": [
        "web"
      ],
      "status": 100,
      "release": 2,
      "createdAt": "2020-05-08T05:38:45.250Z",
      "updatedAt": "2020-05-08T05:38:45.250Z",
      "offlineAt": null
    }
  ]
}
```

<h3 id="读取产品下灰度标签列表，支持分页，按照创建时间倒序-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|灰度标签列表返回结果|Inline|

<h3 id="读取产品下灰度标签列表，支持分页，按照创建时间倒序-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» totalSize|[TotalSize](#schematotalsize)(int64)|false|none|当前分页查询的总数据量|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[LabelInfo](#schemalabelinfo)]|false|none|none|
|»» hid|string|false|none|灰度标签的 hid|
|»» product|string|false|none|灰度标签所属的产品名称|
|»» name|string|false|none|灰度标签名称，同一产品下唯一（不能重名）|
|»» desc|string|false|none|灰度标签描述|
|»» channels|[string]|false|none|灰度标签适用版本通道|
|»» clients|[string]|false|none|灰度标签适用客户端类型|
|»» status|integer(int64)|false|none|使用用户计数（被动异步计算，非精确值）|
|»» release|integer(int64)|false|none|标签发布（被设置）计数|
|»» createdAt|string(date-time)|false|none|灰度标签创建时间|
|»» updatedAt|string(date-time)|false|none|灰度标签更新时间|
|»» offlineAt|string(date-time)|false|none|灰度标签下线时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 添加产品灰度标签，灰度标签 name 在产品下必须唯一

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/products/{product}/labels \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/products/{product}/labels HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`POST /v1/products/{product}/labels`

> Body parameter

```json
{
  "name": "product or module or setting"
}
```

<h3 id="添加产品灰度标签，灰度标签-name-在产品下必须唯一-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|body|body|[NameDescBody](#schemanamedescbody)|true|创建产品、功能模块、配置项的请求数据|
|» name|body|string|false|名称，必须符合正则 /^[0-9a-z][0-9a-z.-]{0,61}[0-9a-z]$/|
|» desc|body|string|false|描述|

> Example responses

> 200 Response

```json
{
  "result": {
    "hid": "AwAAAAAAAAB25V_QnbhCuRwF",
    "product": "teambition",
    "name": "beta",
    "desc": "string",
    "channels": [
      "beta"
    ],
    "clients": [
      "web"
    ],
    "status": 100,
    "release": 2,
    "createdAt": "2020-05-08T05:38:45.250Z",
    "updatedAt": "2020-05-08T05:38:45.250Z",
    "offlineAt": null
  }
}
```

<h3 id="添加产品灰度标签，灰度标签-name-在产品下必须唯一-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|灰度标签列表返回结果|Inline|

<h3 id="添加产品灰度标签，灰度标签-name-在产品下必须唯一-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[LabelInfo](#schemalabelinfo)|false|none|none|
|»» hid|string|false|none|灰度标签的 hid|
|»» product|string|false|none|灰度标签所属的产品名称|
|»» name|string|false|none|灰度标签名称，同一产品下唯一（不能重名）|
|»» desc|string|false|none|灰度标签描述|
|»» channels|[string]|false|none|灰度标签适用版本通道|
|»» clients|[string]|false|none|灰度标签适用客户端类型|
|»» status|integer(int64)|false|none|使用用户计数（被动异步计算，非精确值）|
|»» release|integer(int64)|false|none|标签发布（被设置）计数|
|»» createdAt|string(date-time)|false|none|灰度标签创建时间|
|»» updatedAt|string(date-time)|false|none|灰度标签更新时间|
|»» offlineAt|string(date-time)|false|none|灰度标签下线时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 更新指定 product name 的产品

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/products/{product}/labels/{label} \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/products/{product}/labels/{label} HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`PUT /v1/products/{product}/labels/{label}`

> Body parameter

```json
{
  "desc": null,
  "channels": [
    "beta"
  ],
  "clients": [
    "ios",
    "android"
  ]
}
```

<h3 id="更新指定-product-name-的产品-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|label|path|string|true|灰度标签名称|
|body|body|[LabelUpdateBody](#schemalabelupdatebody)|true|更新灰度标签的请求数据|
|» desc|body|string|false|灰度标签描述|
|» channels|body|[string]|false|该灰度标签适用的产品版本通道，必须是服务端配置的可用版本通道的子集，为空表示适用所有|
|» clients|body|[string]|false|该灰度标签适用的客户端类型，必须是服务端配置的可用客户端类型的子集，为空表示适用所有|

> Example responses

> 200 Response

```json
{
  "result": {
    "hid": "AwAAAAAAAAB25V_QnbhCuRwF",
    "product": "teambition",
    "name": "beta",
    "desc": "string",
    "channels": [
      "beta"
    ],
    "clients": [
      "web"
    ],
    "status": 100,
    "release": 2,
    "createdAt": "2020-05-08T05:38:45.251Z",
    "updatedAt": "2020-05-08T05:38:45.251Z",
    "offlineAt": null
  }
}
```

<h3 id="更新指定-product-name-的产品-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|灰度标签列表返回结果|Inline|

<h3 id="更新指定-product-name-的产品-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[LabelInfo](#schemalabelinfo)|false|none|none|
|»» hid|string|false|none|灰度标签的 hid|
|»» product|string|false|none|灰度标签所属的产品名称|
|»» name|string|false|none|灰度标签名称，同一产品下唯一（不能重名）|
|»» desc|string|false|none|灰度标签描述|
|»» channels|[string]|false|none|灰度标签适用版本通道|
|»» clients|[string]|false|none|灰度标签适用客户端类型|
|»» status|integer(int64)|false|none|使用用户计数（被动异步计算，非精确值）|
|»» release|integer(int64)|false|none|标签发布（被设置）计数|
|»» createdAt|string(date-time)|false|none|灰度标签创建时间|
|»» updatedAt|string(date-time)|false|none|灰度标签更新时间|
|»» offlineAt|string(date-time)|false|none|灰度标签下线时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 将指定产品灰度标签下线，所有设置给用户或群组的对应灰度标签也会被移除！

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/products/{product}/labels/{label}:offline \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/products/{product}/labels/{label}:offline HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`PUT /v1/products/{product}/labels/{label}:offline`

<h3 id="将指定产品灰度标签下线，所有设置给用户或群组的对应灰度标签也会被移除！-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|label|path|string|true|灰度标签名称|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="将指定产品灰度标签下线，所有设置给用户或群组的对应灰度标签也会被移除！-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="将指定产品灰度标签下线，所有设置给用户或群组的对应灰度标签也会被移除！-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 批量为用户或群组设置灰度标签

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/products/{product}/labels/{label}:assign \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/products/{product}/labels/{label}:assign HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`POST /v1/products/{product}/labels/{label}:assign`

> Body parameter

```json
{
  "users": [
    "5c4057f0be825b390667abee"
  ],
  "groups": [
    "5bdc1846cd57df001789c751",
    "5bdc1846cd57df001789c751"
  ],
  "value": "beta"
}
```

<h3 id="批量为用户或群组设置灰度标签-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|label|path|string|true|灰度标签名称|
|body|body|[UsersGroupsBody](#schemausersgroupsbody)|true|批量为用户或群组设置灰度标签或配置项的请求数据|
|» users|body|[string]|false|用户 uid 数组，可以不提供，最大长度受 HTTP 请求最大字节数限制|
|» groups|body|[string]|false|群组 uid 数组，可以不提供，最大长度受 HTTP 请求最大字节数限制|
|» value|body|string|false|配置项值，设置灰度标签时不必提供|

> Example responses

> 200 Response

```json
{
  "result": {
    "release": 2,
    "users": [
      "50c32afae8cf1439d35a87e6",
      "5e69a9bd6ac3cd00213ea969"
    ],
    "groups": [
      "string"
    ]
  }
}
```

<h3 id="批量为用户或群组设置灰度标签-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|设置灰度标签返回结果|Inline|

<h3 id="批量为用户或群组设置灰度标签-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[LabelReleaseInfo](#schemalabelreleaseinfo)|false|none|none|
|»» release|integer(int64)|false|none|发布批次（被设置）计数|
|»» users|[string]|false|none|用户 uid 数组|
|»» groups|[string]|false|none|群组 uid 数组|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 批量撤销对用户或群组设置的产品灰度标签

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/products/{product}/labels/{label}:recall \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/products/{product}/labels/{label}:recall HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`POST /v1/products/{product}/labels/{label}:recall`

> Body parameter

```json
{
  "release": 2
}
```

<h3 id="批量撤销对用户或群组设置的产品灰度标签-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|label|path|string|true|灰度标签名称|
|body|body|[RecallBody](#schemarecallbody)|true|撤销/回滚指定发布批次的灰度标签或配置项|
|» release|body|integer(int64)|false|发布批次（被设置）计数|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="批量撤销对用户或群组设置的产品灰度标签-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="批量撤销对用户或群组设置的产品灰度标签-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 读取指定产品灰度标签的用户列表

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/products/{product}/labels/{label}/users \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/products/{product}/labels/{label}/users HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/products/{product}/labels/{label}/users`

<h3 id="读取指定产品灰度标签的用户列表-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|label|path|string|true|灰度标签名称|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "result": [
    {
      "labelHID": "urbs",
      "assignedAt": "2020-05-08T05:38:45.253Z",
      "release": 2,
      "user": "50c32afae8cf1439d35a87e6"
    }
  ]
}
```

<h3 id="读取指定产品灰度标签的用户列表-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|灰度标签作用的用户列表|Inline|

<h3 id="读取指定产品灰度标签的用户列表-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[[LabelUserInfo](#schemalabeluserinfo)]|false|none|none|
|»» labelHID|string|false|none|灰度标签 hid|
|»» assignedAt|string(date-time)|false|none|灰度标签设置时间|
|»» release|integer(int64)|false|none|发布批次（被设置）计数|
|»» user|string|false|none|用户的 uid|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 读取指定产品灰度标签的群组列表

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/products/{product}/labels/{label}/groups \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/products/{product}/labels/{label}/groups HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/products/{product}/labels/{label}/groups`

<h3 id="读取指定产品灰度标签的群组列表-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|label|path|string|true|灰度标签名称|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "result": [
    {
      "labelHID": "urbs",
      "assignedAt": "2020-05-08T05:38:45.253Z",
      "release": 2,
      "group": "50c32afae8cf1439d35a87e6",
      "kind": "organization",
      "desc": "string",
      "status": 99
    }
  ]
}
```

<h3 id="读取指定产品灰度标签的群组列表-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|灰度标签作用的群组列表|Inline|

<h3 id="读取指定产品灰度标签的群组列表-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[[LabelGroupInfo](#schemalabelgroupinfo)]|false|none|none|
|»» labelHID|string|false|none|灰度标签 hid|
|»» assignedAt|string(date-time)|false|none|灰度标签设置时间|
|»» release|integer(int64)|false|none|发布批次（被设置）计数|
|»» group|string|false|none|群组的 uid|
|»» kind|string|false|none|群组类型|
|»» desc|string|false|none|群组的描述|
|»» status|integer(int64)|false|none|群组成员数量，非精确值|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 读取指定产品灰度标签的灰度发布规则列表

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/products/{product}/labels/{label}/rules \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/products/{product}/labels/{label}/rules HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/products/{product}/labels/{label}/rules`

<h3 id="读取指定产品灰度标签的灰度发布规则列表-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|label|path|string|true|灰度标签名称|

> Example responses

> 200 Response

```json
{
  "result": [
    {
      "hid": "urbs",
      "labelHID": "urbs",
      "kind": "userPercent",
      "rule": "{\"value\": 10}",
      "release": 2,
      "createdAt": "2020-05-08T05:38:45.254Z",
      "updatedAt": "2020-05-08T05:38:45.254Z"
    }
  ]
}
```

<h3 id="读取指定产品灰度标签的灰度发布规则列表-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|灰度标签的发布规则列表|Inline|

<h3 id="读取指定产品灰度标签的灰度发布规则列表-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[[LabelRuleInfo](#schemalabelruleinfo)]|false|none|none|
|»» hid|string|false|none|发布规则的 hid|
|»» labelHID|string|false|none|发布规则所属灰度标签的 hid|
|»» kind|string|false|none|发布规则类型，目前仅支持 "userPercent"|
|»» rule|object|false|none|发布规则内容，JSON 对象，具体格式取决于 kind|
|»»» value|integer(int64)|false|none|当 kind 为 "userPercent" 时，value 为百分比，取值 [0, 100]|
|»» release|integer(int64)|false|none|发布批次（被设置）计数|
|»» createdAt|string(date-time)|false|none|创建时间|
|»» updatedAt|string(date-time)|false|none|更新时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 创建指定产品灰度标签的灰度发布规则，同一个灰度标签同一种 kind 的发布规则只能创建一个

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/products/{product}/labels/{label}/rules \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/products/{product}/labels/{label}/rules HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`POST /v1/products/{product}/labels/{label}/rules`

> Body parameter

```json
{
  "kind": "userPercent",
  "rule": "{\"value\": 10}",
  "value": "x"
}
```

<h3 id="创建指定产品灰度标签的灰度发布规则，同一个灰度标签同一种-kind-的发布规则只能创建一个-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|label|path|string|true|灰度标签名称|
|body|body|[LabelRuleBody](#schemalabelrulebody)|true|创建/更新配置项的发布规则|
|» kind|body|string|false|发布规则类型，目前仅支持 "userPercent"|
|» rule|body|object|false|发布规则内容，JSON 对象，具体格式取决于 kind|
|»» value|body|integer(int64)|false|当 kind 为 "userPercent" 时，value 为百分比，取值 [0, 100]|
|» value|body|string|false|发布规则的配置项值|

> Example responses

> 200 Response

```json
{
  "result": {
    "hid": "urbs",
    "labelHID": "urbs",
    "kind": "userPercent",
    "rule": "{\"value\": 10}",
    "release": 2,
    "createdAt": "2020-05-08T05:38:45.254Z",
    "updatedAt": "2020-05-08T05:38:45.254Z"
  }
}
```

<h3 id="创建指定产品灰度标签的灰度发布规则，同一个灰度标签同一种-kind-的发布规则只能创建一个-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|灰度标签的发布规则结果|Inline|

<h3 id="创建指定产品灰度标签的灰度发布规则，同一个灰度标签同一种-kind-的发布规则只能创建一个-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[LabelRuleInfo](#schemalabelruleinfo)|false|none|none|
|»» hid|string|false|none|发布规则的 hid|
|»» labelHID|string|false|none|发布规则所属灰度标签的 hid|
|»» kind|string|false|none|发布规则类型，目前仅支持 "userPercent"|
|»» rule|object|false|none|发布规则内容，JSON 对象，具体格式取决于 kind|
|»»» value|integer(int64)|false|none|当 kind 为 "userPercent" 时，value 为百分比，取值 [0, 100]|
|»» release|integer(int64)|false|none|发布批次（被设置）计数|
|»» createdAt|string(date-time)|false|none|创建时间|
|»» updatedAt|string(date-time)|false|none|更新时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 更新指定产品灰度标签的灰度发布规则

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/products/{product}/labels/{label}/rules/{hid} \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/products/{product}/labels/{label}/rules/{hid} HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`PUT /v1/products/{product}/labels/{label}/rules/{hid}`

> Body parameter

```json
{
  "kind": "userPercent",
  "rule": "{\"value\": 10}",
  "value": "x"
}
```

<h3 id="更新指定产品灰度标签的灰度发布规则-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|label|path|string|true|灰度标签名称|
|uid|path|string|true|标签/配置项/发布规则 hid|
|body|body|[LabelRuleBody](#schemalabelrulebody)|true|创建/更新配置项的发布规则|
|» kind|body|string|false|发布规则类型，目前仅支持 "userPercent"|
|» rule|body|object|false|发布规则内容，JSON 对象，具体格式取决于 kind|
|»» value|body|integer(int64)|false|当 kind 为 "userPercent" 时，value 为百分比，取值 [0, 100]|
|» value|body|string|false|发布规则的配置项值|

> Example responses

> 200 Response

```json
{
  "result": {
    "hid": "urbs",
    "labelHID": "urbs",
    "kind": "userPercent",
    "rule": "{\"value\": 10}",
    "release": 2,
    "createdAt": "2020-05-08T05:38:45.255Z",
    "updatedAt": "2020-05-08T05:38:45.255Z"
  }
}
```

<h3 id="更新指定产品灰度标签的灰度发布规则-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|灰度标签的发布规则结果|Inline|

<h3 id="更新指定产品灰度标签的灰度发布规则-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[LabelRuleInfo](#schemalabelruleinfo)|false|none|none|
|»» hid|string|false|none|发布规则的 hid|
|»» labelHID|string|false|none|发布规则所属灰度标签的 hid|
|»» kind|string|false|none|发布规则类型，目前仅支持 "userPercent"|
|»» rule|object|false|none|发布规则内容，JSON 对象，具体格式取决于 kind|
|»»» value|integer(int64)|false|none|当 kind 为 "userPercent" 时，value 为百分比，取值 [0, 100]|
|»» release|integer(int64)|false|none|发布批次（被设置）计数|
|»» createdAt|string(date-time)|false|none|创建时间|
|»» updatedAt|string(date-time)|false|none|更新时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 删除指定产品灰度标签的灰度发布规则

> Code samples

```shell
# You can also use wget
curl -X DELETE https://urbs-setting:8443/v1/products/{product}/labels/{label}/rules/{hid} \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
DELETE https://urbs-setting:8443/v1/products/{product}/labels/{label}/rules/{hid} HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`DELETE /v1/products/{product}/labels/{label}/rules/{hid}`

<h3 id="删除指定产品灰度标签的灰度发布规则-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|label|path|string|true|灰度标签名称|
|uid|path|string|true|标签/配置项/发布规则 hid|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="删除指定产品灰度标签的灰度发布规则-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="删除指定产品灰度标签的灰度发布规则-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

<h1 id="urbs-setting-module">Module</h1>

Module 产品功能模块相关接口

## 读取产品功能模块列表，支持分页，按照创建时间倒序

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/products/{product}/modules \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/products/{product}/modules HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/products/{product}/modules`

<h3 id="读取产品功能模块列表，支持分页，按照创建时间倒序-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "totalSize": 1,
  "nextPageToken": "",
  "result": [
    {
      "name": "urbs",
      "desc": "string",
      "status": 0,
      "createdAt": "2020-05-08T05:38:45.256Z",
      "updatedAt": "2020-05-08T05:38:45.256Z",
      "offlineAt": null
    }
  ]
}
```

<h3 id="读取产品功能模块列表，支持分页，按照创建时间倒序-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|产品列表返回结果|Inline|

<h3 id="读取产品功能模块列表，支持分页，按照创建时间倒序-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» totalSize|[TotalSize](#schematotalsize)(int64)|false|none|当前分页查询的总数据量|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[Module](#schemamodule)]|false|none|none|
|»» name|string|false|none|功能模块名称|
|»» desc|string|false|none|功能模块的描述|
|»» status|integer(int64)|false|none|有效配置项计数（被动异步计算，非精确值）|
|»» createdAt|string(date-time)|false|none|功能模块创建时间|
|»» updatedAt|string(date-time)|false|none|功能模块更新时间|
|»» offlineAt|string(date-time)|false|none|功能模块下线时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 添加产品功能模块，功能模块 name 在产品下必须唯一

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/products/{product}/modules \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/products/{product}/modules HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`POST /v1/products/{product}/modules`

> Body parameter

```json
{
  "name": "product or module or setting"
}
```

<h3 id="添加产品功能模块，功能模块-name-在产品下必须唯一-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|body|body|[NameDescBody](#schemanamedescbody)|true|创建产品、功能模块、配置项的请求数据|
|» name|body|string|false|名称，必须符合正则 /^[0-9a-z][0-9a-z.-]{0,61}[0-9a-z]$/|
|» desc|body|string|false|描述|

> Example responses

> 200 Response

```json
{
  "result": {
    "name": "urbs",
    "desc": "string",
    "status": 0,
    "createdAt": "2020-05-08T05:38:45.256Z",
    "updatedAt": "2020-05-08T05:38:45.256Z",
    "offlineAt": null
  }
}
```

<h3 id="添加产品功能模块，功能模块-name-在产品下必须唯一-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|单个产品返回结果|Inline|

<h3 id="添加产品功能模块，功能模块-name-在产品下必须唯一-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[Module](#schemamodule)|false|none|none|
|»» name|string|false|none|功能模块名称|
|»» desc|string|false|none|功能模块的描述|
|»» status|integer(int64)|false|none|有效配置项计数（被动异步计算，非精确值）|
|»» createdAt|string(date-time)|false|none|功能模块创建时间|
|»» updatedAt|string(date-time)|false|none|功能模块更新时间|
|»» offlineAt|string(date-time)|false|none|功能模块下线时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 更新指定 product name 的产品

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/products/{product}/modules/{module} \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/products/{product}/modules/{module} HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`PUT /v1/products/{product}/modules/{module}`

> Body parameter

```json
{
  "desc": "Urbs 产品线 xxx 模块，负责人：XXX"
}
```

<h3 id="更新指定-product-name-的产品-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|body|body|[ModuleUpdateBody](#schemamoduleupdatebody)|true|更新功能模块请求数据|
|» desc|body|string|false|产品描述|

> Example responses

> 200 Response

```json
{
  "result": {
    "name": "urbs",
    "desc": "string",
    "status": 0,
    "createdAt": "2020-05-08T05:38:45.257Z",
    "updatedAt": "2020-05-08T05:38:45.257Z",
    "offlineAt": null
  }
}
```

<h3 id="更新指定-product-name-的产品-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|单个产品返回结果|Inline|

<h3 id="更新指定-product-name-的产品-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[Module](#schemamodule)|false|none|none|
|»» name|string|false|none|功能模块名称|
|»» desc|string|false|none|功能模块的描述|
|»» status|integer(int64)|false|none|有效配置项计数（被动异步计算，非精确值）|
|»» createdAt|string(date-time)|false|none|功能模块创建时间|
|»» updatedAt|string(date-time)|false|none|功能模块更新时间|
|»» offlineAt|string(date-time)|false|none|功能模块下线时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 将指定产品功能模块下线，此操作会将功能模块名下的所有配置项都下线，所有设置给用户或群组的对应配置项也会被移除！

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/products/{product}/modules/{module}:offline \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/products/{product}/modules/{module}:offline HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`PUT /v1/products/{product}/modules/{module}:offline`

<h3 id="将指定产品功能模块下线，此操作会将功能模块名下的所有配置项都下线，所有设置给用户或群组的对应配置项也会被移除！-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="将指定产品功能模块下线，此操作会将功能模块名下的所有配置项都下线，所有设置给用户或群组的对应配置项也会被移除！-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="将指定产品功能模块下线，此操作会将功能模块名下的所有配置项都下线，所有设置给用户或群组的对应配置项也会被移除！-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

<h1 id="urbs-setting-setting">Setting</h1>

Setting 产品功能模块配置项相关接口

## 读取产品的配置项列表，支持分页，按照创建时间倒序

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/products/{product}/settings \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/products/{product}/settings HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/products/{product}/settings`

<h3 id="读取产品的配置项列表，支持分页，按照创建时间倒序-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "totalSize": 1,
  "nextPageToken": "",
  "result": [
    {
      "hid": "urbs",
      "product": "urbs",
      "module": "urbs",
      "name": "urbs",
      "desc": "string",
      "status": 0,
      "release": 0,
      "channels": [
        "beta"
      ],
      "clients": [
        "ios",
        "android"
      ],
      "values": [
        "true",
        "false"
      ],
      "createdAt": "2020-05-08T05:38:45.258Z",
      "updatedAt": "2020-05-08T05:38:45.258Z",
      "offlineAt": null
    }
  ]
}
```

<h3 id="读取产品的配置项列表，支持分页，按照创建时间倒序-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|产品列表返回结果|Inline|

<h3 id="读取产品的配置项列表，支持分页，按照创建时间倒序-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» totalSize|[TotalSize](#schematotalsize)(int64)|false|none|当前分页查询的总数据量|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[SettingInfo](#schemasettinginfo)]|false|none|none|
|»» hid|string|false|none|配置项名称|
|»» product|string|false|none|配置项名称|
|»» module|string|false|none|配置项名称|
|»» name|string|false|none|配置项名称|
|»» desc|string|false|none|配置项的描述|
|»» status|integer(int64)|false|none|使用用户计数（被动异步计算，非精确值）|
|»» release|integer(int64)|false|none|配置项发布（被设置）计数|
|»» channels|[string]|false|none|该配置项适用的产品版本通道，是服务端配置的可用版本通道的子集，为空表示适用所有|
|»» clients|[string]|false|none|该配置项适用的客户端类型，是服务端配置的可用客户端类型的子集，为空表示适用所有|
|»» values|[string]|false|none|该配置项可选值列表，配置项指派给用户或群组时只能从该列表中选择合法值|
|»» createdAt|string(date-time)|false|none|配置项创建时间|
|»» updatedAt|string(date-time)|false|none|配置项更新时间|
|»» offlineAt|string(date-time)|false|none|配置项下线时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 读取产品功能模块配置项列表，支持分页，按照创建时间倒序

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/products/{product}/modules/{module}/settings`

<h3 id="读取产品功能模块配置项列表，支持分页，按照创建时间倒序-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "totalSize": 1,
  "nextPageToken": "",
  "result": [
    {
      "hid": "urbs",
      "product": "urbs",
      "module": "urbs",
      "name": "urbs",
      "desc": "string",
      "status": 0,
      "release": 0,
      "channels": [
        "beta"
      ],
      "clients": [
        "ios",
        "android"
      ],
      "values": [
        "true",
        "false"
      ],
      "createdAt": "2020-05-08T05:38:45.259Z",
      "updatedAt": "2020-05-08T05:38:45.259Z",
      "offlineAt": null
    }
  ]
}
```

<h3 id="读取产品功能模块配置项列表，支持分页，按照创建时间倒序-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|产品列表返回结果|Inline|

<h3 id="读取产品功能模块配置项列表，支持分页，按照创建时间倒序-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» totalSize|[TotalSize](#schematotalsize)(int64)|false|none|当前分页查询的总数据量|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[SettingInfo](#schemasettinginfo)]|false|none|none|
|»» hid|string|false|none|配置项名称|
|»» product|string|false|none|配置项名称|
|»» module|string|false|none|配置项名称|
|»» name|string|false|none|配置项名称|
|»» desc|string|false|none|配置项的描述|
|»» status|integer(int64)|false|none|使用用户计数（被动异步计算，非精确值）|
|»» release|integer(int64)|false|none|配置项发布（被设置）计数|
|»» channels|[string]|false|none|该配置项适用的产品版本通道，是服务端配置的可用版本通道的子集，为空表示适用所有|
|»» clients|[string]|false|none|该配置项适用的客户端类型，是服务端配置的可用客户端类型的子集，为空表示适用所有|
|»» values|[string]|false|none|该配置项可选值列表，配置项指派给用户或群组时只能从该列表中选择合法值|
|»» createdAt|string(date-time)|false|none|配置项创建时间|
|»» updatedAt|string(date-time)|false|none|配置项更新时间|
|»» offlineAt|string(date-time)|false|none|配置项下线时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 添加产品功能模块的配置项，配置项 name 在功能模块下必须唯一

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`POST /v1/products/{product}/modules/{module}/settings`

> Body parameter

```json
{
  "name": "product or module or setting"
}
```

<h3 id="添加产品功能模块的配置项，配置项-name-在功能模块下必须唯一-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|body|body|[NameDescBody](#schemanamedescbody)|true|创建产品、功能模块、配置项的请求数据|
|» name|body|string|false|名称，必须符合正则 /^[0-9a-z][0-9a-z.-]{0,61}[0-9a-z]$/|
|» desc|body|string|false|描述|

> Example responses

> 200 Response

```json
{
  "totalSize": 1,
  "nextPageToken": "",
  "result": [
    {
      "hid": "urbs",
      "product": "urbs",
      "module": "urbs",
      "name": "urbs",
      "desc": "string",
      "status": 0,
      "release": 0,
      "channels": [
        "beta"
      ],
      "clients": [
        "ios",
        "android"
      ],
      "values": [
        "true",
        "false"
      ],
      "createdAt": "2020-05-08T05:38:45.260Z",
      "updatedAt": "2020-05-08T05:38:45.260Z",
      "offlineAt": null
    }
  ]
}
```

<h3 id="添加产品功能模块的配置项，配置项-name-在功能模块下必须唯一-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|产品列表返回结果|Inline|

<h3 id="添加产品功能模块的配置项，配置项-name-在功能模块下必须唯一-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» totalSize|[TotalSize](#schematotalsize)(int64)|false|none|当前分页查询的总数据量|
|» nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|» result|[[SettingInfo](#schemasettinginfo)]|false|none|none|
|»» hid|string|false|none|配置项名称|
|»» product|string|false|none|配置项名称|
|»» module|string|false|none|配置项名称|
|»» name|string|false|none|配置项名称|
|»» desc|string|false|none|配置项的描述|
|»» status|integer(int64)|false|none|使用用户计数（被动异步计算，非精确值）|
|»» release|integer(int64)|false|none|配置项发布（被设置）计数|
|»» channels|[string]|false|none|该配置项适用的产品版本通道，是服务端配置的可用版本通道的子集，为空表示适用所有|
|»» clients|[string]|false|none|该配置项适用的客户端类型，是服务端配置的可用客户端类型的子集，为空表示适用所有|
|»» values|[string]|false|none|该配置项可选值列表，配置项指派给用户或群组时只能从该列表中选择合法值|
|»» createdAt|string(date-time)|false|none|配置项创建时间|
|»» updatedAt|string(date-time)|false|none|配置项更新时间|
|»» offlineAt|string(date-time)|false|none|配置项下线时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 读取指定配置项

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting} \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting} HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/products/{product}/modules/{module}/settings/{setting}`

<h3 id="读取指定配置项-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|setting|path|string|true|配置项名称|

> Example responses

> 200 Response

```json
{
  "result": {
    "hid": "urbs",
    "product": "urbs",
    "module": "urbs",
    "name": "urbs",
    "desc": "string",
    "status": 0,
    "release": 0,
    "channels": [
      "beta"
    ],
    "clients": [
      "ios",
      "android"
    ],
    "values": [
      "true",
      "false"
    ],
    "createdAt": "2020-05-08T05:38:45.260Z",
    "updatedAt": "2020-05-08T05:38:45.260Z",
    "offlineAt": null
  }
}
```

<h3 id="读取指定配置项-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|单个配置项返回结果|Inline|

<h3 id="读取指定配置项-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[SettingInfo](#schemasettinginfo)|false|none|none|
|»» hid|string|false|none|配置项名称|
|»» product|string|false|none|配置项名称|
|»» module|string|false|none|配置项名称|
|»» name|string|false|none|配置项名称|
|»» desc|string|false|none|配置项的描述|
|»» status|integer(int64)|false|none|使用用户计数（被动异步计算，非精确值）|
|»» release|integer(int64)|false|none|配置项发布（被设置）计数|
|»» channels|[string]|false|none|该配置项适用的产品版本通道，是服务端配置的可用版本通道的子集，为空表示适用所有|
|»» clients|[string]|false|none|该配置项适用的客户端类型，是服务端配置的可用客户端类型的子集，为空表示适用所有|
|»» values|[string]|false|none|该配置项可选值列表，配置项指派给用户或群组时只能从该列表中选择合法值|
|»» createdAt|string(date-time)|false|none|配置项创建时间|
|»» updatedAt|string(date-time)|false|none|配置项更新时间|
|»» offlineAt|string(date-time)|false|none|配置项下线时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 更新指定配置项

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting} \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting} HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`PUT /v1/products/{product}/modules/{module}/settings/{setting}`

> Body parameter

```json
{
  "values": [
    "a",
    "b"
  ]
}
```

<h3 id="更新指定配置项-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|setting|path|string|true|配置项名称|
|body|body|[SettingUpdateBody](#schemasettingupdatebody)|true|更新配置项请求数据|
|» desc|body|string|false|产品描述|
|» channels|body|[string]|false|该配置项适用的产品版本通道，必须是服务端配置的可用版本通道的子集，为空表示适用所有|
|» clients|body|[string]|false|该配置项适用的客户端类型，必须是服务端配置的可用客户端类型的子集，为空表示适用所有|
|» values|body|[string]|false|该配置项可选值列表，配置项指派给用户或群组时只能从该列表中选择合法值|

> Example responses

> 200 Response

```json
{
  "result": {
    "hid": "urbs",
    "product": "urbs",
    "module": "urbs",
    "name": "urbs",
    "desc": "string",
    "status": 0,
    "release": 0,
    "channels": [
      "beta"
    ],
    "clients": [
      "ios",
      "android"
    ],
    "values": [
      "true",
      "false"
    ],
    "createdAt": "2020-05-08T05:38:45.261Z",
    "updatedAt": "2020-05-08T05:38:45.261Z",
    "offlineAt": null
  }
}
```

<h3 id="更新指定配置项-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|单个配置项返回结果|Inline|

<h3 id="更新指定配置项-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[SettingInfo](#schemasettinginfo)|false|none|none|
|»» hid|string|false|none|配置项名称|
|»» product|string|false|none|配置项名称|
|»» module|string|false|none|配置项名称|
|»» name|string|false|none|配置项名称|
|»» desc|string|false|none|配置项的描述|
|»» status|integer(int64)|false|none|使用用户计数（被动异步计算，非精确值）|
|»» release|integer(int64)|false|none|配置项发布（被设置）计数|
|»» channels|[string]|false|none|该配置项适用的产品版本通道，是服务端配置的可用版本通道的子集，为空表示适用所有|
|»» clients|[string]|false|none|该配置项适用的客户端类型，是服务端配置的可用客户端类型的子集，为空表示适用所有|
|»» values|[string]|false|none|该配置项可选值列表，配置项指派给用户或群组时只能从该列表中选择合法值|
|»» createdAt|string(date-time)|false|none|配置项创建时间|
|»» updatedAt|string(date-time)|false|none|配置项更新时间|
|»» offlineAt|string(date-time)|false|none|配置项下线时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 将指定配置项下线，所有设置给用户或群组的对应配置项也会被移除！

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}:offline \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}:offline HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`PUT /v1/products/{product}/modules/{module}/settings/{setting}:offline`

<h3 id="将指定配置项下线，所有设置给用户或群组的对应配置项也会被移除！-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|setting|path|string|true|配置项名称|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="将指定配置项下线，所有设置给用户或群组的对应配置项也会被移除！-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="将指定配置项下线，所有设置给用户或群组的对应配置项也会被移除！-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 批量为用户或群组设置配置项

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}:assign \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}:assign HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`POST /v1/products/{product}/modules/{module}/settings/{setting}:assign`

> Body parameter

```json
{
  "users": [
    "5c4057f0be825b390667abee"
  ],
  "groups": [
    "5bdc1846cd57df001789c751",
    "5bdc1846cd57df001789c751"
  ],
  "value": "beta"
}
```

<h3 id="批量为用户或群组设置配置项-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|setting|path|string|true|配置项名称|
|body|body|[UsersGroupsBody](#schemausersgroupsbody)|true|批量为用户或群组设置灰度标签或配置项的请求数据|
|» users|body|[string]|false|用户 uid 数组，可以不提供，最大长度受 HTTP 请求最大字节数限制|
|» groups|body|[string]|false|群组 uid 数组，可以不提供，最大长度受 HTTP 请求最大字节数限制|
|» value|body|string|false|配置项值，设置灰度标签时不必提供|

> Example responses

> 200 Response

```json
{
  "result": {
    "release": 2,
    "users": [
      "50c32afae8cf1439d35a87e6",
      "5e69a9bd6ac3cd00213ea969"
    ],
    "groups": [
      "string"
    ],
    "value": "x"
  }
}
```

<h3 id="批量为用户或群组设置配置项-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|设置配置项返回结果|Inline|

<h3 id="批量为用户或群组设置配置项-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[SettingReleaseInfo](#schemasettingreleaseinfo)|false|none|none|
|»» release|integer(int64)|false|none|发布批次（被设置）计数|
|»» users|[string]|false|none|用户 uid 数组|
|»» groups|[string]|false|none|群组 uid 数组|
|»» value|string|false|none|配置项值|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 批量撤销对用户或群组设置的产品功能配置项

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}:recall \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}:recall HTTP/1.1
Host: urbs-setting:8443
Content-Type: application/json
Accept: application/json
Authorization: string

```

`POST /v1/products/{product}/modules/{module}/settings/{setting}:recall`

> Body parameter

```json
{
  "release": 2
}
```

<h3 id="批量撤销对用户或群组设置的产品功能配置项-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|setting|path|string|true|配置项名称|
|body|body|[RecallBody](#schemarecallbody)|true|撤销/回滚指定发布批次的灰度标签或配置项|
|» release|body|integer(int64)|false|发布批次（被设置）计数|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="批量撤销对用户或群组设置的产品功能配置项-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="批量撤销对用户或群组设置的产品功能配置项-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 读取指定产品功能配置项的用户列表

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}/users \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}/users HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/products/{product}/modules/{module}/settings/{setting}/users`

<h3 id="读取指定产品功能配置项的用户列表-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|setting|path|string|true|配置项名称|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "result": [
    {
      "hid": "urbs",
      "settingHID": "urbs",
      "kind": "userPercent",
      "rule": "{\"value\": 10}",
      "value": "x",
      "release": 2,
      "createdAt": "2020-05-08T05:38:45.264Z",
      "updatedAt": "2020-05-08T05:38:45.264Z"
    }
  ]
}
```

<h3 id="读取指定产品功能配置项的用户列表-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|配置项作用的用户列表|Inline|

<h3 id="读取指定产品功能配置项的用户列表-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[[SettingUserInfo](#schemasettinguserinfo)]|false|none|none|
|»» hid|string|false|none|发布规则的 hid|
|»» settingHID|string|false|none|发布规则所属配置项的 hid|
|»» kind|string|false|none|发布规则类型，目前仅支持 "userPercent"|
|»» rule|object|false|none|发布规则内容，JSON 对象，具体格式取决于 kind|
|»»» value|integer(int64)|false|none|当 kind 为 "userPercent" 时，value 为百分比，取值 [0, 100]|
|»» value|string|false|none|发布规则的配置项值|
|»» release|integer(int64)|false|none|发布批次（被设置）计数|
|»» createdAt|string(date-time)|false|none|创建时间|
|»» updatedAt|string(date-time)|false|none|更新时间|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 读取指定产品功能配置项的群组列表

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}/groups \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}/groups HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/products/{product}/modules/{module}/settings/{setting}/groups`

<h3 id="读取指定产品功能配置项的群组列表-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|setting|path|string|true|配置项名称|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|
|q|query|string|false|部分列表类 API 支持搜索，q 为搜索关键词|

> Example responses

> 200 Response

```json
{
  "result": [
    {
      "settingHID": "urbs",
      "assignedAt": "2020-05-08T05:38:45.265Z",
      "release": 2,
      "group": "50c32afae8cf1439d35a87e6",
      "kind": "organization",
      "desc": "string",
      "status": 99,
      "value": "b",
      "lastValue": "a"
    }
  ]
}
```

<h3 id="读取指定产品功能配置项的群组列表-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|配置项作用的群组列表|Inline|

<h3 id="读取指定产品功能配置项的群组列表-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|[[SettingGroupInfo](#schemasettinggroupinfo)]|false|none|none|
|»» settingHID|string|false|none|配置项的 hid|
|»» assignedAt|string(date-time)|false|none|配置项设置时间|
|»» release|integer(int64)|false|none|发布批次（被设置）计数|
|»» group|string|false|none|群组的 uid|
|»» kind|string|false|none|群组类型|
|»» desc|string|false|none|群组的描述|
|»» status|integer(int64)|false|none|群组成员数量，非精确值|
|»» value|string|false|none|配置项值|
|»» lastValue|string|false|none|上一个配置项值|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 读取指定产品功能配置项发布规则列表

> Code samples

```shell
# You can also use wget
curl -X GET https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}/rules \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
GET https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}/rules HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`GET /v1/products/{product}/modules/{module}/settings/{setting}/rules`

<h3 id="读取指定产品功能配置项发布规则列表-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|setting|path|string|true|配置项名称|

> Example responses

> 200 Response

```json
{
  "result": []
}
```

<h3 id="读取指定产品功能配置项发布规则列表-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|配置项的发布规则列表|Inline|

<h3 id="读取指定产品功能配置项发布规则列表-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|array|false|none|none|
|»» *anonymous*|any|false|none|none|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 创建指定产品功能配置项的灰度发布规则，同一个功能配置项同一种 kind 的发布规则只能创建一个

> Code samples

```shell
# You can also use wget
curl -X POST https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}/rules \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
POST https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}/rules HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`POST /v1/products/{product}/modules/{module}/settings/{setting}/rules`

<h3 id="创建指定产品功能配置项的灰度发布规则，同一个功能配置项同一种-kind-的发布规则只能创建一个-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|setting|path|string|true|配置项名称|

> Example responses

> 200 Response

```json
{
  "result": null
}
```

<h3 id="创建指定产品功能配置项的灰度发布规则，同一个功能配置项同一种-kind-的发布规则只能创建一个-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|配置项的发布规则结果|Inline|

<h3 id="创建指定产品功能配置项的灰度发布规则，同一个功能配置项同一种-kind-的发布规则只能创建一个-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|any|false|none|none|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 更新指定产品功能配置项的灰度发布规则

> Code samples

```shell
# You can also use wget
curl -X PUT https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}/rules/{hid} \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
PUT https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}/rules/{hid} HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`PUT /v1/products/{product}/modules/{module}/settings/{setting}/rules/{hid}`

<h3 id="更新指定产品功能配置项的灰度发布规则-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|setting|path|string|true|配置项名称|
|uid|path|string|true|标签/配置项/发布规则 hid|

> Example responses

> 200 Response

```json
{
  "result": null
}
```

<h3 id="更新指定产品功能配置项的灰度发布规则-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|配置项的发布规则结果|Inline|

<h3 id="更新指定产品功能配置项的灰度发布规则-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|any|false|none|none|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

## 删除指定产品功能配置项的灰度发布规则

> Code samples

```shell
# You can also use wget
curl -X DELETE https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}/rules/{hid} \
  -H 'Accept: application/json' \
  -H 'Authorization: string'

```

```http
DELETE https://urbs-setting:8443/v1/products/{product}/modules/{module}/settings/{setting}/rules/{hid} HTTP/1.1
Host: urbs-setting:8443
Accept: application/json
Authorization: string

```

`DELETE /v1/products/{product}/modules/{module}/settings/{setting}/rules/{hid}`

<h3 id="删除指定产品功能配置项的灰度发布规则-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|product|path|string|true|产品名称|
|module|path|string|true|功能模块名称|
|setting|path|string|true|配置项名称|
|uid|path|string|true|标签/配置项/发布规则 hid|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="删除指定产品功能配置项的灰度发布规则-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|标准 Boolean 类返回结果|Inline|

<h3 id="删除指定产品功能配置项的灰度发布规则-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» result|boolean|false|none|是否成功|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
HeaderAuthorizationJWT
</aside>

# Schemas

<h2 id="tocSnextpagetoken">NextPageToken</h2>

<a id="schemanextpagetoken"></a>

```json
""

```

*nextPageToken*

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|nextPageToken|string|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|

<h2 id="tocStotalsize">TotalSize</h2>

<a id="schematotalsize"></a>

```json
1

```

*totalSize*

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|totalSize|integer(int64)|false|none|当前分页查询的总数据量|

<h2 id="tocSversion">Version</h2>

<a id="schemaversion"></a>

```json
{
  "name": "urbs-setting",
  "version": "v1.2.0",
  "gitSHA1": "cd7e82a",
  "buildTime": "2020-05-08T05:38:45.267Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|false|none|服务名称|
|version|string|false|none|当前版本|
|gitSHA1|string|false|none|git commit hash|
|buildTime|string(date-time)|false|none|打包构建时间|

<h2 id="tocScachelabelinfo">CacheLabelInfo</h2>

<a id="schemacachelabelinfo"></a>

```json
{
  "l": "beta",
  "cls": [
    "ios",
    "android",
    "web"
  ],
  "chs": [
    "stable",
    "beta",
    "dev"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|l|string|false|none|灰度标签名称|
|cls|[string]|false|none|灰度标签适用的 Clients 客户端类型列表，当列表为空时表示全适用|
|chs|[string]|false|none|灰度标签适用的 Channels 版本通道列表，当列表为空时表示全适用|

<h2 id="tocSlabelinfo">LabelInfo</h2>

<a id="schemalabelinfo"></a>

```json
{
  "hid": "AwAAAAAAAAB25V_QnbhCuRwF",
  "product": "teambition",
  "name": "beta",
  "desc": "string",
  "channels": [
    "beta"
  ],
  "clients": [
    "web"
  ],
  "status": 100,
  "release": 2,
  "createdAt": "2020-05-08T05:38:45.267Z",
  "updatedAt": "2020-05-08T05:38:45.267Z",
  "offlineAt": null
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|hid|string|false|none|灰度标签的 hid|
|product|string|false|none|灰度标签所属的产品名称|
|name|string|false|none|灰度标签名称，同一产品下唯一（不能重名）|
|desc|string|false|none|灰度标签描述|
|channels|[string]|false|none|灰度标签适用版本通道|
|clients|[string]|false|none|灰度标签适用客户端类型|
|status|integer(int64)|false|none|使用用户计数（被动异步计算，非精确值）|
|release|integer(int64)|false|none|标签发布（被设置）计数|
|createdAt|string(date-time)|false|none|灰度标签创建时间|
|updatedAt|string(date-time)|false|none|灰度标签更新时间|
|offlineAt|string(date-time)|false|none|灰度标签下线时间|

<h2 id="tocSmylabel">MyLabel</h2>

<a id="schemamylabel"></a>

```json
{
  "hid": "AwAAAAAAAAB25V_QnbhCuRwF",
  "product": "teambition",
  "name": "task-share",
  "desc": "string",
  "release": 1,
  "assignedAt": "2020-05-08T05:38:45.267Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|hid|string|false|none|灰度标签的 hid|
|product|string|false|none|灰度标签所属的产品名称|
|name|string|false|none|灰度标签名称|
|desc|string|false|none|灰度标签描述|
|release|integer(int64)|false|none|被设置批次|
|assignedAt|string(date-time)|false|none|被设置时间|

<h2 id="tocSmysetting">MySetting</h2>

<a id="schemamysetting"></a>

```json
{
  "hid": "AwAAAAAAAAB25V_QnbhCuRwF",
  "product": "teambition",
  "module": "task",
  "name": "task-share",
  "desc": "string",
  "value": "disable",
  "lastValue": "",
  "release": 1,
  "assignedAt": "2020-05-08T05:38:45.268Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|hid|string|false|none|配置项的 hid|
|product|string|false|none|配置项所属的产品名称|
|module|string|false|none|配置项所属的功能模块名称|
|name|string|false|none|配置项名称|
|desc|string|false|none|配置项描述，|
|value|string|false|none|配置项值|
|lastValue|string|false|none|配置项值|
|release|integer(int64)|false|none|被设置批次|
|assignedAt|string(date-time)|false|none|被设置时间|

<h2 id="tocSuser">User</h2>

<a id="schemauser"></a>

```json
{
  "uid": "5e82d747fe02a50021d339f3",
  "labels": "{\"teambition\":[{\"l\":\"dev\"}]}",
  "activeAt": 1585636012,
  "createdAt": "2020-05-08T05:38:45.268Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|uid|string|false|none|用户的 uid|
|labels|string|false|none|缓存的用户当前被设置的 labels|
|activeAt|integer(int64)|false|none|用户活跃时间戳，1970 以来的秒数，非精确值|
|createdAt|string(date-time)|false|none|用户创建时间|

<h2 id="tocSgroup">Group</h2>

<a id="schemagroup"></a>

```json
{
  "uid": "5e82d747fe02a50021d339f3",
  "kind": "organization",
  "desc": "string",
  "syncAt": 1585636012,
  "status": 99,
  "createdAt": "2020-05-08T05:38:45.268Z",
  "updatedAt": "2020-05-08T05:38:45.268Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|uid|string|false|none|群组的 uid|
|kind|string|false|none|群组类型|
|desc|string|false|none|群组的描述|
|syncAt|integer(int64)|false|none|群组成员同步时间点，1970 以来的秒数|
|status|integer(int64)|false|none|群组成员数量，非精确值|
|createdAt|string(date-time)|false|none|灰度标签创建时间|
|updatedAt|string(date-time)|false|none|灰度标签更新时间|

<h2 id="tocSgroupmember">GroupMember</h2>

<a id="schemagroupmember"></a>

```json
{
  "user": "5e82d747fe02a50021d339f3",
  "syncAt": 1585636012,
  "createdAt": "2020-05-08T05:38:45.268Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|user|string|false|none|群组成员的用户 uid|
|syncAt|integer(int64)|false|none|该群组成员同步时间，1970 以来的秒数|
|createdAt|string(date-time)|false|none|该群组成员添加时间|

<h2 id="tocSproduct">Product</h2>

<a id="schemaproduct"></a>

```json
{
  "name": "urbs",
  "desc": "string",
  "status": 0,
  "createdAt": "2020-05-08T05:38:45.268Z",
  "updatedAt": "2020-05-08T05:38:45.268Z",
  "deletedAt": null,
  "offlineAt": null
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|false|none|产品名称|
|desc|string|false|none|产品的描述|
|status|integer(int64)|false|none|产品状态值|
|createdAt|string(date-time)|false|none|产品创建时间|
|updatedAt|string(date-time)|false|none|产品更新时间|
|deletedAt|string(date-time)|false|none|产品删除时间|
|offlineAt|string(date-time)|false|none|产品下线时间|

<h2 id="tocSproductstatistics">ProductStatistics</h2>

<a id="schemaproductstatistics"></a>

```json
{
  "labels": 0,
  "modules": 0,
  "settings": 0,
  "release": 0,
  "status": 0
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|labels|integer(int64)|false|none|产品下灰度标签数量|
|modules|integer(int64)|false|none|产品下功能模块数量|
|settings|integer(int64)|false|none|产品下配置项数量|
|release|integer(int64)|false|none|产品下灰度标签和配置项总发布次数|
|status|integer(int64)|false|none|产品下灰度标签和配置项总作用人数（非精确值）|

<h2 id="tocSmodule">Module</h2>

<a id="schemamodule"></a>

```json
{
  "name": "urbs",
  "desc": "string",
  "status": 0,
  "createdAt": "2020-05-08T05:38:45.268Z",
  "updatedAt": "2020-05-08T05:38:45.268Z",
  "offlineAt": null
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|false|none|功能模块名称|
|desc|string|false|none|功能模块的描述|
|status|integer(int64)|false|none|有效配置项计数（被动异步计算，非精确值）|
|createdAt|string(date-time)|false|none|功能模块创建时间|
|updatedAt|string(date-time)|false|none|功能模块更新时间|
|offlineAt|string(date-time)|false|none|功能模块下线时间|

<h2 id="tocSsettinginfo">SettingInfo</h2>

<a id="schemasettinginfo"></a>

```json
{
  "hid": "urbs",
  "product": "urbs",
  "module": "urbs",
  "name": "urbs",
  "desc": "string",
  "status": 0,
  "release": 0,
  "channels": [
    "beta"
  ],
  "clients": [
    "ios",
    "android"
  ],
  "values": [
    "true",
    "false"
  ],
  "createdAt": "2020-05-08T05:38:45.268Z",
  "updatedAt": "2020-05-08T05:38:45.268Z",
  "offlineAt": null
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|hid|string|false|none|配置项名称|
|product|string|false|none|配置项名称|
|module|string|false|none|配置项名称|
|name|string|false|none|配置项名称|
|desc|string|false|none|配置项的描述|
|status|integer(int64)|false|none|使用用户计数（被动异步计算，非精确值）|
|release|integer(int64)|false|none|配置项发布（被设置）计数|
|channels|[string]|false|none|该配置项适用的产品版本通道，是服务端配置的可用版本通道的子集，为空表示适用所有|
|clients|[string]|false|none|该配置项适用的客户端类型，是服务端配置的可用客户端类型的子集，为空表示适用所有|
|values|[string]|false|none|该配置项可选值列表，配置项指派给用户或群组时只能从该列表中选择合法值|
|createdAt|string(date-time)|false|none|配置项创建时间|
|updatedAt|string(date-time)|false|none|配置项更新时间|
|offlineAt|string(date-time)|false|none|配置项下线时间|

<h2 id="tocSlabelreleaseinfo">LabelReleaseInfo</h2>

<a id="schemalabelreleaseinfo"></a>

```json
{
  "release": 2,
  "users": [
    "50c32afae8cf1439d35a87e6",
    "5e69a9bd6ac3cd00213ea969"
  ],
  "groups": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|release|integer(int64)|false|none|发布批次（被设置）计数|
|users|[string]|false|none|用户 uid 数组|
|groups|[string]|false|none|群组 uid 数组|

<h2 id="tocSlabelgroupinfo">LabelGroupInfo</h2>

<a id="schemalabelgroupinfo"></a>

```json
{
  "labelHID": "urbs",
  "assignedAt": "2020-05-08T05:38:45.269Z",
  "release": 2,
  "group": "50c32afae8cf1439d35a87e6",
  "kind": "organization",
  "desc": "string",
  "status": 99
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|labelHID|string|false|none|灰度标签 hid|
|assignedAt|string(date-time)|false|none|灰度标签设置时间|
|release|integer(int64)|false|none|发布批次（被设置）计数|
|group|string|false|none|群组的 uid|
|kind|string|false|none|群组类型|
|desc|string|false|none|群组的描述|
|status|integer(int64)|false|none|群组成员数量，非精确值|

<h2 id="tocSlabeluserinfo">LabelUserInfo</h2>

<a id="schemalabeluserinfo"></a>

```json
{
  "labelHID": "urbs",
  "assignedAt": "2020-05-08T05:38:45.269Z",
  "release": 2,
  "user": "50c32afae8cf1439d35a87e6"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|labelHID|string|false|none|灰度标签 hid|
|assignedAt|string(date-time)|false|none|灰度标签设置时间|
|release|integer(int64)|false|none|发布批次（被设置）计数|
|user|string|false|none|用户的 uid|

<h2 id="tocSlabelruleinfo">LabelRuleInfo</h2>

<a id="schemalabelruleinfo"></a>

```json
{
  "hid": "urbs",
  "labelHID": "urbs",
  "kind": "userPercent",
  "rule": "{\"value\": 10}",
  "release": 2,
  "createdAt": "2020-05-08T05:38:45.269Z",
  "updatedAt": "2020-05-08T05:38:45.269Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|hid|string|false|none|发布规则的 hid|
|labelHID|string|false|none|发布规则所属灰度标签的 hid|
|kind|string|false|none|发布规则类型，目前仅支持 "userPercent"|
|rule|object|false|none|发布规则内容，JSON 对象，具体格式取决于 kind|
|» value|integer(int64)|false|none|当 kind 为 "userPercent" 时，value 为百分比，取值 [0, 100]|
|release|integer(int64)|false|none|发布批次（被设置）计数|
|createdAt|string(date-time)|false|none|创建时间|
|updatedAt|string(date-time)|false|none|更新时间|

<h2 id="tocSsettingreleaseinfo">SettingReleaseInfo</h2>

<a id="schemasettingreleaseinfo"></a>

```json
{
  "release": 2,
  "users": [
    "50c32afae8cf1439d35a87e6",
    "5e69a9bd6ac3cd00213ea969"
  ],
  "groups": [
    "string"
  ],
  "value": "x"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|release|integer(int64)|false|none|发布批次（被设置）计数|
|users|[string]|false|none|用户 uid 数组|
|groups|[string]|false|none|群组 uid 数组|
|value|string|false|none|配置项值|

<h2 id="tocSsettinggroupinfo">SettingGroupInfo</h2>

<a id="schemasettinggroupinfo"></a>

```json
{
  "settingHID": "urbs",
  "assignedAt": "2020-05-08T05:38:45.270Z",
  "release": 2,
  "group": "50c32afae8cf1439d35a87e6",
  "kind": "organization",
  "desc": "string",
  "status": 99,
  "value": "b",
  "lastValue": "a"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|settingHID|string|false|none|配置项的 hid|
|assignedAt|string(date-time)|false|none|配置项设置时间|
|release|integer(int64)|false|none|发布批次（被设置）计数|
|group|string|false|none|群组的 uid|
|kind|string|false|none|群组类型|
|desc|string|false|none|群组的描述|
|status|integer(int64)|false|none|群组成员数量，非精确值|
|value|string|false|none|配置项值|
|lastValue|string|false|none|上一个配置项值|

<h2 id="tocSsettinguserinfo">SettingUserInfo</h2>

<a id="schemasettinguserinfo"></a>

```json
{
  "hid": "urbs",
  "settingHID": "urbs",
  "kind": "userPercent",
  "rule": "{\"value\": 10}",
  "value": "x",
  "release": 2,
  "createdAt": "2020-05-08T05:38:45.270Z",
  "updatedAt": "2020-05-08T05:38:45.270Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|hid|string|false|none|发布规则的 hid|
|settingHID|string|false|none|发布规则所属配置项的 hid|
|kind|string|false|none|发布规则类型，目前仅支持 "userPercent"|
|rule|object|false|none|发布规则内容，JSON 对象，具体格式取决于 kind|
|» value|integer(int64)|false|none|当 kind 为 "userPercent" 时，value 为百分比，取值 [0, 100]|
|value|string|false|none|发布规则的配置项值|
|release|integer(int64)|false|none|发布批次（被设置）计数|
|createdAt|string(date-time)|false|none|创建时间|
|updatedAt|string(date-time)|false|none|更新时间|

