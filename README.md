# passcrux
PassCrux - never lose your soul again!

[![Build Status](https://travis-ci.com/xkortex/passcrux.svg?branch=master "Travis CI status")](https://travis-ci.com/xkortex/passcrux.svg?branch=master)
[![GoDoc](https://godoc.org/github.com/xkortex/passcrux?status.svg)](https://godoc.org/github.com/xkortex/passcrux)
[![Go Report Card](https://goreportcard.com/badge/github.com/xkortex/passcrux)](https://goreportcard.com/report/github.com/xkortex/passcrux)

If you are like me, you loathe single points of failure, *especially* when it comes to passwords, 
and *extra-especially* when it comes to "master" passwords which guard things like password managers. 
Some password managers have mechanism to reset your password, however ever password reset route is an 
increased attack surface. Plus, there is always the spectre of a password manager provider going under, 
or getting hacked, or whatever. I wanted a way to back up my "master" passwords under my own terms. 

PassCrux gets to the crux of this matter, and has only a transient resemblance to crux-sounding recovery 
schemes found in some magical literature ;). It works by separating your password - or any data - into shards, 
that you can do anything you want with, just provide M of the N shards and you can recover the original data. 
If you haven't been living in a cursed cave for the past two decades, you'll recognize this as 
[Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing) algo, which is indeed 
at the heart of this tool. It's basically a lightweight wrapper around SSS with built-in helpers and formatters. 

# Testimonials 

> Heck, that sounds so fun! I can't wait to lose my master password!
 -- [aeksco](https://github.com/aeksco)

> That's a good way to setup a quest for someone 1000 years in the future. Put a few of those in temples guarded by 
> bosses, and you've got yourself a solid adventure.
 -- [erotemic](https://github.com/erotemic)

# Usage

To split a password from a prompt, into 5 shards with 3 needed to reconstruct, enter 
```bash
passcrux split -r 3/5 -p 
```

Out:
```bash
a38f786f19680cb3
c4f5a36d797e336d
38519b5d42021620
f3ab8d463e182893
e4ea839978cc8eae
```

Copy (at least) 3 of the 5 output shards and paste into a file, `shards.txt`, one shard string per line. 

`shards.txt`
```bash
a38f786f19680cb3
38519b5d42021620
e4ea839978cc8eae
```

Then run:
```bash
cat shards.txt | passcrux combine --stdin  
```

# Building/Installation

#### Turbo-instant docker usage:

```docker run --rm -it xkortex/passcrux [OPTIONS]```

Note: `-it` is required for `-p/--pass` (interactive password prompt) and `-v` for any file-I/O. 

#### Conventional:

Requires a typical golang environment. Simply run `make` to compile `passcrux` to `$GOPATH/bin/passcrux`

## todo
- [x] main IO commands
- [x] primary encode formats: hex, base32, base64
- [x] "abc" encoding
- [x] Stabilize pipe-in interface and flags
- [ ] validate/ensure correct behavior with DOS-style carriage returns `\r and \r\n`
- [x] Travis / CI hooks
- [x] Dockerfile
- [ ] standardize output formatter interface
- [ ] config parsing
- [ ] handling for raw bytes I/O

Stretch goals:
- [ ] QR generator
- [ ] QR parser

# License

PassCrux is licensed under [Mozilla Public License v2.0](http://mozilla.org/MPL/2.0/) \[[FAQ](https://www.mozilla.org/en-US/MPL/2.0/FAQ/)\]. Use it to your heart's content. ¯\\\_(ツ)\_/¯ 

The current implementation relies on [SSS](https://github.com/hashicorp/vault/tree/master/shamir) from [Hashicorp Vault](https://github.com/hashicorp/vault), also MPL2. PassCrux uses SSS wholesale without any modification. 

[Cobra](https://github.com/spf13/cobra) is licensed Apache, Copyright (c) 2015 Steve Francia <spf@spf13.com>

[Viper](https://github.com/spf13/viper) is licensed MIT, Copyright (c) 2014 Steve Francia

IANAL and I *think* I am doing this correctly, but if something is in error, please open an Issue or PR to help rectify. Go makes it astoundingly easy to `go get` code and use it in your project, and if you ask me, if you are advertising your code as go-gettable on Github, your intent is to share, but it's not my call, so please clarify if I am in the wrong. 

Thanks!
