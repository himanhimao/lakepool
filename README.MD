# LAKE POOL
> 这是一个基于golang的开源矿池， 包含[后端](#后端)，前端，中台 三个部分

### 后端

#### 设计初衷
- stratum协议本身是一个消息同步协议，其架构应该尽量同步化
- 矿池支持多个币种本质都是基于stratum协议，其架构因保持一定的复用性，如日志模块，share验证模块， 同时也应该保持一定的独立性，做好资源隔离，例如stratum_server
- 消灭大量无效share_log日志, 减轻后端压力
- 后端只负责提供数据，而不直接负责业务的结算，数据的业务结算由中台处理
- 采用开发效率的go语言实现

#### 架构图
![backend architecture](docs/images/lakepool-backend-arch.png)



#### DOCKER BUILD
> log-server 
```
cd lakepool/backend/log_server
docker build -t lakepool/backend/log_server -f build/Dockerfile .
```
> sphere-server 
```
cd lakepool/backend/sphere_server
docker build -t lakepool/backend/sphere_server -f build/Dockerfile .
```
> stats-server 
```
cd lakepool/backend/stats_server
docker build -t lakepool/backend/stats_server -f build/Dockerfile .
```
> stratum-server
```
cd lakepool/backend/stratum_server
docker build -t lakepool/backend/stratum_server -f build/Dockerfile .
```
#### 中台


#### 前端

