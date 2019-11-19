# IPGW Tool
> Email: aUBzaGFuZ3llcy5uZXQ=

## OS Support
- [x] Windows
- [x] OS X
- [x] Linux



## TODO
- [x] Login
- [x] Account Base Info Fetch
- [ ] Account More Info Fetch
- [ ] Logout
- [ ] Kick Off
- [ ] Multi-user Toggle
- [ ] Traffic Monitor
- [ ] Fail Detail

## Usage
```bash
Usage: ipgw [-s] [-u username] [-p password] [-i configPath]
Options:
  -h    show help info
  -i path
        path to configuration file, default is %USER_PROFILE%/.ipgw
  -p password
        password
  -s    save username and password after login successfully
  -u username
        username
  -v    show version and exit
```
### Login

Use `ipgw -u username -p password -s` to login.

if login successfully, account info would be saved into configuration file `%USER_DIR%/.ipgw`,

or you can specify the path of config file by `ipgw -u username -p password -i path -s`.

Then next time, you can just use `ipgw` or `ipgw -i path` to login with the account saved in configuration file.

Or you can use `ipgw -u username -p password` without `-s` to login temporarily.

*Tips: `ipgw` should be changed to `ipgw.exe` on windows*

**It is more convenient to add executable path to environment variable.**
### More to be done