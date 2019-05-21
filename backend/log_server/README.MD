## log_server
> 公共的日志服务，stats_server 和 stratum_server 产生的share_log, block_log(出块日志)，都会再此进行集中化处理. share_log 会再次聚合
合并不同stratum_server产生的算力

### 主要模块

#### SERVER
> 接口服务器，为上游服务器提供日志服务(查询，写入)

#### PUSHER (TODO)
> 将集群的业务数据，发往数据中台

#### 数据库
> 数据库采用时序数据库influxdb

#### 数据字典
- min_share_log
    - tags
        - worker_name
        - is_right
    - fields
        - compute_power
        - server_ip
        - client_ip
        - user_name
        - ext_name
        - host_name
        - user_agent
        - pid 
        - height 
- block_log
    - fields
        - server_ip
        - client_ip
        - user_name
        - user_agent
        - ext_name
        - host_name
        - user_agent
        - pid 
        - worker_name
        - hash 
        
### 聚合维度
> worker_name, is_right, 1min 

### 前期准备
- 创建数据库

`CREATE DATABASE "mining_share"`

`CREATE DATABASE "mining_log"`

- 创建连续查询
> 多个币种需更具不同的measurement，建立不同的连续查询

`CREATE CONTINUOUS QUERY "cq_log_btc_share_final" ON "mining_share" BEGIN SELECT count("host_name") as count, mean("compute_power") as compute_power, last("server_ip") as server_ip, last("client_ip") as client_ip, last("user_name") as user_name, last("ext_name") as ext_name, last("host_name") as host_name, last("user_agent") as user_agent, last("host_name") as host_name, last("pid") as pid, last("height") as height  INTO "mining_share"."autogen"."log_share_btc_final" FROM "mining_share"."autogen"."log_share_btc" GROUP BY time(1m), worker_name, is_right fill(none) END`

### 结果验证
`SELECT * FROM "mining_share"."autogen"."log_share_btc_final"`


## 部署
> docker build 
```
cd lakepool/backend/log_server
docker build -t lakepool/backend_log_server -f build/Dockerfile .