# passcrux
PassCrux - Immortalize your master password!

[![Build Status](https://travis-ci.com/xkortex/passcrux.svg?branch=master)](https://travis-ci.com/xkortex/passcrux)
[![GoDoc](https://godoc.org/github.com/xkortex/passcrux?status.svg)](https://godoc.org/github.com/xkortex/passcrux)
[![Go Report Card](https://goreportcard.com/badge/github.com/xkortex/passcrux)](https://goreportcard.com/report/github.com/xkortex/passcrux)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fxkortex%2Fpasscrux.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fxkortex%2Fpasscrux?ref=badge_shield)

If you are like me, you loathe single points of failure, *especially* when it comes to passwords, 
and **extra-especially** when it comes to "master" passwords which guard things like password managers. 
Some password managers have mechanisms to reset your password, however every password reset-function is an 
increased attack surface. Plus, there is always the spectre of a password manager provider going under, 
or getting hacked, or whatever. I wanted a way to back up my "master" passwords under my own terms. 

PassCrux gets to the crux of this matter, and has only a transient resemblance to crux-sounding recovery 
schemes found in some magical literature ;). It works by separating your password - or any data - into shards, 
that you can do anything you want with. Just provide `M` of the `N` shards and you can recover the original data. 
If you haven't been living in a cursed cave for the past two decades, you'll recognize this as 
[Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing) algorithm, which is indeed 
at the heart of this tool. It's basically a lightweight wrapper around SSS, with built-in helpers and formatters
for storing the shards in human-readable format. 

## Testimonials 

> Heck, that sounds so fun! I can't wait to lose my master password!
 -- [aeksco](https://github.com/aeksco)

> That's a good way to setup a quest for someone 1000 years in the future. Put a few of those in temples guarded by 
> bosses, and you've got yourself a solid adventure.
 -- [erotemic](https://github.com/erotemic)

## Usage

To split a password from a prompt, into 5 shards with 3 needed to reconstruct, enter 
```bash
passcrux split --ratio 3/5 --prompt 
```

Out:
```bash
a38f786f19680cb3
c4f5a36d797e336d
38519b5d42021620
f3ab8d463e182893
e4ea839978cc8eae
```

Now, distribute your shards. In this case, we are just going to dump them into a file for this demo. 
Copy (at least) 3 of the 5 output shards and paste into a file, `shards.txt`, one shard string per line. 

`shards.txt`
```bash
a38f786f19680cb3
38519b5d42021620
e4ea839978cc8eae
```

Then run this command to print the secret:
```bash
cat shards.txt | passcrux combine
```

## What to do with shards?
That's entirely up to you! Get creative with it! Here are some ideas to get your ideas churning:
- Stamp them into metal. Hide them in various places. Make a good ol' fashioned pirate map 🏴‍☠️
- Use `--enc abc` and mark the letters in your favorite books 📚
- Convert hexadecimal values into notes and make some sweet guitar riffs 🎸
- Give a copy to `M` trusted friends with instructions to delete your browser history 🗑️


## Building/Installation

#### Turbo-instant docker usage:

```docker run --rm -it xkortex/passcrux [OPTIONS]```

Note: `-it` is required for `-p/--prompt` (interactive password prompt) and `-v` for any file-I/O. 

#### Conventional:

Requires a typical golang environment. Simply run `make` to compile `passcrux` to `$GOPATH/bin/passcrux`

# FAQ

### What is the threat model exactly? 

Fallibility and data loss, both digitally and mentally. Modern security practices encourage creating a security bottleneck, using random, unmemorable, high-strength passwords for day-to-day logins, and using some sort of vault or password manager to store those keys. However, you still need to secure the key/password for the manager. This means there is usually at least one Super Ultra Master Password To Rule Them All. You probably shouldn't just write it down, or save it in plain text on a drive, since that makes it trivial for an attacker to recover your entire vault. 

I wanted a system that would allow me to securely back up these Master Passwords, but not in a way that relied on memorizing some sequence of tokens. Humans are generally bad at memorizing long sequences of tokens, even when allowing for the usual tricks (was it correct-staple-horse-battery or correct-horse-battery-staple? Or was it with underscores?). Furthermore, backups are usually not frequently accessed, increasing the risk of forgetting. I needed a mental model which a) provided redundancy and b) leveraged the types of things that humans *are* good at memorizing: meaningful locations and events. A real-world [memory palace/ method of loci](https://en.wikipedia.org/wiki/Method_of_loci). Much like That Guy Who Ought Not Have His Name Mentioned, Passcrux tokens can be tied to something (or someone) personally meaningful. I'm sure you remember where you had your first kiss, or got your first puppy, or your favorite book - maybe even favorite passage. Chosen wisely, these keypoints are very easy for you to recall or access, but extremely improbably for an attacker to recover. And it's 100% offline. 

Passcrux also enables you to easily set up some sort of a trust, as a contingency for your accounts in case of an unexpected event that leaves you incapacitated. Work it into a scavenger hunt if you are the mischievous type. 

### Does this have anything to do with "password horcruxing"?

There is a neat blog post going around called [Double Blind Passwords (aka Horcruxing)](https://kaizoku.dev/double-blind-passwords-aka-horcruxing). This project is wholly unrelated. In fact, this method is actually the *opposite* goal of horcruxing (depending on your perspective). Horcruxes are a form of redundancy - That Bloke Who Won't Be Named For Copyright Reasons split his soul up, such that in the event of his death, he'd have some anchors back to the physical world, such that he could be reanimated from the future. Canon does not say the number of horcruxes requried, but the destruction of the diary shows there is some redundancy, and the plot points suggest that he is using a `1 of N` scheme. Phani Karan's technique of blinded passwords, while really neat, is fail-open: loss of the secret chunk results in data loss. 

### How stable is this program? 

I am aiming for extremely high stability, and so far has proven to be, however I make no guarantees. I wrote this foremost for myself, for use with very-long-term storage, where I may not have access to the same physical systems as I did when I wrote it. It needs to be easy to compile and unchanging with its functionality. I plan on vendoring SSS at some point, which should further mitigate against any changes. 

## todo
- [x] main IO commands
- [x] primary encode formats: hex, base32, base64
- [x] "abc" encoding
- [x] Stabilize pipe-in interface and flags
- [ ] validate/ensure correct behavior with DOS-style carriage returns `\r and \r\n`
- [ ] test on windows
- [x] Travis / CI hooks
- [x] Dockerfile
- [ ] standardize output formatter interface
- [ ] goexpect for testing interactive password prompt
- [ ] config parsing
- [ ] handling for raw bytes I/O

Stretch goals:
- [ ] QR generator
- [ ] QR parser
- [ ] built-in error correcting codes

# License

PassCrux is licensed under [Mozilla Public License v2.0](http://mozilla.org/MPL/2.0/) \[[FAQ](https://www.mozilla.org/en-US/MPL/2.0/FAQ/)\]. Use it to your heart's content. ¯\\\_(ツ)\_/¯ 

The current implementation relies on [SSS](https://github.com/hashicorp/vault/tree/master/shamir) from [Hashicorp Vault](https://github.com/hashicorp/vault), also MPL2. PassCrux uses SSS wholesale without any modification. 

[Cobra](https://github.com/spf13/cobra) is licensed Apache, Copyright (c) 2015 Steve Francia <spf@spf13.com>

[Viper](https://github.com/spf13/viper) is licensed MIT, Copyright (c) 2014 Steve Francia

IANAL and I *think* I am doing this correctly, but if something is in error, please open an Issue or PR to help rectify. Go makes it astoundingly easy to `go get` code and use it in your project, and if you ask me, if you are advertising your code as go-gettable on Github, your intent is to share, but it's not my call, so please clarify if I am in the wrong. 

Thanks!


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fxkortex%2Fpasscrux.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fxkortex%2Fpasscrux?ref=badge_large)
