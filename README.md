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

### Issues
Add multiple accounts

#License
MIT License. Copyright (c) 2013-2015 Vitaliy Drevenchuk.
