# path-sanitizer

<!-- toc -->
- [Purpose](#purpose)
- [Examples](#examples)
<!-- /toc -->

## Purpose

`path-sanitizer` is a tiny utility to make configuring `$PATH` a bit easier in shell resource files, such as `~/.profile`, `~/.zprofile` and so on. `path-sanitizer` examines the current `$PATH` value, adds command-line stated directories to it, and sanitizes it so that there are no double entries.

## Examples

Bash:

```shell
# ~/.profile: Add . (dot), $HOME/bin, /usr/local/bin and /usr/local/sbin to $PATH
source <(/where/ever/path-sanitizer -s bash $HOME /usr/local)
#                                                 ^ adds /usr/local/bin if found and /usr/local/sbin
#                                           ^ adds $HOME/bin if it exists, and $HOME/sbin
#                                    ^ intended shell, generates `export PATH=...`
```

Fish:

```shell
# ~/.config/fish/config.fish example
eval "$(/where/ever/path-sanitizer -s fish $HOME /usr/local)"
#                                     ^ intended shell, generates `set -gx ...`
```

Wheter to use `source` or `eval` depends on your shell. Fortunately, `path-sanitizer -h` shows examples.