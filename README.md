# Mini Wallet Exercise
## Guideline to Run The App
1. Run Migration (Please configure your database config on MAKEFILE)
```sh
make migrateup
```
2. Configure your database configuration on ./config/app/aplication.dev.yaml
3. Run service
```sh
make start
```
3. Access API on localhost:5000 with preffix 'api/v1'
## API Features 
| Feature | Method | API URL |
| ------ | ------ | ------ |
| Init Wallet | POST | /init |
| View Balance | GET | /wallet |
| View Transactions | GET | /wallet/transactions |
| Enable Wallet | POST | /wallet |
| Disable Wallet | PATCH | /wallet |
| Deposit | POST | /wallet/deposit |
| Withdrawal | POST | /wallet/withdrawals |
