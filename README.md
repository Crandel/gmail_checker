# Mail checker for unread messages

### For i3status-rust

I use i3/Sway wm and want to have a notification for new emails in i3 status line.

## Requirements
You need to create [OAuth client ID credentials](https://developers.google.com/workspace/guides/create-credentials#desktop-app) and

In order to create `$XDG_CONFIG_HOME/mail/config.json` file run these commands:

```bash
$ go build -o ./bin/mail ./cmd/main.go
$ bin/mail -add
```

Answering few questions this command will create sample config file with this content:

```json
[
    {
        "client_id":     "<client_id>",
        "mail_type":     "gmail",
        "email":         "<email address>",
        "short":         "A",
    }
]
```

Just edit this file.
You could use several gmail accounts to have a personal and work notifications.
Before adding new account you need to have `<client id>.json` file, with your [credentials](https://developers.google.com/identity/protocols/oauth2/native-app) and specify the path to this file during account creation.
Unread count is available as `dbus` message.

# Install
You need to have an installed go.

Just clone the repository and run

```bash
go build -o ./bin/mail ./cmd/main.go
```

Or use go-task for this

```bash
task build
```

You get binary file mail. You can put it to /usr/local/bin and run

You can use multiple accounts, just make sure `$XDG_CONFIG_HOME/mail/config.json` has valid json format
