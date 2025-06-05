import asyncio
import capnp
import logging
import os

from typing import List


logging.basicConfig(level=logging.DEBUG)  # Add this line
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)

_file_path = os.path.abspath(__file__)
_directory = os.path.dirname(_file_path)
capnp_file = os.path.join(_directory, "cap", "tiktok.capnp")


tiktok_capnp = capnp.load(capnp_file)


logger.debug("Parsed TikTok capability file.")


class TikTok(tiktok_capnp.TikTok.Server):
    def __init__(self):
        super().__init__()

    async def mention(self, **kwargs):
        logger.info("[call] Getting mention.")
        return SAMPLE_COMMENT.cap()


async def new_connection(stream):
    logger.debug("New connection established.")
    await capnp.TwoPartyServer(stream, bootstrap=TikTok()).on_disconnect()


class User:
    def __init__(self, id: str, uid: str, nickname: str):
        self.id = id
        self.uid = uid
        self.nickname = nickname

    def cap(self):
        user_cap = tiktok_capnp.User.new_message(
            id=self.id,
            uid=self.uid,
            nickname=self.nickname
        )
        return user_cap



class Comment:
    def __init__(self, id: str, author: User, text: str, replies):
        self.id = id
        self.author = author
        self.text = text
        self.replies = replies

    def cap(self):
        comment_cap = tiktok_capnp.Comment.new_message(
            id=self.id,
            text=self.text,
            author=self.author.cap(),
        )
        replies = comment_cap.init("replies", len(self.replies))
        for i, reply in enumerate(self.replies):
            replies[i] = reply.cap()
        return comment_cap


SAMPLE_USER_1 = User("1234567890123456789", "uniqueusername", "louis")
SAMPLE_USER_2 = User("1234567890123456789", "otheruniqueusername", "mikel")
SAMPLE_REPLY_3 = Comment("1234567890123456789", SAMPLE_USER_1, "Reply 3", [])
SAMPLE_REPLY_2 = Comment("1234567890123456789", SAMPLE_USER_2, "Reply 2", [SAMPLE_REPLY_3])
SAMPLE_REPLY_1 = Comment("1234567890123456789", SAMPLE_USER_2, "Reply 1", [])
SAMPLE_COMMENT = Comment("1234567890123456789", SAMPLE_USER_1, "Hello World!", [SAMPLE_REPLY_1, SAMPLE_REPLY_2])
