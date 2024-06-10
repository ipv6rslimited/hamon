# Hamon

Hamon is a tool that uses the djb2 (ca. 1991) hash to map IP addresses to readable words and vice versa. It supports both IPv4 and IPv6 addresses.

The name is derived from HAsh MOdulo Name (ハモン not 刃文).

## Why

If you want to remember an IP address or IPv6 address easier by using words, you can use this. It doesn't require internet or a naming system of any kind.

## Features

- Converts words to IP addresses.
- Converts IP addresses to words.
- Displays all word options from an English Words database provided by @dwyl
- Supports both IPv6 and legacy IPv4

## Setup

```
git clone https://github.com/ipv6rslimited/hamon
cd hamon
git submodule init
git submodule update
make os/arch (e.g., linux/arm64 or all)
```

## Usage

Get an IP address from words:
```
./hamon -forward "ipv4.40.ipv6.42"
40.41.42.43
```

Get words from an IP address:
```
./hamon -reverse "40.41.42.43"
protovillain.unclever.griffin.airling
```

## License

Copyright (c) 2024 IPv6rs Limited Company <https://ipv6.rs>

Distributed under the COOL License.
