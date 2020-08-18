# schedule

### 前端 POST /login

```json
{
    "appid": "******************",
    "secret": "********************************",
    "js_code": "*********************************",
    "grant_type": "authorization_code"
}
```

### 后端 return

if login success, will get cookie
```json
{
    "status": 0
}
```
```json
{
    "status": 1,
    "msg": "请求参数错误"
}
```
```json
{
    "status": 2,
    "msg": "用户认证失败"
}
```
status = 3: 服务器内部错误
```json
{
    "status": 3
}
```

### 前端 POST /savePerson 登录成功后存储用户昵称(携带cookie)

```json
{
    "nick_name": "张三"
}
```

### 后端 return

```json
{
    "status": 0
}
```
```json
{
    "status": 1,
    "msg": "请求参数错误"
}
```
status = 2: 服务器内部错误
```json
{
    "status": 2
}
```

### 前端 POST /newRota 创建值班表(携带cookie)

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
    "status": 0,
    "rota_id": "291255583271555074"
}
```
```json
{
    "status": 1,
    "msg": "请求参数错误"
}
```
```json
{
    "status": 2,
    "msg": "限选不得小于班次"
}
```
status = 3: 服务器内部错误
```json
{
    "status": 3,
    "msg": "实例化工作节点错误"
}
```
```json
{
    "status": 3
}
```

### 前端 GET /rotas 获取用户创建的所有值班表(携带cookie)

### 后端 return

```json
{
    "status": 0,
    "rotas": [
        {
            "rota_id": "291255583271555074",
            "title": "一月份值班表",
            "shift": 2,
            "limit_choose": 4,
            "counter": 6
        },
        {
            "rota_id": "291255583271555075",
            "title": "人事部值班表",
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
status = 1: 服务器内部错误
```json
{
    "status": 1
}
```

### 前端 GET /join 获取用户参与的所有值班表(携带cookie)

### 后端 return

```json
{
    "status": 0,
    "joins": [
        {
            "rota_id": "291255583271555074",
            "title": "一月份值班表"
        },
        {
            "rota_id": "291255583271555075",
            "title": "人事部值班表"
        }
    ]
}
```
status = 1: 服务器内部错误
```json
{
    "status": 1
}
```

### 前端 GET /rota/rota_id=291255583271555074 查看某一值班表(携带cookie)

### 后端 return

```json
{
    "status": 0,
    "frees": [0,1,5,6]
}
```
```json
{
    "status": 1,
    "msg": "rotaId错误"
}
```
status = 2: 服务器内部错误
```json
{
    "status": 2
}
```

### 前端 POST /chooseFree/rota_id=291255583271555074 在该值班表下选择时间段(携带cookie)

```json
{
    "frees": [0,1,5,6]
}
```

### 后端 return

```json
{
    "status": 0

}
```
```json
{
    "status": 1,
    "msg": "rotaId错误"
}
```
```json
{
    "status": 2,
    "msg": "请求参数错误"
}
```
```json
{
    "status": 3,
    "msg": "请至少选择<limitChoose>个时间段"
}
```
status = 4: 服务器内部错误
```json
{
    "status": 4
}
```

### 前端 GET /generate/rota_id=291255583271555074 生成值班表

### 后端 return

```json
{
    "status": 0,
    "interval": [
            {
                "free_id": 3,
                "members": ["张三", "李四", "王明"]
            },
            {
                "rota_id": 14,
                "members": ["张三"]
            },
            {
                "rota_id": 21,
                "members": ["李四", "王明"]
            }
        ]
}
```
```json
{
    "status": 1,
    "msg": "rotaId错误"
}
```
```json
{
    "status": 2,
    "msg": "值班表为空"
}
```
status = 3: 服务器内部错误
```json
{
    "status": 3
}
```

### 前端 GET /download/rota_id=291255583271555074 导出值班表(值班表.xlsx)

### 后端 return

"Content-Type", "application/octet-stream; charset=utf-8"
```json
{
    "status": 1,
    "msg": "值班表文件不存在"
}
```

### 前端 DELETE /delete/rota_id=291255583271555074 删除值班表(携带cookie)

### 后端 return
status = 0: delete success
```json
{
    "status": 0
}
```
```json
{
    "status": 1,
    "msg": "rotaId错误"
}
```
status = 2: 服务器内部错误
```json
{
    "status": 2
}
```