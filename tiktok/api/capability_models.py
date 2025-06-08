from tiktok.api import tiktok_capnp


class User:
    def __init__(self, id: str, uid: str, nickname: str):
        self.id = id
        self.uid = uid
        self.nickname = nickname

    def cap(self):
        user_cap = tiktok_capnp.User.new_message(id=self.id, uid=self.uid, nickname=self.nickname)
        return user_cap


class Comment:
    def __init__(self, id: str, media_id: str, author: User, text: str, replies):
        self.id = id
        self.media_id = media_id
        self.author = author
        self.text = text
        self.replies = replies

    def cap(self):
        comment_cap = tiktok_capnp.Comment.new_message(id=self.id, text=self.text, author=self.author.cap())
        comment_cap.init("replies", len(self.replies))
        for i, reply in enumerate(self.replies):
            comment_cap.replies[i] = reply.cap()
        return comment_cap
