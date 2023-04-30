# XLab 抽签机器人

| Log Manager     | Config Manager | Api Documentation  | Feishu Api Client     |
|:---------------:|:--------------:|:------------------:|:---------------------:|
| sirupsen/logrus | spf13/viper    | swaggo/gin-swagger | YasyaKarasu/feishuapi |

## Usage

- `app/controller` 自定义飞书事件处理方法
- `app/global` 全局变量
- `app/dispatcher` 为自定义的 service controller 注册路由
- `config.yaml` 添加自定义的配置字段

## Architecture

- `app` 机器人主体部分
- `config` 机器人配置
- `docs` swagger 生成的 Api 文档

## 使用说明
抽签机器人使用说明：
1. @机器人，启动抽签
2. 输入@机器人 所有人 或者 @机器人 all，抽取所有人。输入@机器人 @xxx @xxx，抽取@的人
3. 输入@机器人 组数
4. 输入@机器人 每组人数
   即可获得抽签结果
 
    输入@机器人 reset 或者 @机器人 重置，重置抽签机器人
 
    输入@机器人 help 或者 @机器人 帮助，查看抽签机器人使用说明

## 部署
在`conf.yaml`填好相关的配置
```text
feishu:
  #   该区域请于飞书开放平台查询本机器人信息,详见
  #   https://open.feishu.cn/document/home/develop-a-bot-in-5-minutes/coding
  appId: 
  appSecret: 
  verificationToken: 
  encryptKey: 
  larkHost: "https://open.feishu.cn"

server:
  # 端口，请根据实际情况为准
  port: 10001
```