# cmk-teamspeak

Check_MK agent check for Teamspeak3 virtual server instances

## Installation (Raw)

### Checkmk site server

```shell
omd su SITE
cd /tmp
wget https://github.com/Marco98/cmk-teamspeak/releases/download/v0.3.1/teamspeak-v0.3.1.mkp
mkp install teamspeak-v0.3.1.mkp
```

### Host with checkmk agent

**/etc/check_mk/teamspeak3.cfg**:

```ini
[serverquery]
address = "127.0.0.1"
user = serveradmin
password = pass
```

```shell
cd /usr/lib/check_mk_agent/local
wget https://github.com/Marco98/cmk-teamspeak/releases/download/v0.3.1/Teamspeak3
chmod +x Teamspeak3
```
