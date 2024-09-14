# Mail checker for unread messages

### For i3status-rust

I use i3/Sway wm and want to have a notification for new emails in i3 status line.

Without **$XDG_CONFIG_HOME/mail_checker/config.json** it will fail.

In order to create `config.json` file with this structure.

```bash
$ bin/mail -add
```

This command will create sample config file with this content

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
Also you need to put `<client id>.json` file, with your [credentials](https://developers.google.com/identity/protocols/oauth2/native-app) to the same directory as `config.json`.
Unread count is available as `dbus` message.

# Install
You need installed go.

Just clone the repository and run

```bash
go build -o ./bin/mail ./cmd/main.go
```

Or use go-task for this

```bash
task build
```

You get binary file mail. You can put it to /usr/local/bin and run

You can use multiple accounts, just make sure `$XDG_CONFIG_HOME/mail_checker/config.json` has valid json format
