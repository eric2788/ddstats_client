# DDStats Client
贡献直播间专用端

## 部署

### Docker 

详见Dockerfile

### 手动

到 Releases 下载对应的平台执行程序


## 使用方式

打开后，到网址栏输入 `http://localhost:9090/` 即可开始设置

## 直播间贡献原理

提交订阅后，将在 `https://blive.ericlamm.xyz` 那边启动监听，从而让 高亮用户统计化网站 获得数据。

## 注意

你需要长期开着本程序而维持该直播间的监听。

关闭程序五分钟后会自动视为放弃监听而清空监控列表。但是，你可以在界面透过按下离线储存按钮来储存正在监听的房间列表，这样下次打开程序就能预先订阅先前离线的房间。
