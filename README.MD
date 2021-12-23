# Raw packet encode and decode in go

Raw TCP packet serialize and deserialize without adding byte padding. Support most of available type in go.

Packet header is 6 byte:
```
Packet length: 2 first byte (uint16)
Packet encrypted?: 1 byte (bool)
Packet compressed?: 1 byte (bool)
Packet type: 2 byte (uint16)
```

Check `packet_test.go` for example.

