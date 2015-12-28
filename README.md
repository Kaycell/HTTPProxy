# HTTPProxy
Simple golang HTTP proxy

###Command Line Arguments

```
Usage of ./HTTPProxy:
  -api-port int
    	Proxy admin api listening port (default 8001)
  -b string
    	Path to blacklist file
  -p int
    	Proxy listening port (default 8000)
  -rl string
    	Request logfile
  -w string
    	Path to whitelist file

Example:  ./HTTPProxy -p 8000 -api-port 8001 -rl log.txt
```
###Design Decisions

This simple http transparent proxy use the Host field in every request to forward requests.

Requests and their execution time are logged into the output specified or by
default on standard output.

By default localhost is blacklisted and if incoming requests have
empty host or host equal to localhost then the proxy return a bad gateway code
The proxy allow to blacklist or whitelist endpoints with regular expressions.
If an expression is blacklisted and whitelisted the priority is on blacklist.

These list could be provided with files, then files are parsed and not use again
during running time to improve performance and avoid i/o access.
The API is here to allow to blacklist or whitelist endpoint during running time
by two routes : '/allow' and '/forbid'.
