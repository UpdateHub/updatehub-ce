![UpdateHub logo](https://github.com/UpdateHub/updatehub/blob/v1/doc/updatehub.png?raw=true)
![CI](https://github.com/UpdateHub/updatehub-ce/workflows/CI/badge.svg)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FUpdateHub%2Fupdatehub-ce.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FUpdateHub%2Fupdatehub-ce?ref=badge_shield)

---

UpdateHub is an enterprise-grade solution which makes simple to remotely update all your Linux-based devices in the field. It handles all aspects related to sending Firmware Over-the-Air (FOTA) updates with maximum security and efficiency, while you focus in adding value to your product.

To learn more about UpdateHub, check out our [documentation](https://docs.updatehub.io).

## Features

* **Yocto Linux support**: Integrate UpdateHub onto your existing Yocto based project
* **Scalable**: Send updates to one device, or one million
* **Reliability and robustness**: Automated rollback in case of update fail
* **Powerful API & SDK**: Extend UpdateHub to fit your needs

# UpdateHub Community Edition

This is a community edition of **UpdateHub Cloud**, so the core concepts and functionality is identical.

See the comparison table below to help you to choose which version fits you need:

| Feature                                      | UpdateHubCE | UpdateHubCloud  |
|:---                                          |    :---:    |      :---:      |
| Secure communication (HTTPS, CoAP over DTLS) | ✘           | ✔              |
| Signed packages                              | ✔           | ✔              |
| Rollouts                                     | ✔           | ✔              |
| Large scale rollouts                         | ✘           | ✔              |
| Multiple organizations                       | ✘           | ✔              |
| Fully monitored updates                      | ✔           | ✔              |
| Teams                                        | ✘           | ✔              |
| HTTP API                                     | ✘           | ✔              |
| Package upload                               | ✔           | ✔              |
| Multiple products                            | ✘           | ✔              |
| Advanced device filter                       | ✘           | ✔              |
| Multiple users                               | ✘           | ✔              |

## Usage

The easyest way to run a Updatehub CE server is to run a Docker image to start a ready-to-use server.

```
$ docker run updatehub/updatehub-ce --help
Usage:
  updatehub-ce [flags]

Flags:
      --db string         Database file (default "updatehub.db")
  -h, --help              help for updatehub-ce
      --password string   Admin password (default "admin")
      --port int          Port (default 8080)
      --username string   Admin username (default "admin")

Example:
docker run -d -p 8080:8080 updatehub/updatehub-ce:latest

```

On the example above, a Docker image will be automatically downloaded and run on 8080 port.
Now you can access the UpdateHub CE dashboard through the host IP on your web browser.

```
http://<Host_IP_Address>:8080/
```

The default login and password is `admin`.

## Building

If you want to build `updatehub-ce` by yourself, follow these steps:

The `updatehub-ce` uses `go mod` to manage its dependencies and
`yarn` to build the web UI for the server. `go mod` should be
included in your default instalation of the `go` toolchain.

You can refer to yarn's documentation page for informations on
how to install it in yor system: https://classic.yarnpkg.com/en/docs/install/

After that, need to install `packr` that is a simple solution for
bundling static assets inside of Go binaries use by
`updatehub-ce`.

To install Packr utility and the dependencies:

```
$ go install github.com/gobuffalo/packr/v2/packr2@latest
```


Finally, you can build `updatehub-ce` as:

```
$ cd <YOUR-UPDATEHUB-CE-PATH>/ui/ && yarn install && yarn run build && cd ..
$ ~/go/bin/packr2 install
$ go build
$ go install
```

Now you can run the `updatehub-ce` as:

```
~/go/bin/updatehub-ce --http 8080
```

## Yocto Project

To integrate updatehub to an Yocto image, you need inherit the `updatehub-image`class
from `meta-updatehub`. Also some variables must be set to work properly.

Add these variables to `conf/local.conf`:

```
UPDATEHUB_SERVER_URL = “http://<Host_IP_Address>:8080"
UPDATEHUB_POLLING_INTERVAL = “1m”
```

* `UPDATEHUB_SERVER_URL` will point to the `updatehub-ce` IP address. If you don't set this variable, your updatehub agent will not communicate with your `updatehub-ce` server, but with [Updatehub Cloud](https://updatehub.io/) instead.
* `UPDATEHUB_POLLING_INTERVAL` will set the interval between each polling between the agent and the server. If you don't set this variable, your updatehub agent will poll every 24 hours.

If you are not familiar with it, you can read the [Yocto Project Reference](https://docs.updatehub.io/yocto-project/yocto-project-reference/) or a [quick start example](https://docs.updatehub.io/quick-starting-with-raspberrypi3/).

## Contributing

UpdateHub is an open source project and we love to receive contributions from our community.
If you would like to contribute, please read our [contributing guide](https://github.com/UpdateHub/updatehub/blob/v1/CONTRIBUTING.md).

## License

UpdateHub Community Edition is licensed under the MIT license. See [LICENSE](LICENSE) for the full license text.


[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FUpdateHub%2Fupdatehub-ce.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FUpdateHub%2Fupdatehub-ce?ref=badge_large)

## Getting in touch

* Reach us on [Gitter](https://gitter.im/UpdateHub/community)
* All source code are in [Github](https://github.com/UpdateHub)
* Email us at [contact@updatehub.io](mailto:contact@updatehub.io)
