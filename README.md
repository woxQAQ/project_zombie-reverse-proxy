# Project Zombie Reverse Proxy

## Prepare

- machine with public IP as proxy machine(the more bandwidth the better)
- machine as the service
- (option) machine as client

## How to use

1. install frp in your proxy machine and run 

```shell
cd proxy
frps -c frps.toml
```

2. install frp in your service machine and run

```
cd service 
frpc -c frpc.toml
```

3a. (windows) install frp in your client and run 
``` powershell
cd client
frpc -c fprc.toml
```
or run the bat

```
cd client 
./run.bat
```


