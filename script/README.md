# Script transfer random

### Gen key và tạo file
    
```
./tool-transfer transfer  gen-key  -p <private key> --rpc  <rpc> --receiver-number <amount key>  -f <file-name>
```
- example: 

```./tool-transfer transfer  gen-key  -p bb900295e2a1a6a84bbf6747e89f60602a0677970a391a4087c5c6b81280aab7 --rpc https://l1testnet.trustkeys.network --receiver-number 1000 -f key ```
    

### Tạo file và transfer giữa các key trong file
```
 ./tool-transfer transfer random-transfer --prvkey <private key> --rpc <rpc> --receiver-number <amount key> --max-value-per-transaction <value per transaction> --time-cross-transfer < second> time-between-trans <mili second> -f <fileName >
```
- file được lưu trong folder tmp
- random user và value mỗi giao dịch

- example: 

```./tool-transfer transfer random-transfer  --prvkey bb900295e2a1a6a84bbf6747e89f60602a0677970a391a4087c5c6b81280aab7 --rpc https://l1testnet.trustkeys.network --receiver-number 20 --max-value-per-transaction 20 --time-cross-transfer 36000 time-between-trans 1000 -f key```



### Import file và transfer giữa các key trong file
```
./tool-transfer transfer import-file-and-transfer  --prvkey  <private key> --rpc <rpc> --time-cross-transfer <second> time-between-trans <mili second> -f <fileName >
```
- example: 

```./tool-transfer transfer import-file-and-transfer  --prvkey bb900295e2a1a6a84bbf6747e89f60602a0677970a391a4087c5c6b81280aab7 --rpc https://l1testnet.trustkeys.network  --time-cross-transfer 36000 time-between-trans 1000 -f key```