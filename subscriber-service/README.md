```
curl -X POST http://localhost:8080/subscribe \
-H "Content-Type: application/json" \
-d '{
  "tickers": [
    {
      "symbol": "BTC/USD",
      "changeThreshold": 6.0
    },
    {
      "symbol": "ETH/USD",
      "changeThreshold": 6.0
    }
  ],
  "notificationOptions": {
    "email": "user@example.com"
  }
}'

```