golang 开发的，多云供应商的云资源 sdk 的混合，包括 aws,tencent 等。

```bash
go get -u github.com/xops-infra/multi-cloud-sdk@main
```

### 说明

由于时间问题，很多借口没有完全实现，只实现了业务用到的部分，如果你需要用到其他接口，欢迎 PR，共建。

### 支持功能

[model_service.go](pkg/model/model_service.go) 封装的 service 接口，可以屏蔽云差异；
[model_io.go](pkg/model/model_io.go) 封装的各个云 io 接口，如果你只用到某个云的话，可以直接用这里的借口；

### 开发日志

- 2024-11:
  - feat: add 对象存储生命周期管理，初步调试几个接口。
- 2024-04:
  - feat: add s3&cos 比原版 SDK 增强，支持 tag 返回以及 s3 location 返回。
  - feat: add s3&cos presign url 生成，支持对象不存在报错。
- 2024-03:
  - feat: remove reigon 的初始配置，改为传参后 new client. 这样不需要初始化 N 多 client.
  - feat: 支持查询分页
  - feat: add privateDNS for tencent
- 2024-01:
  - feat: 支持 emr list&describe
  - feat: 支持服务器操作
- 2023-12:
  - 新增服务器单查询
- 2023-11:
  - 支持创建解析时候增加 info 备注
  - 更新服务器过滤 filter 方法
- 2023-10：
  - 新增 DnsPod 域名解析
- 2023-09：
  - 完成 VPC 查询整合,新增 OCR 腾讯云接口
  - 新增 TIIA 图搜图接口
