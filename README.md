# Go MobileQQ API

> An open-source RPC-based fast MobileQQ client API in go

## Try MobileQQ client API

Waiting for `v0.1.0-alpha`, or you can build it from source.

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

### Account

1. Update status

### Auth

1. Check captcha (click twice! open the link and drag the slider) **[recommand]**
2. Check picture (integrated with iTerm2)
3. Check SMS (confirm to refresh SMS)
4. Password/Non-Password sign in
5. Auto-processing auth response
6. Auto-unlocking device

### Message

1. Delete message
2. Get message
3. Send message
4. Handle online push message
5. Handle push notify
6. Handle push readed (not fully implement)

## Packages

### Crypto

1. ECDH key exchange
2. TEA cipher encrypt/decrypt

### Encoding

1. JCE Marshal/Unmarshal
2. Markdown (with emoticon, not fully implement)
3. OICQ Marshal/Unmarshal
4. UNI Marshal/Unmarshal

### RPC

1. Codec interface
2. Heartbeat alive
3. Server notify
4. TCP dialing test

### Others

1. Bytes (not fully implement)
2. Protobuf (not fully implement)
3. TLVs (not fully implement)

## TODO

> just a plan here

1. Full support `Markdown` messages (Release `v0.1.0-alpha`)
2. Multi-users login (Release `v0.1.0-alpha`)
3. Use local database (Release `v0.1.0`)
4. Support more protocols (Release `v0.1.0`)
5. `Telegram Bot API`

## License

This project is licensed under the GNU Affero General Public License version 3.0.
