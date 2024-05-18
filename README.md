# README

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