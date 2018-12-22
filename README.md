# Aqua File Traverser Assignment

Aqua file traverser is a cli that can travers and send information about the existing file.  
The cli can also listen to a port and store the received file info to a mysql db.

## Installation
Run 
```
go get -u github.com/dimanevelev/aqua
go install github.com/dimanevelev/aqua
```

The executable will be at `$GOPATH/bin/aqua`

## Server
This command will start a server an expose two api endpoints:  
* `/api/v1/file` - Receives and stores POST requests with file information payload of the format `{"Path":"/home/example.zip","FileInfo":{"Name":"example.zip","Size":6848,"Mode":436,"ModTime":"2018-12-18T08:12:46.83861937+02:00","IsDir":false}}`.
* `/api/v1/stats` - Receives GET requests and will return statistics of the received files. 
		Result example: `{"code":200,"data":{"TotalFiles":2,"MaxFile":{"Size":3,"Path":"/foo/bar.abc"},"AvgFileSize":1.5,"Extensions":[".abc",".txt"],"TopExtension":".txt"}}`

### Usage:  
* aqua server [flags]

### Flags:  
- -n / --db-name [string]: The MySql DB name. (default "files")
- -p / --db-password [string]: The MySql DB password. (default "password")
* --db / -port [string]: The MySql DB port. (default "3306")
* --db-url [string]: The MySql DB url. (default "127.0.0.1")
* -u / --db-username [string]: The MySql DB username. (default "root")
* -h / --help: help for server
* --port [string]: The server port. For https use 443. (default "8080")

## Travers
Traverses the file system and sends file info to a server. example of file info: 
	{"Path":"/home/example.zip","FileInfo":{"Name":"example.zip","Size":6848,"Mode":436,"ModTime":"2018-12-18T08:12:46.83861937+02:00","IsDir":false}}.

### Usage:
  aqua travers [flags]

### Flags:
* -h / --help:        help for travers
* -p / --path [string]: The path to start traversing. (default ".")
* -t / --threads [int]: Number of threads that will send the requests (default 3)
* -u / --url [string]: The target of the requests. (default "http://localhost:8080/api/v1/file")

## TODO
* Tests