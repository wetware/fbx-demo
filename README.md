# Flashbots_X Demo

> [!IMPORTANT]
> Make sure to clone the repository with the `--recursive` option.

## Project Structure

This demo app consists of a TikTok reply bot split into four services:

* `dstack`: Running on the host machine. The rest of services run on Docker.
* `app`: The TikTok bot app running on Wetware with a TikTok [object capability](https://capnproto.org/). The app is written in Go and compiled to WASM.
* `tiktok`: a service receiving TikTok notifications, sending them with some context to a bot, and posting the bot's replies. The service is using TikApi,
because [the suggested alternative](https://github.com/davidteather/TikTok-Api) does not provide the notification or tagging functionality
described in the document.
* `wetware`: our [Wetware](https://github.com/wetware/pkg) platform running on the TDX emulator. The service will be running Wetware processes:
WASM guest processes that can only reach outside their Sandbox through Object Capabilities.

## Building

Create a `.env` file with these fields. `API_KEY` and `ACCOUNT_KEY` can be omited if `TIKTOK_MOCK` is `True`.

```shell
API_KEY=<tikapi api key>
ACCOUNT_KEY=<tikapi account key>
TIKTOK_MOCK=True
TIKTOK_HOST=0.0.0.0
TIKTOK_PORT=6060
```

Then run

```
make build
```

### Running:

The following command will create a Dstack simulator and then spawn the rest of the services with `docker-compose up`.

```
make run
```

---

---

## Requiements

* Cargo
* Phala
* Dstack
* A TikTok API key and an account key

## Setup

### Phala

Just run:

```shell
npm install -g phala
```

### Dstack

We need to set up the simulator to get working. [The guide](https://docs.phala.network/phala-cloud/references/tee-cloud-cli/phala/simulator) works
with [the Python example](https://github.com/Phala-Network/python-tee-sim-template), but it returns `404` errors for the Go SDK calls.
Perhaps I'm doing something wrong.

Following [the simulation instructions from the Go SDK](https://github.com/Phala-Network/dstack/tree/0dd93563763a3fb2da3b5fb52c953894abcd3ebf/sdk/go#development) does work.
Make sure you *don't* have this in your `~/.gitconfig` when running the `./build.sh` step, as it is a common setting to have when working with private Go repos.

```toml
[url "git@github.com:"]
        insteadOf = https://github.com/
```

The process is:

```shell
git clone git@github.com:Phala-Network/dstack.git
cd dstack/sdk/simulator
# Save this path to and environment variable
echo "export DSTACK_SIMULATOR_DIR=\"$PWD\"" >> ~/.zshrc
echo "export DSTACK_SOCKET=\"\$DSTACK_SIMULATOR_DIR/dstack.sock\""  >> ~/.zshrc
echo "export TAPPD_SOCKET=\"\$DSTACK_SIMULATOR_DIR/tappd.sock\""  >> ~/.zshrc
./build.sh  # requires Cargo
./dstack-simulator
```

### Wetware Container

Make sure the `DSTACK_SOCKET` and `TAPPD_SOCKET` environment variables are set, then run:

```shell
make build
make run
```

> Note: the `make build` command builds the image with `docker` instead of `phala docker` because
> it let's us set the image author to `wetware` without logging into docker with the account.
