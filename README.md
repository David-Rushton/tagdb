# TagDb

> ⚠️⚠️⚠️ This repository is under construction ⚠️⚠️⚠️  
> Most of what is here is unfinished and does not work

A persist key-value store with searchable tags.

## Things to make and do

## Backend

- ✅ In-mem only
- ✅ Transactions
- ✅ WAL backed
- ✅ WAL rolling
- ✅ Restore in-mem at start up

## Web Server

### Web API

- ✅ Implement all endpoints
- ✅ Containerize 
  - ✅ Run db on start up via Docker restart policy

### Web App

- ✅ Don't CSS forever
- ✅ Support all endpoints

## Link Ext

- ✅ Rewrite using tagdb
- ✅ Support custom titles

## Command Box

- ✅ Rewrite using tagdb

## CLI

- ✅ Command parsing
- ✅ Routes
- ⌛ Arguments - variadic for last only
- ⌛ Options - both -t a -t b and -t a b 
- ⌛ Add -- to switch to posistional mode only
- 🆕 Move to interface
- 🆕 context.Context support
- 🆕 Multi flag support -a -b -c == -abc
- 🆕 Array arguments
- 🆕 Array options
- 🆕 Better Arg/Option errors
- 🆕 Arg/option validation
- 🆕 POSIX IEEE Std 1003.2-1992
- 🆕 Don't error - panic
- 🆕 Components
- 🆕 Support all endpoints
