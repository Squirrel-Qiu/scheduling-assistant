# schedule

### 前端 post /schedule/login 登录成功后获取个人创建的值班表信息

```json
{
    "openid": "openid",
    "nick_name": "张三"
}
```

### 后端 return

```json
{
    "rotas": [
        {
            "rota_id": "291255583271555074",
            "title": "IT协会值班表",
            "shift": 2,
            "limit_choose": 4,
            "counter": 6
        },
        {
            "rota_id": "291255583271555075",
            "title": "人事部",
            "shift": 4,
            "limit_choose": 5,
            "counter": 8
        },
        {
            "rota_id": "291255583271555076",
            "title": "xxx值班表",
            "shift": 1,
            "limit_choose": 2,
            "counter": 4
        }
    ]
}
```

If login success, will return json with token. Otherwise will return http BadRequest

### 前端 post /schedule/newRota

```json
{
    "title": "xxx值班表",
    "shift": 2,
    "limit_choose": 4,
    "counter": 6
}
```

### 后端 return

```json
{
    "status": 1,
    "result": "not defined(滑稽)"
}
```

token is UUID

status: OK is 0, failed is 1, jsonErr is 2, requestErr is 3

if status is OK, result is null

if status is failed, means the student is signed up, result is null

if status is jsonErr, means the json query has some format error, result is null

if status is requestErr, means the query has some problem, result is 

```json
{
    "field": "the problem field",
    "reason": "the reason"
}
```
//////////////////////
### 前端 查看某一值班表 get /schedule/detail?rota_id=291255583271555074

### 后端 return

```json
{
    "frees": [0,1,5,6]
}
```

status: OK is 0, timeout is 1, student not found is 2, jsonErr or url query error is 3