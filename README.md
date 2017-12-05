QED - QUIC with enhanced distribution
=====================================
QED is a experimental multipath extension for the QUIC protocol. It is based on minq (https://www.github.com/ekr/minq) and comes with all its features and flaws.


## WARNING

QED is absolutely not suitable for any kind of production use. Is is an academic research project and will most likely break at this point.


## Logging

To enable logging, set the ```MINQ_LOG``` environment variable, as
in ```MINQ_LOG=connection go test```. Valid values are:

    // Pre-defined log types
    const (
    	logTypeAead       = "aead"
    	logTypeCodec      = "codec"
    	logTypeConnBuffer = "connbuffer"
    	logTypeConnection = "connection"
    	logTypeAck        = "ack"
    	logTypeFrame      = "frame"
    	logTypeHandshake  = "handshake"
    	logTypeTls        = "tls"
    	logTypeTrace      = "trace"
    	logTypeServer     = "server"
    	logTypeUdp        = "udp"
    	logTypeMultipath  = "mp"
    )

Multiple log levels can be separated by commas.

## Minq

QED depends on Minq (https://www.github.com/ekr/minq) for QUIC.

## Mint

Minq depends on Mint (https://www.github.com/bifurcation/mint) for TLS.
Currently Mint master should work, but occasionally I will have to be on
a branch. Will try to keep this updated.
