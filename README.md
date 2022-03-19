# tftp

A toy tftp server written in Golang.

## Four types of packet

```
          2 bytes    string   1 byte     string   1 byte
          -----------------------------------------------
   RRQ/  | 01/02 |  Filename  |   0  |    Mode    |   0  |
   WRQ    -----------------------------------------------
   
          2 bytes    2 bytes       n bytes
          ---------------------------------
   DATA  | 03    |   Block #  |    Data    |
          ---------------------------------
          
          2 bytes    2 bytes
          -------------------
   ACK   | 04    |   Block #  |
          --------------------
          
          2 bytes  2 bytes        string    1 byte
          ----------------------------------------
   ERROR | 05    |  ErrorCode |   ErrMsg   |   0  |
          ----------------------------------------
```

## Key process

### Host A reads a file from Host B

- Host A sends RRQ to Host B
- Host B sends DATA to Host A
- Host A sends ACK to Host B
- ...
- Host B sends DATA with data less than 512 bytes
- Host A sens ACK to Host B

### Host B writes a file to Host B

- Host A sends WRQ to Host B
- Host B sends ACK to Host A
- Host A sends DATA to Host B
- Host B sends ACK to Host A
- ...
- Host A sens DATA with data less than 512 bytes
- Host B sends ACK to Host A

## Run and test

```
go run cmd/main.go
```

Test on Mac:

```
$ tftp 127.0.0.1 69
tftp> get 123.png
Received 102490 bytes in 0.8 secondsk

tftp> get 123.png
Received 102490 bytes in 0.2 seconds 

tftp> quit
```

## Reference
- [RFC 1350](https://datatracker.ietf.org/doc/html/rfc1350)
- [Trivial File Transfer Protocol (TFTP)](http://tcpip.marcolavoie.ca/tftp.html)
