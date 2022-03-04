# filetransfer

这是一个负责传输文件的应用程序。

# 目前的功能

通过http请求上传文件到目标linux服务器。

# 第一阶段的目标

### 文件上传

#### 上传任务初始化

POST /file/upload/initialization

**请求头**

|参数     |是否必选|类型|描述|
|:-------:|:-----:|:-----:|:----:|
|Content-Type|是|string|“application/json;charset=utf8” |

**请求体**

|参数     |是否必选|类型|描述|
|:-------:|:-----:|:-----:|:----:|
|resource|是|Object|目标资源信息|
|path|是|string|传输路径，绝对路径|
|filename|是|string|文件名|

resource参数

|参数     |是否必选|类型|描述|
|:-------:|:-----:|:-----:|:----:|
|address|是|string|目标资源地址|
|port|是|number|目标资源端口号|
|account|是|object|登录资源的账号信息|

account参数

|参数     |是否必选|类型|描述|
|:-------:|:-----:|:-----:|:----:|
|name|是|string|目标资源登录用户名|
|password|是|string|登录密码|


**响应体**

**正常响应**

Response 200 OK

|参数     |类型|描述|
|:-------:|:-----:|:----:|
|data|object|正常响应内容|

data参数

|参数     |类型|描述|
|:-------:|:-----:|:----:|
|taskId|string|初始化后的任务id|

**通用异常响应**

Response 400 BadRequest

|参数     |类型|描述|
|:-------:|:-----:|:----:|
|error|object|异常响应内容|

error参数

|参数     |类型|描述|
|:-------:|:-----:|:----:|
|message|string|异常响应错误信息|
|code|string|异常响应错误代码|

其中错误代码使用英文大驼峰缩写。

#### 上传文件

POST /file/upload

**URL参数**

|参数     |是否必选|类型|描述|
|:-------:|:-----:|:-----:|:----:|
|taskId|是|string|任务id|

**请求体**
- 文件流

**响应体**

**正常响应**

Response 204 NoContent

**异常响应**

- 通用异常响应

### 文件下载

#### 文件下载初始化

POST /file/download/initialization

- 请求头与请求体与**上传任务初始化**一致
- 响应与**上传任务初始化**一致

#### 下载文件

GET /file/download

- URL参数与**上传文件**一致

**正常响应**

Response 200 OK

**异常响应**
- 通用异常响应

# Q&A

Q: 功能就这？

A: 啊对对对，还在做嘛，抽空做没时间嘛。

Q: 为啥要做这个？

A: 本人在学习golang，整一个项目练手用。
