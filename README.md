# passcrux
PassCrux - never lose your soul again!

If you are like me, you loathe single points of failure, *especially* when it comes to passwords, and *extra-especially* when it comes to "master" passwords which guard things like password managers. Some password managers have mechanism to reset your password, however ever password reset route is an increased attack surface. Plus, there is always the spectre of a password manager provider going under, or getting hacked, or whatever. I wanted a way to back up my "master" passwords under my own terms. 

PassCrux gets to the crux of this matter, and has only a transient resemblance to crux-sounding recovery schemes found in some magical literature ;). It works by separating your password - or any data - into shards, that you can do anything you want with, just provide M of the N shards and you can recover the original data. If you haven't been living in a cursed cave for the past two decades, you'll recognize this as [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing) algo, which is indeed at the heart of this tool. It's basically a lightweight wrapper around SSS with built-in helpers and formatters. 

## todo
- main IO commands
- output formatter interface
- config parsing
- QR generator
- QR parser?

# License

PassCrux is licensed under [Mozilla Public License v2.0](http://mozilla.org/MPL/2.0/) \[[FAQ](https://www.mozilla.org/en-US/MPL/2.0/FAQ/)\]. Use it to your heart's content. ¯\\\_(ツ)\_/¯ 

The current implementation relies on [SSS](https://github.com/hashicorp/vault/tree/master/shamir) from [Hashicorp Vault](https://github.com/hashicorp/vault), also MPL2. PassCrux uses SSS wholesale without any modification. 

[Cobra](https://github.com/spf13/cobra) is licensed Apache, Copyright (c) 2015 Steve Francia <spf@spf13.com>

[Viper](https://github.com/spf13/viper) is licensed MIT, Copyright (c) 2014 Steve Francia

IANAL and I *think* I am doing this correctly, but if something is in error, please open an Issue or PR to help rectify. Go makes it astoundingly easy to `go get` code and use it in your project, and if you ask me, if you are advertising your code as go-gettable on Github, your intent is to share, but it's not my call, so please clarify if I am in the wrong. 

Thanks!
