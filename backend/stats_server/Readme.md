## stats_server
> 聚合器，和单个stratum_server搭配使用，负责统计单个stratum_server的收益情况


### 主要模块

#### server
> 接口服务器， stratum_server 会将日志数据通过接口上传至stats_server, 并持久化存储

#### pusher 
> pusher进程，定期将聚合后的统计数据，传递到log_server


### 数据库
> 数据库采用时序数据库influxdb

### 数据字典
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

### 聚合维度
> worker_name, is_right, 1min 


### 前期准备
- 创建数据库

` CREATE DATABASE "mining_stats"`

- 创建保留策略

`CREATE RETENTION POLICY "three_days" on "mining_stats" DURATION 72h  REPLICATION 1 DEFAULT`

`CREATE RETENTION POLICY "two_weeks" on "mining_stats" DURATION 336h  REPLICATION 1`
- 创建连续查询

`CREATE CONTINUOUS QUERY "cq_stats_btc_share_1min" ON "mining_stats" BEGIN SELECT count("host_name") as count, mean("compute_power") as compute_power, last("server_ip") as server_ip, last("client_ip") as client_ip, last("user_name") as user_name, last("ext_name") as ext_name, last("host_name") as host_name, last("user_agent") as user_agent, last("host_name") as host_name, last("pid") as pid, last("height") as height  INTO "mining_stats"."two_weeks"."stats_share_btc_1min" FROM "mining_stats"."three_days"."stats_share_btc" GROUP BY time(1m), worker_name, is_right fill(none) END`

### 结果验证
- 查询结果
`SELECT * FROM select * from mining_stats.two_weeks.stats_share_btc_1min`