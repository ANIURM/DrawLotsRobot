# XLab 飞书机器人

| Web Framework | Log Manager     | Config Manager | Api Documentation  | Feishu Api Client     |
|:-------------:|:---------------:|:--------------:|:------------------:|:---------------------:|
| gin-gonic/gin | sirupsen/logrus | spf13/viper    | swaggo/gin-swagger | YasyaKarasu/feishuapi |

## Usage

- `app/controller` 自定义飞书事件处理方法
- `app/global` 全局变量
- `app/dispatcher` 为自定义的 service controller 注册路由
- `config.yaml` 添加自定义的配置字段

## Architecture

- `app` 机器人主体部分
- `config` 机器人配置
- `docs` swagger 生成的 Api 文档
