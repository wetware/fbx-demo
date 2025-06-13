# Flashbots_X Demo

> [!IMPORTANT]
> Make sure to clone the repository with the `--recursive` option.

## Project Structure

This demo app consists of a TikTok reply bot split into four services:

* `dstack`: Running on the host machine. The rest of services run on Docker.
* `wetware`: our [Wetware](https://github.com/wetware/pkg) platform running on the TDX emulator. The service will be running Wetware processes:
WASM guest processes that can only reach outside their Sandbox through Object Capabilities.
* `app`:
  - `guest`: The TikTok bot running on Wetware with a TikTok [object capability](https://capnproto.org/). The app is written in Go and compiled to WASM.
  - `host`: A simple Go program that connects to a Wetware node, bootstraps the TikTok object capability, and starts the bot process (`guest`) on that node with that capability.
* `tiktok`: a service receiving TikTok notifications, sending them with some context to a bot, and posting the bot's replies. The service is using TikApi,
because [the suggested alternative](https://github.com/davidteather/TikTok-Api) does not provide the notification or tagging functionality
described in the document.

## Setup

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

### Environment Variables

Besides the `DSTACK_SOCKET` and `TAPPD_SOCKET` environment variables, you can use the ones in the `.env` to configure the project.

Make sure to `source` the `.env` file included in this repository, and modify it to have values for the `API_KEY` and `ACCOUNT_KEY` variables
if not mocking TikTok responses. Here's the sample `.env` file:

```shell
API_KEY=<tikapi api key>
ACCOUNT_KEY=<tikapi account key>
TIKTOK_MOCK=True
TIKTOK_HOST=0.0.0.0
TIKTOK_PORT=6060
```

## Build

After setting up Dstack, just run:

```
make build
```

This will:

* Compile Dstack (in case you hadn't already)
* Build an image with the TikTok object capability server
* Build the Wetware docker image
* Build the bot app, both `host` and `guest` programs.

> Note: the `make build` command builds the image with `docker` instead of `phala docker` because
> it let's us set the image author to `wetware` without logging into docker with the account.

## Run

Then run:

```
make run
```

This will:

* Start the Dstack simulator
* [container] Start the TikTok object capability server
* [container] Start a Wetware node
* Start the bot app:
  1. [container] Run the `host` program to bootstrap the TikTok capability, connect to the Wetware node, load the WASM bytecode and start the bot program in the Wetware node.
  2. [wetware] Run the `guest` program in a loop: fetch a mention, grab some context, generate a response, and reply to the mention.
