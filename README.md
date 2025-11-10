# Turn Server

1. Create `secret/proxy.json` file with proxy account: 
```
    {"username":"your_proxy_email","password":"your_proxy_pwd"}
```
2. Open `service.sh`, change path of `ExecStart`, `WorkingDirectory` to current directory
```
    ExecStart           = /home/my_linux/turn/turn
    WorkingDirectory    = /home/my_linux/turn
```
3. Deploy turn server
```
    $ sudo sh service.sh
```
4. Check status turn server:
```
    $ sudo systemctl status edge-turn.service
```
---
Notice
- Your server deploy turn need to open port from  49152 to 65535
- We had set up a monitor system for turn, it will notify in logging system (Discord) if all turn server shutdown
