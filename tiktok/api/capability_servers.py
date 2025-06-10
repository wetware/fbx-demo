import capnp
import logging
import os

from asyncio import Queue
from typing import List

from tikapi import TikAPI, ValidationException, ResponseException
from tikapi.api import APIResponse


from tiktok.api import tiktok_capnp
from tiktok.api.capability_models import Comment
from tiktok.api.tiktok_client import TikTokClient


logging.basicConfig(level=logging.DEBUG)  # Add this line
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)


class TikTok(tiktok_capnp.TikTok.Server):
    def __init__(self, tiktok_client: TikTokClient):
        super().__init__()
        self.tiktok_client = tiktok_client

    async def mention(self, **kwargs) -> Comment:
        logger.info("Getting mention.")
        return await self.tiktok_client.get_mention()

    async def comments(self, mediaId: str, _context=None) -> List[Comment]:
        logger.info(f"Getting comments for media {mediaId}.")
        return await self.tiktok_client.get_comments(mediaId)

    async def reply(self, media_id: str, comment_id: str, response: str, **kwargs):
        logger.info(f"Replying to comment {media_id}:{comment_id}: {response}.")
        return await self.tiktok_client.reply(media_id, comment_id, response)
