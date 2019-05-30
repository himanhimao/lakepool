## sphere_server
> 负责为stratum_server,提供各币种的Job推送，share验证，出块提交的工作.
> 目前暂时只支持币种BTC

### 数据库
> 依赖内存数据库redis
- 建议配置(单个币种), 可根据业务配置自行调整
    - cpu 1core
    - memory 4G
- 需修改redis内存回收策略为volatile-random, redis默认不会做内存回收

### 数据节点
> 依赖币种的原始节点