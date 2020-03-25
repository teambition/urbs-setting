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

<h1 id="urbs-setting">urbs-setting v1.2.0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Urbs 灰度平台灰度策略服务

Base URLs:

* <a href="https://urbs-setting:8443">https://urbs-setting:8443</a>

<h1 id="urbs-setting-default">Default</h1>

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
  "db_connect": true
}
```

<h3 id="服务健康检查接口-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|请求成功|[Healthz](#schemahealthz)|

<aside class="success">
This operation does not require authentication
</aside>

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
  "buildTime": "2020-03-25T11:44:33.996Z"
}
```

<h3 id="获取版本信息-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|请求成功|[Version](#schemaversion)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="urbs-setting-user">User</h1>

User 用户相关接口

## 该接口用于灰度网关。获取指定 uid 用户灰度标签在指定 product 产品下的所有（未分页，最多 400 条）灰度标签，包括从 group 群组继承的灰度标签，按照 label 指派时间反序。网关只会取匹配 client 和 channel 的第一条。标签列表不是实时数据，会被服务缓存，缓存时间在 config.cache_label_expire 配置，默认为 1 分钟

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

<h3 id="该接口用于灰度网关。获取指定-uid-用户灰度标签在指定-product-产品下的所有（未分页，最多-400-条）灰度标签，包括从-group-群组继承的灰度标签，按照-label-指派时间反序。网关只会取匹配-client-和-channel-的第一条。标签列表不是实时数据，会被服务缓存，缓存时间在-config.cache_label_expire-配置，默认为-1-分钟-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|uid|path|string|true|用户 UID，当 uid 对应用户不存在时，该接口会返回空灰度标签列表|
|product|query|string|true|产品名称，当 product 对应产品不存在时，该接口会返回空灰度标签列表|

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

<h3 id="该接口用于灰度网关。获取指定-uid-用户灰度标签在指定-product-产品下的所有（未分页，最多-400-条）灰度标签，包括从-group-群组继承的灰度标签，按照-label-指派时间反序。网关只会取匹配-client-和-channel-的第一条。标签列表不是实时数据，会被服务缓存，缓存时间在-config.cache_label_expire-配置，默认为-1-分钟-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|请求成功|Inline|

<h3 id="该接口用于灰度网关。获取指定-uid-用户灰度标签在指定-product-产品下的所有（未分页，最多-400-条）灰度标签，包括从-group-群组继承的灰度标签，按照-label-指派时间反序。网关只会取匹配-client-和-channel-的第一条。标签列表不是实时数据，会被服务缓存，缓存时间在-config.cache_label_expire-配置，默认为-1-分钟-responseschema">Response Schema</h3>

<aside class="success">
This operation does not require authentication
</aside>

## 获取指定 uid 用户所有灰度标签，不包括从群组继承的灰度标签，支持分页。考虑到灰度标签不会很多，暂未支持根据 product 过滤

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

<h3 id="获取指定-uid-用户所有灰度标签，不包括从群组继承的灰度标签，支持分页。考虑到灰度标签不会很多，暂未支持根据-product-过滤-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Authorization|header|string|true|请求 JWT token, 格式如: `Bearer xxx`|
|uid|path|string|true|用户 UID，当 uid 对应用户不存在时，该接口会返回 404 错误|
|pageSize|query|integer(int32)|false|分页大小，默认为 10，(1-1000]|
|pageToken|query|string|false|分页请求标记，来自于响应结果的 nextPageToken|

> Example responses

> 200 Response

```json
{
  "totalSize": 99,
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
      "status": 0,
      "created_at": "2020-03-25T11:44:34.000Z",
      "updated_at": "2020-03-25T11:44:34.000Z",
      "offline_at": null
    }
  ]
}
```

<h3 id="获取指定-uid-用户所有灰度标签，不包括从群组继承的灰度标签，支持分页。考虑到灰度标签不会很多，暂未支持根据-product-过滤-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|请求成功|[LabelsInfoRes](#schemalabelsinfores)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|请求失败|[ErrorResponse](#schemaerrorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|请求失败|[ErrorResponse](#schemaerrorresponse)|

<aside class="success">
This operation does not require authentication
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
|uid|path|string|true|用户 UID，当 uid 对应用户不存在时，该接口会返回 404 错误|

> Example responses

> 200 Response

```json
{
  "result": true
}
```

<h3 id="强制刷新指定用户的灰度标签列表缓存-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|请求成功|[BoolRes](#schemaboolres)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|请求失败|[ErrorResponse](#schemaerrorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|请求失败|[ErrorResponse](#schemaerrorresponse)|

<aside class="success">
This operation does not require authentication
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
99

```

*totalSize*

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|totalSize|integer(int64)|false|none|当前分页查询的总数据量|

<h2 id="tocSerrorresponse">ErrorResponse</h2>

<a id="schemaerrorresponse"></a>

```json
{
  "error": "NotFound",
  "message": "user 50c32afae8cf1439d35a87e6 not found"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|error|string|false|none|错误代号|
|message|string|false|none|错误详情|

<h2 id="tocSboolres">BoolRes</h2>

<a id="schemaboolres"></a>

```json
{
  "result": true
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|result|boolean|false|none|是否成功|

<h2 id="tocSversion">Version</h2>

<a id="schemaversion"></a>

```json
{
  "name": "urbs-setting",
  "version": "v1.2.0",
  "gitSHA1": "cd7e82a",
  "buildTime": "2020-03-25T11:44:34.001Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|false|none|服务名称|
|version|string|false|none|当前版本|
|gitSHA1|string|false|none|git commit hash|
|buildTime|string(date-time)|false|none|打包构建时间|

<h2 id="tocShealthz">Healthz</h2>

<a id="schemahealthz"></a>

```json
{
  "db_connect": true
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|db_connect|boolean|false|none|是否连接了数据库|

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

<h2 id="tocSlabelsinfores">LabelsInfoRes</h2>

<a id="schemalabelsinfores"></a>

```json
{
  "totalSize": 99,
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
      "status": 0,
      "created_at": "2020-03-25T11:44:34.002Z",
      "updated_at": "2020-03-25T11:44:34.002Z",
      "offline_at": null
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|totalSize|[TotalSize](#schematotalsize)|false|none|当前分页查询的总数据量|
|nextPageToken|[NextPageToken](#schemanextpagetoken)|false|none|用于分页查询时用于获取下一页数据的 token，当为空值时表示没有下一页了|
|result|[[LabelInfo](#schemalabelinfo)]|false|none|none|

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
  "status": 0,
  "created_at": "2020-03-25T11:44:34.002Z",
  "updated_at": "2020-03-25T11:44:34.002Z",
  "offline_at": null
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|hid|string|false|none|灰度标签的 HID|
|product|string|false|none|灰度标签所属的产品名称|
|name|string|false|none|灰度标签名称，同一产品下唯一（不能重名）|
|desc|string|false|none|灰度标签描述|
|channels|[string]|false|none|灰度标签适用版本通道|
|clients|[string]|false|none|灰度标签适用客户端类型|
|status|integer(int64)|false|none|灰度标签状态（暂未支持）|
|created_at|string(date-time)|false|none|灰度标签创建时间|
|updated_at|string(date-time)|false|none|灰度标签更新时间|
|offline_at|string(date-time)|false|none|灰度标签下线时间|

