# XBank [Golang + Postgres + gRPC + Kubernetes]


## Setup


### Tools 
 - Docker
 - Postgres
 - TablePlus/Azure Data Studio
 - [Golang Migrate](https://github.com/golang-migrate/migrate)
 - [sqlc](https://sqlc.dev/) 
```shell 
brew install sqlc
```
 - [Golang Postgres Driver for database/sql](https://github.com/lib/pq)


### Section 1 Highlights
 - DB Schema can be created using dbdiagram.io
    ![dbdiagram.io](assets/db_schema.png)
 - Docker, Postgres(Docker Version), TablePlus or Azure Data Studio Installation
 - `sqlc` can be used to generate the Golang Code (check db/query folder for sql queries and db/sqlc for generated Golang code)
 - Added Unit Tests for all CRUD operations (Accounts, Entries, Transfers TABLES)
 - What is a DB Transaction
   - A transaction in SQL is a sequence of one or more operations such as insertions, updates, or deletions, performed on a database as a single unit of work.
   <img src="assets/db_transaction.png" alt="image" width="380" height="auto">
   <img src="assets/transaction_example.png" alt="image" width="380" height="auto">
 - Deadlock can occur while executing the transactions
  <img src="assets/transaction_deadlock_queries.png" alt="image" width="380" height="auto"> 
  <img src="assets/transaction_deadlock_shell.png" alt="image" width="380" height="auto">
 -  


### Section 2 Highlights
