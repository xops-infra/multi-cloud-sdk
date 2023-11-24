golang开发的，多云供应商的云资源 sdk的混合，包括 aws,tencent等。

```bash
go get -u github.com/xops-infra/multi-cloud-sdk
```

### 支持功能
- 服务器资源
    - 服务器
- 虚拟网络
    - VPC
    - NAT
    - EIP
    - SUBNET
- 对象存储
- EMR
- 消息队列
- OCR(tencent)
- TIIA 图搜图(tencent)
- DnsPod 域名解析(aws,tencent)


### 开发日志
- 2023-11:
    - 支持创建解析时候增加info备注
    - 更新服务器过滤filter方法
- 2023-10：
    - 新增 DnsPod 域名解析
- 2023-09：
    - 完成VPC查询整合,新增 OCR 腾讯云接口
    - 新增TIIA图搜图接口
