## Gmail checker for conky
Create gmail checker written on go

Put your nickname and password in *```~/.gmail.json```* such way
```json
{"account":"ACCOUNT","short_conky":"SHORT","username":"USERNAME","password":"PASSWORD"}
```
> You must use first part of your email like this

> **USERNAME**@gmail.com


account is label for account

short_conky - this text will show in conky

if *```~/.gmail.json```* is not exist it will be created with example values

# Install
You need installed go.

Just clone the repository and run
```
go build gmail.go
```
You get binary file gmail. You can put it to /usr/local/bin and run

First time you run it, it create config file in *```~/.gmail.json```

After this just put your data into config and create command for conky
### Issues
Add multiple accounts

#License
MIT License. Copyright (c) 2013-2015 Vitaliy Drevenchuk.
