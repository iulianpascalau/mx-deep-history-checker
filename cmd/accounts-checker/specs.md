1. have a defined list of addresses: like this one: `erd103r4tfg6x00jtcyzvara4nwjegrs4mvzmtfvxen3c3688z44d7yqmfs0gj`, `erd105vcnjmzaaw2pd6awpfj6amwcnkm67wljucka7vefnwsk4e62cys4pydst`, `erd100z5n2u5gre3fqzdhnu3twhjjt4m0ym8v8s5t2hg803ujrfhn7dqjs9fnf`, `erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqllls0lczs7`
2. have a const URL that will be hardcoded to `https://mvx-deep-history.jls-software.net/`
3. have a const token that will be hardcoded to `<token>`
4. have a const increment value hardcoded to `5000`
5. have a const of starting nonce hardcoded to `1`
6. query the path: `https://mvx-deep-history.jls-software.net/v1/<token>/network/status/4294967295` and take the erd_nonce value - this will be the highest value until we can check data upon. Response sample:
```json
{
   "data": {
   "status": {
   "erd_block_timestamp": 1775679222,
   "erd_block_timestamp_ms": 1775679222000,
   "erd_cross_check_block_height": "0: 29907817, 1: 29897451, 2: 29903216, ",
   "erd_current_round": 29926937,
   "erd_epoch_number": 2078,
   "erd_highest_final_nonce": 29889808,
   "erd_nonce": 29889808,
   "erd_nonce_at_epoch_start": 29888155,
   "erd_nonces_passed_in_current_epoch": 1653,
   "erd_round_at_epoch_start": 29925284,
   "erd_rounds_passed_in_current_epoch": 1653,
   "erd_rounds_per_epoch": 14400
   }
   },
   "error": "",
   "code": "successful"
}
```
7start querying all defined addresses using this pattern:
`https://mvx-deep-history.jls-software.net/v1/<token>/address/<address>?blockNonce=<nonce>`
The `<nonce>` value will start from the defined starting nonce and incremented by the defined constant.
The status response will be 200, resembling the following:
```json
{
    "data": {
        "account": {
            "address": "erd103r4tfg6x00jtcyzvara4nwjegrs4mvzmtfvxen3c3688z44d7yqmfs0gj",
            "nonce": 3,
            "balance": "241606849315068493",
            "username": "",
            "code": "",
            "codeHash": null,
            "rootHash": null,
            "codeMetadata": null,
            "developerReward": "0",
            "ownerAddress": ""
        },
        "blockInfo": {
            "nonce": 1000000,
            "hash": "0d1a5ee9d8a702cc7d7fd88b9a98984262a5f407525c712b46c9ba7da7c4c06f",
            "rootHash": "2f322cd33ebc27b6107e286984b50eba6fd550ce3a356378c5c311a896b2be2d"
        }
    },
    "error": "",
    "code": "successful"
}
```
The response is valid if one ot the following fields are not default: `rootHash`, `nonce`, `balance`.
If the response is not valid, the program will halt, displaying the address, nonce and the returned fields.
8. the program will compute the time took to check the data, printing some lines once a threshold is reached. 
Let's say we define that once 100000 nonces, we can display, 4 requests/second as an average. 