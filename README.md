# Go MobileQQ API

> An open-source RPC-based fast MobileQQ client API in go

## Try MobileQQ client API

Waiting for `v0.1.0-pre`, or you can build it from source.

```bash
# build from source
go get github.com/elap5e/go-mobileqq-api/cmd/go-mobileqq-echo@latest

# generate config.yaml template
go-mobileqq-echo

# add your account into config.yaml
vi ~/.goqq-dev/config.yaml

# enjoy it!
go-mobileqq-echo
```

> **NOTE:** Unstable version, you may lose everything if you **_DO NOT KNOW_** what you did!

### Auth

1. Captcha **[recommand]** (click twice! open the link and drag the slider)
2. Picture verification
3. SMS Code verification
4. Auto-processing auth response
5. Auth without password (signature auth)

### Message

1. Receiving and sending messages
2. Markdown message decoding

## TODO

> just a plan here

### Release `v0.1.0-pre`

1. `Markdown` messages full support
2. Multi-users login

### Future Plans

1. Use local database (Release `v0.1.0`)
2. Support more protocols (Release `v0.1.0`)
3. `Telegram Bot API`

## License

This project is licensed under the GNU Affero General Public License version 3.0.
