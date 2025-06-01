# fbx-demo

## Requiements

* Cargo
* Phala
* Dstack

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
