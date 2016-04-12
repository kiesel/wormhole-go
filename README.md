Wormhole
=========

Wormhole is a an application that allows to open files from a commandline within a VM in your favorite editor(s) / applications in the host system.

Prerequisites are that:

* the part of the VM filesystem that is hosting the files in question is mounted in your host OS
* you are logging in via SSH (though that limitation is only relevant for the client part.)

Installation
------------

1. Download the latest release from the [GitHub Releases](https://github.com/kiesel/wormhole-go/releases) page.
2. Extract `.wormhole.yml` from the release zip into your home directory.
3. Run `wormhole` / `wormhole.exe`

To start wormhole with your shell, put this line into `.bashrc` / `.zshrc`:

```sh
(nohup $HOME/wormhole.exe >>wormhole.log 2>&1 &)
```
