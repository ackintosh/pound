# pound

A dummy POP3 server.


## Usage

```
$ mkdir .pound
$ echo foobar > .pound/maildata
$ pound
# dummy pop3 server has started.
$ telnet localhost {port}
# authentication is not required.
> list
+OK 1 messages (1 octets)
> retr 1
+OK 1 octets
foobar
.
> quit
```

## Author

Akihito Nakano