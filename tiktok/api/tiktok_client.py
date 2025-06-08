from abc import ABC, abstractmethod
import logging
import os

from asyncio import Queue
from typing import Union

from tikapi import TikAPI, ValidationException, ResponseException
from tikapi.api import APIResponse

from tiktok.api.capability_models import Comment, User


# fmt: off
API_FETCH_LIMIT = 20
SAMPLE_USER_1 = User("1234567890123456789", "uniqueusername", "louis")
SAMPLE_USER_2 = User("1234567890123456789", "otheruniqueusername", "mikel")
SAMPLE_REPLY_3 = Comment("1234567890123456789", "9876543210123456789", SAMPLE_USER_1, "Reply 3", [])
SAMPLE_REPLY_2 = Comment("1234567890123456789", "9876543210123456789", SAMPLE_USER_2, "Reply 2", [SAMPLE_REPLY_3])
SAMPLE_REPLY_1 = Comment("1234567890123456789", "9876543210123456789", SAMPLE_USER_2, "Reply 1", [])
SAMPLE_COMMENT = Comment("1234567890123456789", "9876543210123456789", SAMPLE_USER_1, "Hello World!", [SAMPLE_REPLY_1, SAMPLE_REPLY_2])
# fmt: on


logging.basicConfig(level=logging.DEBUG)  # Add this line
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)


class TikTokClient(ABC):
    @abstractmethod
    async def get_mention(self) -> Comment:
        pass

    @abstractmethod
    async def reply(self, media_id: str, comment_id: str, text: str):
        pass


class TikTokMockClient(TikTokClient):
    async def get_mention(self) -> Comment:
        logger.info("[mock] Getting mention")
        return SAMPLE_COMMENT.cap()

    async def reply(self, media_id: str, comment_id: str, text: str):
        logger.info(f"[mock] Replying to comment {media_id}:{comment_id}: {text}")
        pass


class TikApiClient(TikTokClient):

    def __init__(self):
        api_key: str = os.getenv("API_KEY", "")
        account_key: str = os.getenv("ACCOUNT_KEY", "")
        if api_key == "" or account_key == "":
            raise ValueError("api_key and/or account_key is not set")

        api = TikAPI(api_key)
        self.user = api.user(accountKey=account_key)

        self.mention_queue = Queue()

    async def get_mention(self) -> Comment:
        if self.mention_queue.qsize() == 0:
            logger.info("No mention available, fetching from TikTok.")
            await self.fetch_mentions()
            # TODO: case in which no new mentions are fetched
        return await self.mention_queue.get()

    async def fetch_mentions(self):
        logger.info("Fetching mentions from TikTok.")
        try:
            response = self.user.notifications(filter="mentions", count=API_FETCH_LIMIT)
            while response:
                parsed_mention, next_response = parse_mention(response)
                response = next_response
                await self.mention_queue.put(parsed_mention)
        except Exception as e:
            logger.error(f"Error fetching mentions: {e}")

    async def reply(self, media_id: str, comment_id: str, text: str):
        try:
            self.user.posts.comments.post(media_id=media_id, reply_comment_id=comment_id, text=text)
        except ValidationException as e:
            logger.error(f"{media_id}:{comment_id}: Validation error: {e} - {e.field}")

        except ResponseException as e:
            logger.error(f"{media_id}:{comment_id}: Response error: {e} - {e.response.status_code}")


def parse_mention(mention: Union[APIResponse, bytes]) -> APIResponse:
    if type(mention) == bytes:
        raise ValueError(f"Mention format: `{mention}`")
    # TODO: extract comment body and send it to Wetware agent.
    try:
        json_response = mention.json()
        notice = json_response["notice_list"][0]["notice_list"][0]["at"]
        author = User(
            id=notice["comment"]["user"]["unique_id"],
            uid=notice["comment"]["user"]["uid"],
            nickname=notice["comment"]["user"]["nickname"],
        )

        comment = Comment(
            author=author,
            id=notice["comment"]["cid"],
            media_id=notice["comment"]["aweme_id"],  # TODO mikel: validate this field is media_id
            text=notice["content"],
            replies=[],  # TODO mikel: parse replies
        )
    except (KeyError, IndexError) as e:
        logger.error(f"Error parsing mention: {e}")
    return mention.next_items()
