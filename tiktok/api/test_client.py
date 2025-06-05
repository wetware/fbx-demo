#!/usr/bin/env python3

# Based on: https://github.com/capnproto/pycapnp/blob/master/examples/async_client.py

from pprint import pprint
import asyncio
import logging
import os
import time

import capnp

from . import capability


logging.basicConfig(level=logging.DEBUG)  # Add this line
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)


host: str = "127.0.0.1"
port: int = int(os.getenv("TIKTOK_PORT", "6060"))


async def print_mentions(cap):
    logging.info("[call] Client get mentions.")
    mention = await cap.mention()
    pprint(mention)


async def main(host):
    logging.info("Start client connection.")
    connection = await capnp.AsyncIoStream.create_connection(host=host, port=port)
    logging.info("Bootstrap capability from client.")
    client = capnp.TwoPartyClient(connection)
    cap = client.bootstrap().cast_as(capability.tiktok_capnp.TikTok)

    # Start background task for subscriber
    task = asyncio.ensure_future(print_mentions(cap))
    await task


if __name__ == "__main__":
    asyncio.run(capnp.run(main(host)))

    # Test that we can run multiple asyncio loops in sequence. This is particularly tricky, because
    # main contains a background task that we never cancel. The entire loop gets cleaned up anyways,
    # and we can start a new loop.
    asyncio.run(capnp.run(main(host)))
