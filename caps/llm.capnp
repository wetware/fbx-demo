@0xf641f3e21bc55b4e;

struct Input {
    mention @0 :Text;
    context @1 :List(Text);
}

interface LLM {
    generateResponse @0 (input :Input) -> (output: Text);
}
