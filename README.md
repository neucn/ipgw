# IPGW Tool
> Email: aUBzaGFuZ3llcy5uZXQ=

## OS Support
- [x] Windows
- [x] OS X
- [x] Linux



## TODO
- [x] Login
- [x] Account Base Info Fetch
- [x] Logout
- [x] Kick out
- [ ] Account More Info Fetch
- [ ] Multi-user Toggle
- [ ] Traffic Monitor
- [ ] Fail Detail

## Document
### Begin to use
#### Linux or OS X
1.
    a. download source code and build

    `git clone https://github.com/iMyOwn/ipgw && cd ipgw && go build`
    
    b. or you can download releases in `https://github.com/iMyOwn/ipgw/releases`

    `wget https://github.com/iMyOwn/ipgw/releases/download/(here is version)/ipgw_(here is os)`

2. set permission and move `./ipgw` to one of paths in `$PATH`

    `chmod +x ipgw && mv ipgw path/to/path`
   
3. then you can use `ipgw` anywhere

#### Windows
1.
    a. download source code and build

    `git clone https://github.com/iMyOwn/ipgw && cd ipgw && go build`
    
    b. or you can download releases `ipgw.exe` in `https://github.com/iMyOwn/ipgw/releases`

2. move the `ipgw.exe` into any path in environment variable or add it's path to environment variable

3. then you can use `ipgw` anywhere

### Usage
```
Usage: ipgw [-s] [-u username] [-p password] [-i configPath]
Options:
  -h    show help info
  -i path
        path to configuration file, default is %USER_PROFILE%/.ipgw
  -k sid
        kick out the specific device by sid
  -o    log out
  -p password
        password
  -s    save username and password after login successfully
  -u username
        username
  -v    show version and exit
```

### Login

- Use `ipgw -u username -p password -s` when login account for the first time.

    - if login successfully, username and password would be saved into configuration file `%USER_DIR%/.ipgw`,

- or you can specify the path of config file by `ipgw -u username -p password -i path/to/config -s`.

- Then next time, you can just use `ipgw` or `ipgw -i path/to/config` to login with the account saved in configuration file.


**For temporarily use**: you can use `ipgw -u username -p password` without `-s` to login without saving username and password.

### Logout

- Use `ipgw -u username -p password -o` to logout

- If you have saved username and password then you can just use `ipgw -o` or `ipgw -o -i path/to/config` to logout

### Kick out

- Use `ipgw -k sid` to kick out specify device

- WARNING: It can be used to kick out everyone's device without any validation, the user shall be responsible for any consequences

### More to be done