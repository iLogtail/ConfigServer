# ConfigServer

## 简介

目前可观测采集 Agent，例如 LoongCollector，主要提供了本地采集配置管理模式，当涉及实例数较多时，需要逐个实例进行配置变更，管理比较复杂。此外，多实例场景下 Agent 的版本信息、运行状态等也缺乏统一的监控。因此，需要提供全局管控服务用于对可观测采集 Agent 的采集配置、版本信息、运行状态等进行统一的管理。

ConfigServer 就是这样的一款可观测 Agent 管控工具，目前支持：

* 采集 Agent 注册到 ConfigServer
* 以 Agent 组的形式对采集 Agent 进行统一管理
* 远程批量配置采集 Agent 的采集配置
* 监控采集 Agent 的运行状态

## 术语表

* 采集 Agent：数据采集器。可以是 LoongCollector，或者其他的采集器。
* 采集配置：一个数据采集 Pipeline，对应一组独立的数据采集配置。
* AgentGroup：可以将相同属性的Agent划分为一个组，只需要绑定采集配置到 AgentGroup 即可在组内所有 Agent 生效。目前仅支持单 AgentGroup（即默认的 AgentGroup `default`）。
* ConfigServer：采集配置管控的服务端。

## 功能描述

对采集 Agent 进行全局管控。任何服从协议的Agent，都可以受到统一的管控。

### 实例注册

* Agent 侧配置已部署的 ConfigServer 信息。
* Agent 启动后，定期向 ConfigServer 进行心跳注册，证明存活性。
* 上报包括但不限于如下的信息，供 ConfigServer 汇集后统一通过API对外呈现。
  * Agent 的 instance_id，作为唯一标识。
  * 版本号
  * 启动时间
  * 运行状态

### AgentGroup 管理

* Agent 管理的基本单元是 AgentGroup ，采集配置通过关联 AgentGroup 生效到组内的 Agent 实例上。
* 系统默认创建 default 组，所有 Agent 也会默认加到 default 组。
* 支持通过 AgentGroup 的 Tag 实现自定义分组。如果 Agent 的 Tag 与 AgentGroup 的 Tag 匹配，则加入该组。
* 一个 AgentGroup 可以包含多个 Agent，一个 Agent 可以属于多个组。

### 采集配置全局管控

* 采集配置通过API进行服务端配置，之后与 AgentGroup 绑定后，即可生效到组内 Agent 实例上。
* 通过采集配置版本号区分采集配置差异，增量变更到 Agent 侧。

### 状态监控

* Agent 定期向 ConfigServer 发送心跳，上报运行信息。
* ConfigServer 汇集后统一通过API对外呈现。

## 快速开始

使用docker，以默认配置快速部署验证 ConfigServer。更多信息请参考[使用介绍](./docs/usage-instructions.md)，[通信协议](./docs/communication-protocol.md)和[开发指南](./docs/config-server-developer-guide.md)。

1. 启动 ConfigServer

    ```shell
      docker run -d \
          --name config-server \
          -p 8899:8899 \
          ghcr.io/ilogtail/config-server:latest
    ```

2. 启动 UI

    ```shell
      docker run -d \
          --name config-server-ui \
          --network="host" \
          -e CONFIG_SERVER_ADDRESS=http://127.0.0.1:8899 \
          ghcr.io/ilogtail/config-server-ui:latest
    ```

3. 启动 LoongCollector

    创建 LoongCollector 配置文件（假设路径为 `/var/lib/loongcollector/conf/instance_config/local/loongcollector_config.json`，集群中可以用cm挂载）

    ```json
      {
        "ilogtail_configserver_address" : [
        "127.0.0.1:8899"
        ]
      }
    ```

    启动 LoongCollector，将配置文件挂载到容器中。（目前镜像仅支持amd64）

    ```shell
      docker run -d \
        --name loongcollector \
        --network="host" \
        -v /:/logtail_host:ro \
        -v /var/run:/var/run \
        -v /var/lib/loongcollector/checkpoint:/usr/local/loongcollector/data/checkpoint \
        -v /var/lib/loongcollector/conf/instance_config/local:/usr/local/loongcollector/conf/instance_config/local \
        ghcr.io/alibaba/loongcollector:edge
    ```

4. 验证

    访问 UI `http://127.0.0.1:80/` ， 即可配置 Agent Group 与采集配置。
