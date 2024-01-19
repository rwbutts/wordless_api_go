# wordless_api_go
v1.1.0
Backend API and static file server for wordless.vue written in GO.  
Serves static content out of ./wwwroot, with default file of "index.html".

./wwwroot is where the deployment of project "wordless_vue" is copied.

v1.2.0
Augment http.FileServer functionality to return 404 error when a directory listing would normally be returned by FileServer. Just for extra security.

v1.3.0 
Add HTTP listen configuration via config.json file or environment variable. 

First, reads "./config.json", expecting, e.g.:
    { 
        "listenAddress" : ":8080" 
    }

Second (overriding json, if present) attempts to set address from environment variable LISTEN_ADDRESS, e.g. exported like this:
    $ export LISTEN_ADDRESS=:80

Fallback default listen address is ":5090"

v1.3.1
Bugfix: on Windows platform, index.html url assembly accidentally used backslash 
separator, resulting in HTTP 500 failure loading the index page.
