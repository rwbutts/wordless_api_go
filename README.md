# wordless_api_go
v1.1.0
Backend API and static file server for wordless.vue written in GO.  
Serves static content out of ./wwwroot, with default file of "index.html".

./wwwroot is where the deployment of project "wordless_vue" is copied.

v1.2.0
Augment http.FileServer functionality to return 404 error when a directory listing would normally be returned by FileServer. Just for extra security.
