## Freshgo, the golang version manager

![fresh go with Freshgo!](https://github.com/nikfot/freshgo/blob/gh-pages/freshgo.png?raw=true)

Use freshgo to **easily** install and manage your golang versions.

freshgo can install the latest go version, or any given version of go. 

### Clone the repo
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
### Download the binary
```
wget wget https://github.com/nikfot/freshgo/blob/gh-pages/freshgo?raw=true > freshgo && chmod ug+x freshgo
```

## Get easily started

Freshgo is an easy way to avoid manual actions for installing go and keeping it up to date.
You do not need **git** or any other tool to make this work.



#### Get latest go version
```
freshgo latest
```
#### List available go versions
```
freshgo list
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
