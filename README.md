# README

- [README](#readme)
  - [architecture](#architecture)
    - [data flow =\> a flow](#data-flow--a-flow)
    - [monitor flow =\> x, y, z flow](#monitor-flow--x-y-z-flow)
  - [data flow](#data-flow)
  - [api](#api)
  - [service](#service)
  - [Dockerfile](#dockerfile)


## architecture

![architecture](asset/architecture.png)

### data flow => a flow

i.  
kafka 每個 topic 有 3 個 partition.  

ii.  
3 個 app sever 同時為 api server, kafka producer, kafka consumer. 搭配 partition 可以加速發布 message 及消費 message 的速度.  

iii.  
app server 和 kafka 都有多個節點, 具有 ha 效果.  

iv.  
3 個 app server 組成一個 kafka consumer group, 精確分配 kafka partition.

v.  
produce evnet 加上 kafka key, 每筆特定的資料, 只會有單一消費者處理. 當多個消費者處理資料時, 也可以保證 message 的順序.  

### monitor flow => x, y, z flow

i.  
kafka 預期使用 Kafka Exporter 發送 metric 給 prometheus

ii.  
app server 預期使用 第三方客戶端庫, 主動推送 api endpoint 服務狀態

iii.  
grafana 讀取 prometheus, 進行 monitor 數值視覺化.

## data flow

用戶利用 http 觸發 CreateOrder 後, 產生 CreatedOrderEvent  

發布事件到外部系統 shipping service, 利用 kafka 消費事件, 進行 CreateShipping  

簡單一句話說明業務流程:  
`when 訂單被創建, then 進行運輸作業`  

![data_flow](asset/data_flow.png)


## api

```bash
curl -X POST http://localhost:8168/api/v1/orders \
     -H "Content-Type: application/json" \
     -d '{
           "customer_name": "caesar",
           "total_price": 27000
         }'
```

## service

```bash
cd ./deploy && docker compose up -d
```

[docker-compose.yml](deploy/docker-compose.yml)


- app_server:  
    <localhost:8168>
    
- prometheus:  
    <localhost:9090>
    
- grafana:  
    <localhost:3000>  
    user: `root`  
    pw: `1234`
    
-  kafka-ui:  
    <localhost:18080>
      
## Dockerfile

```bash
docker build -f Dockerfile -t x246libra/garmin2024:v0.1.0 . && \
    docker rmi `docker images --filter label=stage=builder -q`
```