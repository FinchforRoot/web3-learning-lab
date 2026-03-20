项目的运行环境：
go: 1.25.2
mysql: 5.7.44-log
依赖安装步骤:
切换到项目目录下执行：go mod tidy
启动方式

```
1.修改配置文件的服务器用户名密码以及jwt.secret
2.my-blog-project> go run .\cmd\main.go
```

接口测试：
用户注册:
![img_5.png](img_5.png)
```
curl --location --request POST 'http://127.0.0.1:8910/api/v1/auth/register' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: 127.0.0.1:8910' \
--header 'Connection: keep-alive' \
--data-raw '{
    "username":"user01",
    "email":"12345678910@qq.com",
    "password":"123456"
}'

{
    "code": 200,
    "message": "success",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJ1c2VybmFtZSI6InVzZXIwMSIsImV4cCI6MTc3MzkwODY5NiwibmJmIjoxNzczODIyMjk2LCJpYXQiOjE3NzM4MjIyOTZ9.ORhNhz-N3jcuWzFF0ZAjzW7LojwlFzQd4e12kcMwlpc",
        "user": {
            "id": 4,
            "username": "user01",
            "email": "12345678910@qq.com",
            "created_at": "2026-03-18T16:24:56.307+08:00",
            "updated_at": "2026-03-18T16:24:56.307+08:00"
        }
    }
}
```
用户登录
![img_6.png](img_6.png)
```text
curl --location --request POST 'http://127.0.0.1:8910/api/v1/auth/login' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: 127.0.0.1:8910' \
--header 'Connection: keep-alive' \
--data-raw '{
    "username":"user01",
    "password":"123456"
}'

{
    "code": 200,
    "message": "success",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJ1c2VybmFtZSI6InVzZXIwMSIsImV4cCI6MTc3MzkwODc1OCwibmJmIjoxNzczODIyMzU4LCJpYXQiOjE3NzM4MjIzNTh9.kTJwfHa-7YnG6JxHMBH4xP6vEmE26GKlcD2e8sKYdsc",
        "user": {
            "id": 4,
            "username": "user01",
            "email": "12345678910@qq.com",
            "created_at": "2026-03-18T16:24:56.307+08:00",
            "updated_at": "2026-03-18T16:24:56.307+08:00"
        }
    }
}
```
查看个人基本信息
![img_7.png](img_7.png)
```text
curl --location --request GET 'http://127.0.0.1:8910/api/v1/profile' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJ1c2VybmFtZSI6InVzZXIwMSIsImV4cCI6MTc3MzkwODc1OCwibmJmIjoxNzczODIyMzU4LCJpYXQiOjE3NzM4MjIzNTh9.kTJwfHa-7YnG6JxHMBH4xP6vEmE26GKlcD2e8sKYdsc' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Accept: */*' \
--header 'Host: 127.0.0.1:8910' \
--header 'Connection: keep-alive'

{
    "code": 200,
    "message": "success",
    "data": {
        "id": 4,
        "username": "user01",
        "email": "12345678910@qq.com",
        "created_at": "2026-03-18T16:24:56.307+08:00",
        "updated_at": "2026-03-18T16:24:56.307+08:00"
    }
}
```
以user01用户创建文章，header中配置Authorization，值为Bearer + token
![img_8.png](img_8.png)
```text
curl --location --request POST 'http://127.0.0.1:8910/api/v1/posts' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJ1c2VybmFtZSI6InVzZXIwMSIsImV4cCI6MTc3MzkwODc1OCwibmJmIjoxNzczODIyMzU4LCJpYXQiOjE3NzM4MjIzNTh9.kTJwfHa-7YnG6JxHMBH4xP6vEmE26GKlcD2e8sKYdsc' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: 127.0.0.1:8910' \
--header 'Connection: keep-alive' \
--data-raw '{
    "title":"学好WEB3",
    "content":"这是WEB3"
}'

{
    "code": 200,
    "message": "success",
    "data": {
        "id": 5,
        "user_id": 4,
        "title": "学好WEB3",
        "content": "这是WEB3",
        "created_at": "2026-03-18T16:31:44.995+08:00",
        "updated_at": "2026-03-18T16:31:44.995+08:00",
        "user": {
            "id": 4,
            "username": "user01",
            "email": "12345678910@qq.com",
            "created_at": "2026-03-18T16:24:56.307+08:00",
            "updated_at": "2026-03-18T16:24:56.307+08:00"
        }
    }
}
```
更新文章,先试用user01更新TEST4的文章，返回只有本人可以更新自己的文章。后面将入参的id改为5后，修改成功
![img_9.png](img_9.png)
![img_10.png](img_10.png)
![img_11.png](img_11.png)
```text
curl --location --request PUT 'http://127.0.0.1:8910/api/v1/posts/5' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJ1c2VybmFtZSI6InVzZXIwMSIsImV4cCI6MTc3MzkwODc1OCwibmJmIjoxNzczODIyMzU4LCJpYXQiOjE3NzM4MjIzNTh9.kTJwfHa-7YnG6JxHMBH4xP6vEmE26GKlcD2e8sKYdsc' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: 127.0.0.1:8910' \
--header 'Connection: keep-alive' \
--data-raw '{
    "title":"学习WEB4",
    "content":"这是WEB4"
}'

{
    "code": 200,
    "message": "success",
    "data": {
        "id": 5,
        "user_id": 4,
        "title": "学习WEB4",
        "content": "这是WEB4",
        "created_at": "2026-03-18T16:31:44.995+08:00",
        "updated_at": "2026-03-18T16:33:53.108+08:00",
        "user": {
            "id": 4,
            "username": "user01",
            "email": "12345678910@qq.com",
            "created_at": "2026-03-18T16:24:56.307+08:00",
            "updated_at": "2026-03-18T16:24:56.307+08:00"
        }
    }
}
```
登录用户为文章添加评论，代码中
database.DB.Preload("User").Preload("Post").Preload("Post.User").First(&comment, comment.ID)
这样修改后才会完整返回带有post和post中的user的数据
![img_12.png](img_12.png)
```text
curl --location --request POST 'http://127.0.0.1:8910/api/v1/posts/5/comments' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJ1c2VybmFtZSI6InVzZXIwMSIsImV4cCI6MTc3MzkwODc1OCwibmJmIjoxNzczODIyMzU4LCJpYXQiOjE3NzM4MjIzNTh9.kTJwfHa-7YnG6JxHMBH4xP6vEmE26GKlcD2e8sKYdsc' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: 127.0.0.1:8910' \
--header 'Connection: keep-alive' \
--data-raw '{
    "content":"这篇文章写的很优秀"
}'

{
    "code": 200,
    "message": "success",
    "data": {
        "id": 3,
        "content": "这篇文章写的很优秀",
        "user_id": 4,
        "post_id": 5,
        "created_at": "2026-03-18T16:46:11.182+08:00",
        "updated_at": "2026-03-18T16:46:11.182+08:00",
        "user": {
            "id": 4,
            "username": "user01",
            "email": "12345678910@qq.com",
            "created_at": "2026-03-18T16:24:56.307+08:00",
            "updated_at": "2026-03-18T16:24:56.307+08:00"
        },
        "post": {
            "id": 5,
            "user_id": 4,
            "title": "学习WEB4",
            "content": "这是WEB4",
            "created_at": "2026-03-18T16:31:44.995+08:00",
            "updated_at": "2026-03-18T16:33:53.108+08:00",
            "user": {
                "id": 4,
                "username": "user01",
                "email": "12345678910@qq.com",
                "created_at": "2026-03-18T16:24:56.307+08:00",
                "updated_at": "2026-03-18T16:24:56.307+08:00"
            }
        }
    }
}
```
登录的账户删除文章，只能删除自己的文章
![img_13.png](img_13.png)
![img_14.png](img_14.png)
```text
curl --location --request DELETE 'http://127.0.0.1:8910/api/v1/posts/5' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJ1c2VybmFtZSI6InVzZXIwMSIsImV4cCI6MTc3MzkwODc1OCwibmJmIjoxNzczODIyMzU4LCJpYXQiOjE3NzM4MjIzNTh9.kTJwfHa-7YnG6JxHMBH4xP6vEmE26GKlcD2e8sKYdsc' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Accept: */*' \
--header 'Host: 127.0.0.1:8910' \
--header 'Connection: keep-alive'

{
    "code": 200,
    "message": "success",
    "data": {
        "message": "文章删除成功"
    }
}
```
未登录用户获取文章列表
![img_15.png](img_15.png)
```text
curl --location --request GET 'http://127.0.0.1:8910/api/v1/posts' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Accept: */*' \
--header 'Host: 127.0.0.1:8910' \
--header 'Connection: keep-alive'

{
    "code": 200,
    "message": "success",
    "data": {
        "page": 1,
        "page_size": 10,
        "posts": [
            {
                "id": 4,
                "user_id": 3,
                "title": "学习WEB4",
                "content": "这是WEB4",
                "created_at": "2026-03-18T16:16:33.291+08:00",
                "updated_at": "2026-03-18T16:23:53.102+08:00",
                "user": {
                    "id": 3,
                    "username": "TEST4",
                    "email": "123456789@qq.com",
                    "created_at": "2026-03-18T16:07:45.63+08:00",
                    "updated_at": "2026-03-18T16:07:45.63+08:00"
                }
            },
            {
                "id": 3,
                "user_id": 2,
                "title": "标题3",
                "content": "这是内容3",
                "created_at": "2026-03-18T16:07:21.211+08:00",
                "updated_at": "2026-03-18T16:07:21.211+08:00",
                "user": {
                    "id": 2,
                    "username": "TEST2",
                    "email": "1234526@qq.com",
                    "created_at": "2026-03-18T14:36:37.536+08:00",
                    "updated_at": "2026-03-18T14:36:37.536+08:00"
                }
            },
            {
                "id": 2,
                "user_id": 2,
                "title": "标题2",
                "content": "这是内容2",
                "created_at": "2026-03-18T16:07:13.87+08:00",
                "updated_at": "2026-03-18T16:07:13.87+08:00",
                "user": {
                    "id": 2,
                    "username": "TEST2",
                    "email": "1234526@qq.com",
                    "created_at": "2026-03-18T14:36:37.536+08:00",
                    "updated_at": "2026-03-18T14:36:37.536+08:00"
                }
            },
            {
                "id": 1,
                "user_id": 2,
                "title": "标题1",
                "content": "这是内容",
                "created_at": "2026-03-18T16:07:05.236+08:00",
                "updated_at": "2026-03-18T16:07:05.236+08:00",
                "user": {
                    "id": 2,
                    "username": "TEST2",
                    "email": "1234526@qq.com",
                    "created_at": "2026-03-18T14:36:37.536+08:00",
                    "updated_at": "2026-03-18T14:36:37.536+08:00"
                }
            }
        ],
        "total": 4
    }
}
```
未登录用户特定某文章读取
![img_16.png](img_16.png)
```text
curl --location --request GET 'http://127.0.0.1:8910/api/v1/posts/1' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Accept: */*' \
--header 'Host: 127.0.0.1:8910' \
--header 'Connection: keep-alive'

{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "user_id": 2,
        "title": "标题1",
        "content": "这是内容",
        "created_at": "2026-03-18T16:07:05.236+08:00",
        "updated_at": "2026-03-18T16:07:05.236+08:00",
        "user": {
            "id": 2,
            "username": "TEST2",
            "email": "1234526@qq.com",
            "created_at": "2026-03-18T14:36:37.536+08:00",
            "updated_at": "2026-03-18T14:36:37.536+08:00"
        }
    }
}
```
未登录用户某文章评论读取
![img_17.png](img_17.png)
```text
curl --location --request GET 'http://127.0.0.1:8910/api/v1/comments/post/5' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Accept: */*' \
--header 'Host: 127.0.0.1:8910' \
--header 'Connection: keep-alive'

{
    "code": 200,
    "message": "success",
    "data": {
        "comments": [
            {
                "id": 1,
                "content": "这篇文章写的不错",
                "user_id": 4,
                "post_id": 5,
                "created_at": "2026-03-18T16:41:28.987+08:00",
                "updated_at": "2026-03-18T16:41:28.987+08:00",
                "user": {
                    "id": 4,
                    "username": "user01",
                    "email": "12345678910@qq.com",
                    "created_at": "2026-03-18T16:24:56.307+08:00",
                    "updated_at": "2026-03-18T16:24:56.307+08:00"
                },
                "post": {
                    "id": 5,
                    "user_id": 4,
                    "title": "学习WEB4",
                    "content": "这是WEB4",
                    "created_at": "2026-03-18T16:31:44.995+08:00",
                    "updated_at": "2026-03-18T16:33:53.108+08:00",
                    "user": {
                        "id": 4,
                        "username": "user01",
                        "email": "12345678910@qq.com",
                        "created_at": "2026-03-18T16:24:56.307+08:00",
                        "updated_at": "2026-03-18T16:24:56.307+08:00"
                    }
                }
            },
            {
                "id": 2,
                "content": "这篇文章写的很好",
                "user_id": 4,
                "post_id": 5,
                "created_at": "2026-03-18T16:44:53.174+08:00",
                "updated_at": "2026-03-18T16:44:53.174+08:00",
                "user": {
                    "id": 4,
                    "username": "user01",
                    "email": "12345678910@qq.com",
                    "created_at": "2026-03-18T16:24:56.307+08:00",
                    "updated_at": "2026-03-18T16:24:56.307+08:00"
                },
                "post": {
                    "id": 5,
                    "user_id": 4,
                    "title": "学习WEB4",
                    "content": "这是WEB4",
                    "created_at": "2026-03-18T16:31:44.995+08:00",
                    "updated_at": "2026-03-18T16:33:53.108+08:00",
                    "user": {
                        "id": 4,
                        "username": "user01",
                        "email": "12345678910@qq.com",
                        "created_at": "2026-03-18T16:24:56.307+08:00",
                        "updated_at": "2026-03-18T16:24:56.307+08:00"
                    }
                }
            },
            {
                "id": 3,
                "content": "这篇文章写的很优秀",
                "user_id": 4,
                "post_id": 5,
                "created_at": "2026-03-18T16:46:11.182+08:00",
                "updated_at": "2026-03-18T16:46:11.182+08:00",
                "user": {
                    "id": 4,
                    "username": "user01",
                    "email": "12345678910@qq.com",
                    "created_at": "2026-03-18T16:24:56.307+08:00",
                    "updated_at": "2026-03-18T16:24:56.307+08:00"
                },
                "post": {
                    "id": 5,
                    "user_id": 4,
                    "title": "学习WEB4",
                    "content": "这是WEB4",
                    "created_at": "2026-03-18T16:31:44.995+08:00",
                    "updated_at": "2026-03-18T16:33:53.108+08:00",
                    "user": {
                        "id": 4,
                        "username": "user01",
                        "email": "12345678910@qq.com",
                        "created_at": "2026-03-18T16:24:56.307+08:00",
                        "updated_at": "2026-03-18T16:24:56.307+08:00"
                    }
                }
            }
        ],
        "page": 1,
        "page_size": 20,
        "total": 3
    }
}
```