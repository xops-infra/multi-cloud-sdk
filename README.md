golang开发的，多云供应商的云资源 sdk的混合，包括 aws,tencent等。

```bash
go get -u github.com/xops-infra/multi-cloud-sdk@main
```


### 支持功能
[cloud.go](pkg/model/cloud.go)

### 开发日志
- 2024-03:
    - feat: remove reigon的初始配置，改为传参后 new client. 这样不需要初始化 N多 client.
    - feat: 支持查询分页
    - feat: add privateDNS for tencent
- 2024-01:
    - feat: 支持 emr list&describe
    - feat: 支持服务器操作
- 2023-12:
    - 新增服务器单查询
- 2023-11:
    - 支持创建解析时候增加info备注
    - 更新服务器过滤filter方法
- 2023-10：
    - 新增 DnsPod 域名解析
- 2023-09：
    - 完成VPC查询整合,新增 OCR 腾讯云接口
    - 新增TIIA图搜图接口
