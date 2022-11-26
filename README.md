# wike-service

About wiki project  backend for front & admin service

```
.
├── application         // 应用层
│   ├── dto                 // Data Transfer Object 数据传输对象
│   ├── event               // 后台事件
│   └── service             // 对外应用服务
│       ├── grpc                // GRPC 服务
│       └── http                // HTTP 
├── config              // 配置中心
│   ├── config.go
│   └── files
│       └── config.toml.go
├── domain              // 领域服务层
│   └── game                // 领域
│       ├── entity              // 实体
│       │   ├── object_value.go // 值对象
│       │   └── do.go               // Domain Object 领域对象 
│       ├── repository  // 仓储层
│       └── service.go // 领域服务对外API
├── global // 公共
│   ├── init.go // 文件
│   └── router.go // http服务路由
├── main.go // 项目启动文件
├── Makefile //Makefile
└── README.md
```
