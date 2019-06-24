# 4gt (Forget)

4gt is a Search Engine for personal notes. The goal is to have a
search engine that is fast and ranks results based on
recency.

4gt is backed by the Bleve Search Engine. At this time, Bleve is the
most mature search library for Go. For personal use, Bleve can
find and rank documents within milliseconds.

I use 4gt as a personal note keeper, as a replacement for
Evernote. Most of my documents are simple text files - a combination
of Git and Bleve gives me full ownership of my data.

4gt has custom parsing for org-mode files - Each heading in a document
is scored separately and only relevant headings are returned. It runs
in the CLI, similar to grep. Future iterations will use a curses-like
console to allow for search-while-you-type UI.

Architecturally, 4gt is separated into server and client components,
communicating over Go RPC. This allows flexibility in developing new
client UIs for this project.

I usually run the 4gt in a terminal window:
```
$ 4gt svr
Using config file: /home/jcheng/.forget.toml
Starting rpc on port 8181
```

When I need to remember something, I simply run 4q (bash alias of `4gt qc`):
```
$ 4q stuff
Using config file: /home/jcheng/.forget.toml
Found 6 notes in 1.190284ms
/home/jcheng/org/writing/writing.org:Give the reader some candy: In a story, candy is the stuff that people will talk about, and the meal is the things they'll dwell on and process more
```

# tldr;

Go version 1.12+ is required. 4gt uses Go modules for dependencies management and 1.12-specific APIs.

Install 4gt
```
$ make all     # Runs tests and creates out/grs
$ make install # Installs grs in $HOME/bin
```

Create a configuration file in `~/.forget.toml`
```
test -f ~/.forget.toml || cat > ~/.forget.toml <<ENDL
# where to place the search engine index
indexDir = "/tmp/forget-index"

# the directories on your machine that you wish to index
dataDirs = [
  "/home/jcheng/org",
  "/home/jcheng/something",
  "/home/jcheng/project_x",
]
ENDL
```

Run the 4gt server, rebuilding the index from scratch
```
$ 4gt svr --rebuild
```

Run the 4gt client (in a different console):
```
4q stuff
Using config file: /home/jcheng/.forget.toml
Found 6 notes in 1.3541ms
/home/jcheng/org/writing/writing.org:Give the reader some candy: In a story, candy is the stuff that people will talk about, and the meal is the things they'll dwell on and process more
```
