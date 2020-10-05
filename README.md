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

## Building

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
$ go get -u github.com/gobuffalo/packr/packr
$ go get -u github.com/gobuffalo/packr
```


Finally, you can build `updatehub-ce` as:

```
$ cd <YOUR-UPDATEHUB-CE-PATH>/ui/ && yarn install && yarn run build && cd ..
$ packr install
$ go build
$ go install
```

Now you can run the `updatehub-ce` as:

```
./updatehub-ce --http 8080
```

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
