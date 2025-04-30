# 使用介绍

## 配置介绍

ConfigServer 分为 UI 和 Service 两部分，可以分别独立运行。

### Agent 配置

Agent 侧需要配置 ConfigServer 信息，才能使用管控功能。

#### LoongCollector 配置 ConfigServer

打开 LoongCollector 目录下的 `conf/instance_config/local/loongcollector_config.json` 文件，配置 ConfigServer 相关参数 `ilogtail_configserver_address` 和 `ilogtail_tags`。

`ilogtail_configserver_address` 是 ConfigServer 部署的地址与端口号，可以配置多个 ConfigServer 地址， LoongCollector 将自动切换选择可以链接的 ConfigServer。需要注意的是，目前的 ConfigServer 仅支持单机版，`ilogtail_configserver_address` 即使配置了多个地址，多个 ConfigServer 之间也并不支持数据同步。我们预留了 ConfigServer 支持分布式部署的扩展性，欢迎社区积极贡献开发。

`ilogtail_tags` 是 LoongCollector 在 ConfigServer 处的标签，支持配置多个。虽然该参数暂时无法使用，但我们同样预留了支持通过自定义标签分组管理 Agent 的扩展性。

下面是一个简单的配置示例。

```json
{
    ...
    "ilogtail_configserver_address" : [
      "127.0.0.1:8899"
      ],
    ...
}
```

### Service

Service 为分布式的结构，支持多地部署，负责与采集 Agent 和用户/ui 通信，实现了管控的核心能力。

#### 启动

从 GitHub 下载 LoongCollector 源码，进入 ConfigServer 后端目录下，编译运行代码。

``` bash
cd config_server/service
go build -o ConfigServer
nohup ./ConfigServer > stdout.log 2> stderr.log &
```

#### 配置选项

配置文件为 `config_server/service/seeting/setting.json`。可配置项如下：

* `ip`：Service 服务启动的 ip 地址，默认为 `127.0.0.1`。
* `port`：Service 服务启动的端口，默认为 `8899`。
* `store_mode`：数据持久化的工具，默认为 `leveldb`。当前仅支持基于 leveldb 的单机数据持久化。
* `db_path`：数据持久化的存储地址或数据库连接字符串，默认为 `./DB`。
* `agent_update_interval`：将采集 Agent 上报的数据批量写入存储的时间间隔，单位为秒，默认为 `1` 秒。
* `config_sync_interval`：Service 从存储同步采集 Config 的时间间隔，单位为秒，默认为 `3` 秒。

配置样例：

```json
{
    "ip":"127.0.0.1",
    "store_mode": "leveldb",
    "port": "8899",
    "db_path": "./DB",
    "agent_update_interval": 1,
    "config_sync_interval": 3
}
```

### UI

UI 为一个 Web 可视化界面，与 Service 连接，方便用户对采集 Agent 进行管理。

#### 快速开始

```shell
git clone https://github.com/iLogtail/config-server-ui
cd config-server-ui
yarn install
yarn start
```

#### 更多信息

请参考[这里](https://github.com/iLogtail/config-server-ui/blob/master/README.md)
