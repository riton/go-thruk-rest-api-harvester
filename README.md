# Thruk REST API Harvester

## Description

This simple tool has been designed to create fake traffic at configurable rate to the Thruk REST API
for demonstration purpose.

This tools supports using multiple goroutines to query the Thruk REST API and observe its behavior.

This tool currently only supports listing Nagios/Naemon hosts using Thruk REST API, and all the gouroutines
will perform the exact same query.

## Build and install

```
$ go get https://github.com/riton/go-thruk-rest-api-harvester
$ cd $GOPATH/src/github.com/riton/go-thruk-rest-api-harvester
$ go install .
```

## Usage

```
$ THRUK_ENDPOINT=http://mythruk-install.eu/thruk thruk-rest-api-harvester
[...]
INFO[0001] hosts = hplj2605dn,linksys-srw224p,localhost,winserver  worker-id=4
INFO[0001] hosts = hplj2605dn,linksys-srw224p,localhost,winserver  worker-id=0
INFO[0001] hosts = hplj2605dn,linksys-srw224p,localhost,winserver  worker-id=0
INFO[0001] hosts = hplj2605dn,linksys-srw224p,localhost,winserver  worker-id=0
[...]
```

## Configuration

This tool uses environment vars to find its configuration

**THRUK_ENDPOINT**

Set the Thruk endpoint the requests will be made against.
Defaults to `http://localhost:8080/thruk`.

**THRUK_USER**

Set the Thruk username for Basic authentication.
Defaults to `thrukadmin`.

**THRUK_PASSWORD**

Set the Thruk password for Basic authentication.
Defaults to `thrukadmin`.

**HTTP_TIMEOUT**

Set the HTTP Timeout in seconds used to query Thruk REST API.
Defaults to `3`.

**WORKERS**

Set the number of concurrent goroutines that will harvest the Thruk REST API.
Defaults to `5`.

**WORKER_MAX_WAIT**

Maximum number of seconds a goroutine can wait before performing another request.
Defaults to `5`.

## Author

* Rémi Ferrand ([@riton](https://github.com/riton))
