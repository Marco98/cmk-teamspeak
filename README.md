# cmk-teamspeak

Check_MK agent check for Teamspeak3 virtual server instances

## Installation (Raw)

### Checkmk site server

```shell
omd su SITE
cd /tmp
wget https://github.com/Marco98/cmk-teamspeak/releases/download/v0.3.0/teamspeak-v0.3.0.mkp
mkp install teamspeak-v0.3.0.mkp
```

### Host with checkmk agent

```shell
cd /usr/lib/check_mk_agent/local
wget https://github.com/Marco98/cmk-teamspeak/releases/download/v0.3.0/Teamspeak3
chmod +x Teamspeak3
```
