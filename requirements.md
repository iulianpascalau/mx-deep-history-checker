Ok, I'm back. We will need to develop an app that will check (as much as it can) the integrity of the MultiversX nodes level-DB structures.
There are 4 shards (data partitions) for the MultiversX network. Each partition has more or less the same configuration: 
shard 0, shard 1 and shard 2 having the exact same directory structure, shard metachain being different. 

├── └──  │

```
node-0
   └── 1
       │
       ├── Epoch_0
       │      └── Shard_0
       │            ├── AccountsTrie 
       │            ├── AccountsTrieCheckpoints
       │            ├── BlockHeaders
       │            ├── BootstrapData
       │            ├── DbLookupExtensions
       │            │         └── MiniblocksMetadata
       │            ├── DbLookupExtensions_ResultsHashesByTx
       │            ├── Logs
       │            ├── MetaBlock
       │            ├── MiniBlocks
       │            ├── PeerAccountsTrie
       │            ├── PeerAccountsTrieCheckpoints
       │            ├── PeerBlocks
       │            ├── Receipts
       │            ├── RewardTransactions
       │            ├── ScheduledSCRs
       │            ├── Transactions
       │            ├── TxLookupExtensions
       │            └── UnsignedTransactions
       ├── Epoch_1
       │      └── Shard_0
       │            ├── AccountsTrie 
       │            ├── AccountsTrieCheckpoints
       │            ├── BlockHeaders
       │            ├── BootstrapData
       │            ├── DbLookupExtensions
       │            │         └── MiniblocksMetadata
       │            ├── DbLookupExtensions_ResultsHashesByTx
       │            ├── Logs
       │            ├── MetaBlock
       │            ├── MiniBlocks
       │            ├── PeerAccountsTrie
       │            ├── PeerAccountsTrieCheckpoints
       │            ├── PeerBlocks
       │            ├── Receipts
       │            ├── RewardTransactions
       │            ├── ScheduledSCRs
       │            ├── Transactions
       │            └── UnsignedTransactions
       .......
       ├── Epoch_999
       │      └── Shard_0
       │            ├── AccountsTrie 
       │            ├── AccountsTrieCheckpoints
       │            ├── BlockHeaders
       │            ├── BootstrapData
       │            ├── DbLookupExtensions
       │            │         └── MiniblocksMetadata
       │            ├── DbLookupExtensions_ResultsHashesByTx
       │            ├── Logs
       │            ├── MetaBlock
       │            ├── MiniBlocks
       │            ├── PeerAccountsTrie
       │            ├── PeerAccountsTrieCheckpoints
       │            ├── PeerBlocks
       │            ├── Receipts
       │            ├── RewardTransactions
       │            ├── ScheduledSCRs
       │            ├── Transactions
       │            └── UnsignedTransactions
       .......
       └── Static
              └── Shard_0
                    ├── DbLookupExtensions_EpochByHash
                    ├── DbLookupExtensions_ESDTSupplies
                    ├── DbLookupExtensions_MiniblockHashByTxHash
                    ├── DbLookupExtensions_RoundHash
                    ├── MetaHdrHashNonce
                    ├── ShardHdrHashNonce0
                    └── StatusMetricsStorageDB
node-1 (same as node 0)
node-2 (same as node 0)
node-m 
   └── 1
       │
       ├── Epoch_0
       │      └── Shard_metachain
       │            ├── AccountsTrie 
       │            ├── AccountsTrieCheckpoints
       │            ├── BlockHeaders
       │            ├── BootstrapData
       │            ├── DbLookupExtensions
       │            │         └── MiniblocksMetadata
       │            ├── DbLookupExtensions_ResultsHashesByTx
       │            ├── Logs
       │            ├── MetaBlock
       │            ├── MiniBlocks
       │            ├── PeerAccountsTrie
       │            ├── PeerAccountsTrieCheckpoints
       │            ├── PeerBlocks
       │            ├── Receipts
       │            ├── RewardTransactions
       │            ├── ScheduledSCRs
       │            ├── Transactions
       │            └── UnsignedTransactions
       ├── Epoch_1
       │      └── Shard_metachain
       │            ├── AccountsTrie 
       │            ├── AccountsTrieCheckpoints
       │            ├── BlockHeaders
       │            ├── BootstrapData
       │            ├── DbLookupExtensions
       │            │         └── MiniblocksMetadata
       │            ├── DbLookupExtensions_ResultsHashesByTx
       │            ├── Logs
       │            ├── MetaBlock
       │            ├── MiniBlocks
       │            ├── PeerAccountsTrie
       │            ├── PeerAccountsTrieCheckpoints
       │            ├── PeerBlocks
       │            ├── Receipts
       │            ├── RewardTransactions
       │            ├── ScheduledSCRs
       │            ├── Transactions
       │            └── UnsignedTransactions
       .......
       ├── Epoch_999
       │      └── Shard_metachain
       │            ├── AccountsTrie 
       │            ├── AccountsTrieCheckpoints
       │            ├── BlockHeaders
       │            ├── BootstrapData
       │            ├── DbLookupExtensions
       │            │         └── MiniblocksMetadata
       │            ├── DbLookupExtensions_ResultsHashesByTx
       │            ├── Logs
       │            ├── MetaBlock
       │            ├── MiniBlocks
       │            ├── PeerAccountsTrie
       │            ├── PeerAccountsTrieCheckpoints
       │            ├── PeerBlocks
       │            ├── Receipts
       │            ├── RewardTransactions
       │            ├── ScheduledSCRs
       │            ├── Transactions
       │            └── UnsignedTransactions
       .......
       └── Static
              └── Shard_metachain
                    ├── DbLookupExtensions_EpochByHash
                    ├── DbLookupExtensions_ESDTSupplies
                    ├── DbLookupExtensions_MiniblockHashByTxHash
                    ├── DbLookupExtensions_RoundHash
                    ├── MetaHdrHashNonce
                    ├── ShardHdrHashNonce0
                    ├── ShardHdrHashNonce1
                    ├── ShardHdrHashNonce2
                    └── StatusMetricsStorageDB
```

We need a tool that will walk through this structure and will try to check if the containing level-DB directories contain a valid level-DB DB.
We should be able to provide the target check directory (node-0, node-1, node-2, node-m) and a range of epochs to check. Also, we should provide a flag to check the static directory or not.
This is the repo containing the level-DB wrapper, if needed: https://github.com/multiversx/mx-chain-storage-go
