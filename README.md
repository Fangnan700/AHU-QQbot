# AHU-QQbot



最近在学习Golang，就把学校的教务系统接口收集了一下，集成到了这个QQ机器人中，实现在QQ上操作教务系统。

目前实现的功能：

1. 用户登录
2. 查询课表
3. 查询考试



使用指南：

1. clone本项目或从`release`页面下载最新版可执行文件
2. 授权：`chmod +x ahu-qqbot`
3. 运行程序，初始化配置文件并填写
4. 启动



配置文件示例：

```yaml
GptHost: ""
GptProxy: ""
GptModel: ""
GptKeys: 
- "xxxxxxxxxxxxxxxxxxxxxx"
RedisHost: "172.0.0.1:6379"
RedisPass: "123456"
CqHttpHost: "http://127.0.0.1:5700"
CqHttpPath: "/home/cqhttp"
AhuCalendarStartDate: "2023-03-06"
```

