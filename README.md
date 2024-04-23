# Start app locally

```
 export GOWORK=off
 make app
```

This app uses port `1323`. In case your local computer has port `1323` used already, change it at `service.payment/cmd/api.go`, line 38.

# Run test locally

```
make test
```
