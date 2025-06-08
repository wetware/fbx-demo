import asyncio
import capnp
import logging
import os

from tiktok.api.tiktok_client import TikApiClient, TikTokMockClient
from tiktok.api.capability_servers import TikTok

logging.basicConfig(level=logging.DEBUG)  # Add this line
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)


host: str = os.getenv("TIKTOK_HOST", "0.0.0.0")
port: int = int(os.getenv("TIKTOK_PORT", "6060"))
mock: bool = bool(os.getenv("TIKTOK_MOCK", "True"))


tt_client = TikApiClient() if not mock else TikTokMockClient()


async def new_connection(stream):
    await capnp.TwoPartyServer(stream, bootstrap=TikTok(tt_client)).on_disconnect()


async def main():
    logger.info("Start capability server.")
    server = await capnp.AsyncIoStream.create_server(new_connection, host, port)
    async with server:
        await server.serve_forever()


if __name__ == "__main__":
    asyncio.run(capnp.run(main()))
