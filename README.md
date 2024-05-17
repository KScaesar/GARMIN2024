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