# MX Deep History Checker - Specifications

## 1. Overview
The **mx-deep-history-checker** is a Go-based command-line tool designed to verify the integrity of level-DB database structures within MultiversX nodes. It walks through the complex, epoch-based directory structure of node data partitions (shards and metachain) and systematically checks the readability and corruption status of the underlying level-DB databases.

## 2. Requirements & Inputs
The tool will be executed from the command line and will accept the following input parameters:

- `--node-dir` (string, required): The root path to the node data directory (e.g., `/path/to/node-0`, `/path/to/node-m`). The tool will look inside the `1/` subdirectory by default to start the traversal.
- `--start-epoch` (uint32, optional): The starting epoch number to check (inclusive). Default: `0`.
- `--end-epoch` (uint32, optional): The ending epoch number to check (inclusive). If omitted, the tool will check all available epochs starting from `--start-epoch` up to the highest one found in the directory.
- `--check-static` (bool, optional): A flag indicating whether to check the `Static` directory databases (e.g., `DbLookupExtensions_EpochByHash`, `StatusMetricsStorageDB`, etc.). Default: `true`.
- `--parallel-epochs` (uint, optional): The number of epochs to process in parallel using goroutines. Default: `4`.

## 3. Core Behavior & Workflow

### 3.1 Directory Traversal
1. The tool will begin traversal at `<node-dir>/1/`.
2. It will identify and iterate over `Epoch_X` directories that fall within the range defined by `--start-epoch` and `--end-epoch`.
3. Within each valid `Epoch_X/Shard_Y` directory, it will detect all subdirectories containing databases (e.g., `AccountsTrie`, `MiniBlocks`, `Transactions`, etc.).
   - **Sharded Database Handling**: Complex directories like `AccountsTrie` may be further sharded into multiple subdirectories (e.g., `0`, `1`, `2`, `3`...) instead of directly containing level-DB files. The tool must detect this by checking for the presence of a `config.toml` file (defining properties like `NumShards` and `Type = "LvlDBSerial"`). When encountered, the tool should verify each of the indicated DB shards according to the definitions provided in the `mx-chain-storage-go` repository.
4. If `--check-static` is set to true, it will additionally traverse the `Static/Shard_Y` directory and identify its database components.

### 3.2 Database Verification
1. For each identified database directory, the tool will attempt to instantiate a level-DB connection.
2. It will utilize the MultiversX level-DB wrapper provided by `github.com/multiversx/mx-chain-storage-go`.
3. The integrity check process includes:
   - Initializing and opening the level-DB.
   - Using the underlying leveldb library's built-in repair/verification capabilities or executing basic read operations to ensure the database is not corrupted and can be successfully mounted.
4. Immediately closing the database to avoid resource leaks and locking issues.

### 3.3 Error Handling and Reporting
- **Console Output**: A live progress stream detailing which epoch and database is currently being verified.
- **Validation Results**: Explicitly logging `[OK]` for valid databases and `[ERROR]` + the underlying level-DB error for corrupted/unopenable directories. The application should stop at the first encountered error.

## 4. Software Architecture & Design
To adhere to the project's golang rules and ensure high decoupling, the architecture will be interface-driven:

### 4.1 Key Interfaces
Instead of exposing functional structs, the tool will expose interfaces to separate responsibilities:

- **`ConfigHandler`**: Parses and validates the command-line flags.
- **`DirectoryTraverser`**: Responsible for listing the relevant directories based on the configuration and epochs given.
- **`DatabaseChecker`**: Responsible for invoking the DB library to attempt opening/verifying a given level-DB path.
- **`Reporter`**: Aggregates the results and handles the stdout/logging formats.

### 4.2 Error Management
Custom error types will be defined for specific failure scenarios, such as `ErrInvalidEpochRange`, `ErrDirectoryNotFound`, and `ErrDatabaseCorrupted`.

## 5. Extensibility
The interface-driven design allows dropping in future storage backends (e.g., if MultiversX nodes transition to a different KV store like Pebble, we would simply implement a new `DatabaseChecker` interface).
