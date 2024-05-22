# Start

Start with command

```
 fastapi dev api.py
```

# Test cases

Once app is up, run curl command to test

```
curl --location 'http://127.0.0.1:8000/payment' \
--header 'Content-Type: application/json' \
--data '{
    "card": "5555555555554443",
    "expiry_date": "2023-12-23",
    "cvv": "123",
    "amount": "18.8",
    "currency": "CNY",
    "merchant_id": "12333"
}'
```
