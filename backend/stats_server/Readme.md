## influxdb 
CREATE DATABASE "mining_stats"

### RETENTION POLICY
CREATE RETENTION POLICY "three_days" on "mining_stats" DURATION 72h  REPLICATION 1 DEFAULT
CREATE RETENTION POLICY "two_weeks" on "mining_stats" DURATION 336h  REPLICATION 1


### CONTINUOUS QUERIES
CREATE CONTINUOUS QUERY "cq_stats_btc_share_1min" ON "mining_stats" BEGIN SELECT mean("compute_power") as compute_power INTO "mining_stats"."two_weeks"."stats_share_btc_1min" FROM "mining_stats"."three_days"."stats_share_btc" GROUP BY time(1m), * fill(none) END



SELECT mean("compute") as compute FROM "stats_share_btc" GROUP BY time(1m), * fill(none)


SELECT mean("compute") as compute INTO "mining_stats"."two_weeks"."stats_share_btc_1min" FROM "mining_stats"."three_days"."stats_share_btc" GROUP BY time(1m), * fill(none) 


select count(*) from "mining_stats"."two_weeks"."stats_share_btc_1min"
select count(*) from "mining_stats"."three_days"."stats_share_btc"

select * from "mining_stats"."two_weeks"."stats_share_btc_1min" where "host_name"='windows2008' and "pid"='6' and "time">0m limit 200  offset 0

### 获取分钟内用户总算力
select sum(compute) from "mining_stats"."two_weeks"."stats_share_btc_1min" group by "user_name"
SELECT last("compute_power") FROM "two_weeks".stats_share_btc_1min GROUP BY "hostname","pid"



