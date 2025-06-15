using Go = import "/go.capnp";

$Go.package("tiktok");
$Go.import("github.com/wetware/fbx-demo/app");
@0xfe4bf44cd2322767;

struct User {
    # Unique username, e.g. `mikel112358`
    id       @0 :Text;
    # Unique numerical ID assigned by TikTok, e.g. `1234567890`
    uid      @1 :Text;
    # Nickname chosen by the user, e.g. `mikel`
    nickname @2 :Text;
}

struct Comment {
    id       @0 :Text;
    mediaId @1 :Text;
    author   @2 :User;
    text     @3 :Text;
    replies  @4 :List(Comment);
}

interface TikTok {
    # Get a comment mentioning the current user.
    mention @0 () -> (comment :Comment);
    # Get every comment on a video.
    comments @1 (mediaId :Text) -> (comments :List(Comment));
    # Reply to a mention with the specified ID.
    reply @2 (mediaId :Text, commendId: Text, response :Text) -> ();
}

# Alternative design: create an Object capability and reply by interacting directly with the object.
# Pros: more intuitive, low-grain capability control.
# Cons: more complex to implement, requires one capability server (in-memory object) per comment.
#interface Comment {
#    # Getters.
#    id @0 () -> (id :Text);
#    media_id @1 () -> (media_id :Text);
#    author @2 () -> (author :User);
#    text @3 () -> (text :Text);
#    replies @4 () -> (replies :List(Comment));
#    # Methods to interact with the comment!
#    reply @5 (response :Text) -> ();
#}
