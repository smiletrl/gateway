# Start app locally

```
 export GOWORK=off
 make app
```

This app uses port `1323`. In case your local computer has port `1323` used already, change it at `service.payment/cmd/api.go`, line 38, and `docker-compose.yml`, line 10.

# Run test locally

```
make test
```

# Check code

```
make lintcheck
```

# Run as docker container

```
make docker
```

# Test cases

Once app is up, run curl command to test

```
curl --location 'http://localhost:1323/payment' \
--header 'Content-Type: application/json' \
--data '{
    "card": "5555555555554444",
    "expiry_date": "2023-12-23",
    "cvv": "123",
    "amount": "18.89",
    "currency": "CNY",
    "merchant_id": "12333"
}'
```
