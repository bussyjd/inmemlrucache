# Inmemlrucache

## Description
This is a golang implementation of a LRU cache for storage of ~1mb image files. Data is stored in a tmpfs filesystem to leverage concurrency safety.
Tested on Ubuntu Linux.

## Context restrictions

* Thread safety
* Stored in Memory
* Capacity of 10 entries for ~1Mb images
* Simple API
* Well tested

## API spec

### /SET
_POST_

Stores provided image on the top of the LRU cache.
The image data is provided in the body of the http request

HTTP 200 on success

### /GET/{id}
_GET_

Returns the image stored in the {id} entry of the LRU cache and promotes it on the front of the cache list.

HTTP 200 on success

###  /DEL/{id}
Deletes the key {id} out of the cache and promote the olders entries

returns HTTP 204 on success
returns HTTP 404 if the entry was empty beforehand

### /RESET
Completely clears the cache in a irreversible manner.

returns HTTP 204 on success

### /COUNT
Returns the number of entries stored in the cache.

returns HTTP 200 on success

## How to run

### Docker
```
docker run --rm -v $PWD:/app -w /app treeder/go remote https://github.com/bussyjd/inmemlrucache.git
```
### Local env
```
go build
./inmemlrucache
```
### Development testing
```
cd inmemlrucache
docker run --rm -v $PWD:/app -w /app treeder/go vendor
docker run --rm -v $PWD:/app -w /app treeder/go build
docker run --rm -v $PWD:/app -w /app -p 8080:8080 iron/base ./app
```

## Testing
```
go test
``` 

## Q&A
_Why use a tmfps partition?_

Tmpfs take care of the R/W concurrency and a good candidate makes senses for a small i-memory image storage solution

_Why golang?_

* Good performance
* Concurrency at the core of the language
* Easy to pbuild and API on top of golang