## Freshgo, the golang version manager 🔖
- ***lightweight*** go binary 🤸‍♀️
- ***no git*** installation required ⚙️

![fresh go with Freshgo!](https://github.com/nikfot/freshgo/blob/gh-pages/freshgo.png?raw=true)

Use freshgo to **easily** install and manage your golang versions.
Add freshgo in a cron or your .profile to remind you when a new go version is available!

_freshgo_ can install the latest go version, or any given version of go. 

### Download the binary ⬇️
Download the latest binary and get started with Go in no time!

```
wget https://github.com/nikfot/freshgo/blob/gh-pages/freshgo?raw=true -O freshgo && chmod ug+x freshgo && sudo mv freshgo /usr/local/bin/freshgo
```
_This will download the latest binary and install it in your $PATH_


### Contribute 🤝
You are welcome to contribute any time!
Inform me on an issue or make suggestions or add your own repo.

Clone the repo.
```
git clone git@github.com:nikfot/freshgo.git
```
or 
```
git clone git clone https://github.com/nikfot/freshgo
```

If you already have go installed build the binary:
```
cd freshgo && make
```

## Get easily started 🏁

Freshgo is an easy way to avoid manual actions for installing go and keeping it up to date.
You do not need **git** or any other tool to make this work.



#### Get latest go version 📌
```
freshgo latest
```
#### List available go versions 📌
```
freshgo list
```
#### Select specific version 📌
```
freshgo select -v 1.17.7
```

### Get **Notified** for latest go version ⏰

Freshgo can run everytime your shell opens and check if there is a latest version than the one installed. If there is, it can ask you wether you'd like to update to the latest.

```
echo "${PATH to freshgo repo}/check_ver.sh >>  ~/.profile
```
