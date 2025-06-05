@0xf50adb7bac8b628f;

struct User {
    id       @0 :Text;
    uid      @1 :Text;
    nickname @2 :Text;
}

struct Comment {
    id      @0 :Text;
    author  @1 :User;
    text    @2 :Text;
    replies @3 :List(Comment);
}

interface TikTok {
    mention @0 () -> (comment :Comment);
}
