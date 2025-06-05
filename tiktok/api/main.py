import asyncio
import capnp
import logging
import os

from tikapi import TikAPI, ValidationException, ResponseException
from tikapi.api import APIResponse

from . import capability
from pip._vendor.distro import info


logging.basicConfig(level=logging.DEBUG)  # Add this line
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)


host: str = os.getenv("TIKTOK_HOST", "0.0.0.0")
port: int = int(os.getenv("TIKTOK_PORT", "6060"))


async def main():
    logger.info("Start server.")
    server = await capnp.AsyncIoStream.create_server(capability.new_connection, host, port)
    async with server:
        await server.serve_forever()


if __name__ == "__main__":
    asyncio.run(capnp.run(main()))

# api_key: str = os.getenv("API_KEY", "")
# account_key: str = os.getenv("ACCOUNT_KEY", "")


# if api_key == "" or account_key == "":
#     raise ValueError("api_key and/or account_key is not set")


# api = TikAPI(api_key)
# User = api.user(accountKey=account_key)


# def notification_loop():
#     from json import dumps
#     while True:
#         try:
#             response = User.notifications(filter="mentions", count=5)
#             print(dumps(response.json()))
#             while response:
#                 response = process_notification(response)
#         except ValidationException as e:
#             print(e, e.field)

#         except ResponseException as e:
#             print(e, e.response.status_code)
#         finally:
#             exit(0)


# def process_notification(response: APIResponse) -> APIResponse:
#     # TODO: extract comment body and send it to Wetware agent.
#     min_time = response.json().get("min_time")  # TODO: change field
#     return response.next_items()


# if __name__ == "__main__":
#     notification_loop()
