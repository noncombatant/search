# Search

This program is a simplified, easier-to-use implementation of the idea of `find
... -print0 | xargs -0 grep ...`. I write that kind of command line so often,
and `find` and `xargs` are so awkward, and shell glob syntax is so annoyingly
different from regular expressions, that it was easier to write and use this.

This Go implementation replaces an older C++ implementation I used to have up on
GitHub. Go is obviously easier and more fun to extend than C++, has modern
regular expressions, and other such niceness.

You might also enjoy other fine command-line search programs, such as
[ripgrep](https://github.com/BurntSushi/ripgrep), `git grep ...`, and [other
such delights](https://beyondgrep.com/feature-comparison/). My `search` is
nowhere near as fancy, fast, or feature-riffic as most of those. But it does
exactly what I wanted, and is small and simple. So, you know.

## Usage

Just run `search -h` to get the latest command-line help. Here is an example of
how I frequently use `search`:

```
cd ~/src/some-project
search -x out/ -n '\.(cc|h)$' -c FrobulateGrommets -v
```

This shows me all uses of the word `FrobulateGrommets` in my C++ project, except
in the output directory (which typically contains binaries and generated code).
(It’s important to know precisely **how** and **where** the grommets are being
frobulated.)

## Building And Installing

You’ll need to [install the Go programming
language](https://golang.org/doc/install). If your computer has `make`, you can
just run `make` to build:

```
git clone https://github.com/noncombatant/search.git
cd search
make
```

If not, you can run the steps in the Makefile manually:

```
git clone https://github.com/noncombatant/search.git
cd search
go build
go vet
go test
```

To install, just copy the `search` binary to somewhere in your `$PATH`. I like
to use `$HOME/bin` for these toys. I always have `$HOME/bin` in my `$PATH`, so
this works for me:

```
mkdir -p ~/bin
cp search ~/bin
```

I have only tried this program on Ubuntu and macOS. In theory it should work on
Windows, too. Let me know what happens if you try! I’d take any PRs you’ve got
to ensure it works on Windows. (And any other feature requests?)
