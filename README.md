# HomeKit Bluetooth Low Energy Occupancy Detection

[![Go Workflow][go_workflow_badge]][go_workflow]
[![Coverage Status][coverage_badge]][coverage]
[![Go Report][report_badge]][report]
[![Latest Release][release_badge]][release]
[![License][license_badge]][license]

---

## Table of Contents

1. [Introduction](#introduction)
1. [Usage](#usage)
1. [Contributing](#contributing)
1. [License](#license)

## Introduction

_HomeKit Bluetooth Low Energy Occupancy Detection_ utilizes Bluetooth Low Energy
(BLE) to detect nearby devices and derive the occupancy status of a room from
their presence and signal strength.

## Installation

### Download the pre-compiled and archived binary manually

Binary releases are available on [GitHub Releases][1].

  [1]: https://github.com/lukasmalkmus/homekit-ble-occupancy/releases/latest

### Install from source

```shell
$ git clone https://github.com/lukasmalkmus/homekit-ble-occupancy.git
$ cd homekit-ble-occupancy
$ make build
```

## Usage

1. Run without an argument to scan for nearby devices:

```shell
$ homekit-ble-occupancy
```

2. Provide the device(s) to track plus their signal strength at which they are
   considered close:

```shell
$ homekit-ble-occupancy 3577d0ee-61da-445b-bcf1-704265437842+75
```

## Contributing

Feel free to submit PRs or to fill issues. Every kind of help is appreciated. 

Before committing, `make` should run without any issues.

## License

&copy; Lukas Malkmus, 2021

Distributed under MIT License (`The MIT License`).

See [LICENSE](LICENSE) for more information.

<!-- Badges -->

[go_workflow]: https://github.com/lukasmalkmus/homekit-ble-occupancy/actions?query=workflow%3Ago
[go_workflow_badge]: https://img.shields.io/github/workflow/status/lukasmalkmus/homekit-ble-occupancy/go?style=flat-square&dummy=unused
[coverage]: https://codecov.io/gh/lukasmalkmus/homekit-ble-occupancy
[coverage_badge]: https://img.shields.io/codecov/c/github/lukasmalkmus/homekit-ble-occupancy.svg?style=flat-square&dummy=unused
[report]: https://goreportcard.com/report/github.com/lukasmalkmus/homekit-ble-occupancy
[report_badge]: https://goreportcard.com/badge/github.com/lukasmalkmus/homekit-ble-occupancy?style=flat-square&dummy=unused
[release]: https://github.com/lukasmalkmus/homekit-ble-occupancy/releases/latest
[release_badge]: https://img.shields.io/github/release/lukasmalkmus/homekit-ble-occupancy.svg?style=flat-square&dummy=unused
[license]: https://opensource.org/licenses/MIT
[license_badge]: https://img.shields.io/github/license/lukasmalkmus/homekit-ble-occupancy.svg?color=blue&style=flat-square&dummy=unused
