Wormhole
=========
[![Build Status on TravisCI](https://secure.travis-ci.org/kiesel/wormhole-go.png)](http://travis-ci.org/kiesel/wormhole-go)
[![GitHub release](https://img.shields.io/github/release/kiesel/wormhole-go.svg?maxAge=2592000)](https://github.com/kiesel/wormhole-go/releases)
[![license](https://img.shields.io/github/license/kiesel/wormhole-go.svg?maxAge=2592000)](https://github.com/kiesel/wormhole-go/blob/master/LICENSE)

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
(nohup $HOME/wormhole.exe -quiet 2>&1 &)
```

Client installation
-------------------
You'll need a client, too. Given you're using bash / zsh or a compatible shell, you can use [kiesel/wormhole](https://github.com/kiesel/wormhole).


Security
--------

Wormhole opens a port on a designated interface for you; when binding to public network interfaces, you might expose yourself to serious security risks.

The recommended setup is therefore, to bind it to the loopback address 127.0.0.1 / ::1 and use SSH to make that port available to your client system (which would usually be a local VM). The configuration file `.ssh/config` would look like this:

    Host vbox vb 127.0.0.1
      Hostname 127.0.0.1
      Port 2222

      RemoteForward 127.0.0.1:5115 127.0.0.1:5115
