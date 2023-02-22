# simple-demo

## 抖音项目服务端简单示例

具体功能内容参考飞书说明文档

工程无其他依赖，直接编译运行即可

```shell
go build && ./fastdouyin[fastdouyin-shliang.zip](..%2Ffastdouyin-shliang.zip)
```

### 功能说明

实现了基础接口全部、互动接口全部、社交接口的关系操作、关注列表以及粉丝列表。

* 用户登录数据，点赞关注等信息都会保存在Mysql中，点赞关注等操作会同时缓存到Redis中，来优化读取速度。
* 视频上传后会保存到本地 public\videos 目录中，视频第一帧截图作为封面会保存到本地 public\covers 目录中，头像背景墙保存在本地public的avatar和background-image中。

### 测试

test 目录下为不同场景的功能测试case，可用于验证功能实现正确性

其中 common.go 中的 _serverAddr_ 为服务部署的地址，默认为本机地址，可以根据实际情况修改

测试数据写在 demo_data.go 中，用于列表接口的 mock 测试