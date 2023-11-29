## API Endpoints

Create New Wallet
```bash
$ curl -X POST http://localhost:80/create_wallet -H 'Content-Type: application/json' -d '{"coin":"eth"}'
```

```json
{"id":22,"coin":"eth","balance":0}
```


Process Transaction
```bash
$ curl -X POST http://localhost:80/process_transaction -H 'Content-Type: application/json' -d '{"transaction_type":1, "wallet_id":1, "amount": 1.123}'
```

```json
{"id":43,"transaction_type":1,"wallet_id":1,"amount":1123000000000000000,"updated_balance":1123000000000000000}
```


Transfer Tips
```bash
$ curl -X POST http://localhost:80/transfer_tips -H 'Content-Type: application/json' -d '{"from_wallet_id":1, "to_wallet_id":22, "amount": 0.5}'
```

```json
[{"id":46,"transaction_type":2,"wallet_id":1,"amount":500000000000000000,"updated_balance":3123000000000000000},{"id":47,"transaction_type":1,"wallet_id":22,"amount":500000000000000000,"updated_balance":500000000000000000}]
```

## Usage

run server
```bash
$ make build
$ make run
```

## Docker Usage

```bash
$ make docker-build
$ make docker-run
```