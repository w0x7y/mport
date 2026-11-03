# mport

See what's listening on your ports, and kill it if you want.

![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white)
![Platform](https://img.shields.io/badge/platform-Linux-lightgray?logo=linux&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green)

## Install

```sh
git clone https://github.com/w0x7y/mport.git
cd mport
go install
```

Requires `ss` (from `iproute2`).

## Usage

| Command | Description |
|---|---|
| `mport` | List all listening TCP ports |
| `mport 3000 8080` | Show only these ports |
| `mport -i 3000` | Show the full command of the process on a port |
| `mport -k 3000` | Kill the process on a port (asks first) |

`-i` and `-k` also have long forms: `--info` and `--kill`.

```
$ mport 5432
127.0.0.1:5432 - postgres, 1234
```

If the process shows as `unknown`, run with `sudo` to see processes owned by other users.

## License

[MIT](LICENSE)
