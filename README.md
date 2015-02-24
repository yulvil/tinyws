tinyws
======

Command line:
```
tinyws --cfg tinyws.cfg
tinyws.cfg --cfg '{"http_port": 80, "https_port": 443, "ssl_certificate": "cert.pem", "ssl_certificate_key": "key.pem", "handlers": [{"type": "file", "context_path": "/static/", "directory": "/home/www/static/"}]}'
```

tinyws.cfg
```
{
    "http_port": 80,
    "https_port": 443,
    "ssl_certificate": "cert.pem",
    "ssl_certificate_key": "key.pem",
    "handlers": [
        {
            "type": "file",
            "context_path": "/static/",
            "directory": "/home/www/static/"
        },
        {
            "type": "file",
            "context_path": "/downloads/",
            "directory": "/home/downloads/"
        },
        {
            "type": "proxy",
            "backend": "http://my.server.com:8080",
            "context_path": "/myapp/"
        },
        {
            "type": "proxy",
            "backend": "http://my2.server.com:8888",
            "context_path": "/otherapp/"
        }
    ]
}
```
