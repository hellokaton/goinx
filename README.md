# goinx

💞 domain proxy server written in golang

## Feature

- Support static server
- Support multi domain proxy
- Support HTTPS
- Support GFW reverse proxy

## Usage

**By Binary**

Go [Releases](https://github.com/biezhi/goinx/releases) download the corresponding platform.

**By Golang**

```bash
» ./goinx
💖  Goinx 0.0.1
Author: biezhi
Github: https://github.com/biezhi/goinx

Usage: goinx [start|stop|restart]

Options:

    --config    Configuration path
    --help      Help info
```

[Document](https://github.com/biezhi/goinx/wiki)

## Config File

```bash
log_level: info
access_log:
http:
  servers:
    - name:       demo1
      listen:     ":9001"
      domains:    [localhost, www.biezhi.com]
      proxy_pass: http://127.0.0.1:8080
      cert_file:
      key_file:
    - name:       demo2
      listen:     ":9002"
      domains:    [localhost]
      root: /Users/biezhi/workspace/wwwroot/www.jq22.com/demo/bootstrap-moban20150917
      # ssl: true
      # cert_file: /Users/biezhi/workspace/ssl/cert.pem
      # key_file: /Users/biezhi/workspace/ssl/key.pem
    - name:       demo3
      listen:     ":9003"
      gfw: true
      domains:    [www.biezhi.com]
      proxy_pass: "https://www.google.com"
```
