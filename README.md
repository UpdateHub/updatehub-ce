UpdateHub Community Edition
===========================

> This is community edition of [UpdateHub](https://updatehub.io).
An end-to-end solution for large scale over-the-air update of devices.

## Demo

<img align="center" src="docs/device_list.png"/>

You can try out the public demo instance: http://demo.updatehub.io/
(login with **admin**/**admin**).

Please note that the public demo instance is **reset every 15min**.

## Usage

```
$ ./updatehub-ce-server --help
Usage:
  updatehub-ce-server [flags]

Flags:
      --db string         Database file (default "updatehub.db")
  -h, --help              help for updatehub-ose-server
      --password string   Admin password (default "admin")
      --port int          Port (default 8080)
      --username string   Admin username (default "admin")
```