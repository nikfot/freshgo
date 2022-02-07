## Freshgo, the golang version manager
https://quasilyte.dev/gopherkon/?state=0c0305100307090203070a000000000000#

![frehs go with Freshgo!](https://github.com/nikfot/freshgo/blob/gh-pages/freshgo.png?raw=true)

Use freshgo to **easily** install and manage your golang versions.

freshgo can install the latest go version, or any given version of go. 

### Get easily started

Freshgo is an easy way to avoid manual actions for installing go and keeping it up to date.
You do not need **git** or any other tool to make this work.

#### Get latest go version
```
freshgo latest
```
#### List available go versions
```
freshgo latest
```
#### Select specific version
```
freshgo select -v 1.17.1
```

### Get **Notified** for latest go version

Freshgo can run everytime your shell opens and check if there is a latest version than the one installed. If there is, it can ask you wether you'd like to update to the latest.

```
echo "${PATH to freshgo repo}/check_ver.sh >>  ~/.profile
```
