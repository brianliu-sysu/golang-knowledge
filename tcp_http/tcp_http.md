# **TCP/IP Networking Stack**

## The Four-layer Model
```mermaid
flowchart TB
    subgraph TCP/IP["TCP/IP Four-Layer Model"]
        direction TB
        
        APP["<b>Application Layer</b><br/>HTTP, HTTPS, FTP, DNS, SMTP, SSH, WebSocket"]
        
        APP -->|"Data"| TRANS
        
        TRANS["<b>Transport Layer</b><br/>TCP, UDP"]
        
        TRANS -->|"Segment"| NET
        
        NET["<b>Internet Layer</b><br/>IP, ICMP, ARP"]
        
        NET -->|"Packet"| LINK
        
        LINK["<b>Network Interface Layer</b><br/>Ethernet, Wi-Fi, PPP"]
        
        LINK -->|"Frame"| PHY
        
        PHY["<b>Physical Medium</b>"]
    end

    style APP fill:#e1f5fe
    style TRANS fill:#fff3e0
    style NET fill:#f3e5f5
    style LINK fill:#e8f5e9
    style PHY fill:#fce4ec
```

## Layer Responsibilities
|layer|responsibilities|protocol|data unit|
|---|---|---|---|
|application|provide services to applications|HTTP,DNS,FTP|message/data|
|transaction|End-to-End communication|TCP,UDP|Segment|
|Internet|Routing&Addressing|IP, ICMP|Packet|
|Network Interface|Physical Transmission|Ethernet|Frame|

## Data Encapsulation Process
```mermaid
flowchart TB
    subgraph Encapsulation["Data Encapsulation Process"]
        direction TB
        
        A["<b>Application Layer</b><br/>[ Data ]"]
        
        A -->|"Add TCP/UDP Header"| B
        
        B["<b>Transport Layer</b><br/>[ TCP Header | Data ]"]
        
        B -->|"Add IP Header"| C
        
        C["<b>Internet Layer</b><br/>[ IP Header | TCP Header | Data ]"]
        
        C -->|"Add Frame Header & Trailer"| D
        
        D["<b>Network Interface Layer</b><br/>[ Frame Header | IP Header | TCP Header | Data | Frame Trailer ]"]
        
        D -->|"Physical Signal"| E
        
        E["<b>Network Transmission</b><br/>→ Receiver decapsulates layer by layer"]
    end

    style A fill:#e1f5fe
    style B fill:#fff3e0
    style C fill:#f3e5f5
    style D fill:#e8f5e9
    style E fill:#fce4ec
```

## Core Protocol Explain
1. IP protocol(Internet Layer)
IP header (20 bytes)
### IP Header (20 bytes)

```
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |Version|  IHL  |    ToS        |          Total Length         |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |         Identification        |Flags|      Fragment Offset    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  TTL  |    Protocol   |         Header Checksum               |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                       Source IP Address                       |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Destination IP Address                     |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

| Field | Bits | Description |
|-------|------|-------------|
| Version | 4 | IP version (4 for IPv4) |
| IHL | 4 | Header length in 32-bit words |
| ToS | 8 | Type of Service / DSCP |
| Total Length | 16 | Total packet length in bytes |
| Identification | 16 | Fragment identification |
| Flags | 3 | DF (Don't Fragment), MF (More Fragments) |
| Fragment Offset | 13 | Position of fragment in original packet |
| **TTL** | 8 | Time To Live - decremented at each hop |
| **Protocol** | 8 | Upper layer protocol (6=TCP, 17=UDP, 1=ICMP) |
| Header Checksum | 16 | Header error checking |
| Source IP | 32 | Sender's IP address |
| Destination IP | 32 | Receiver's IP address |

2. TCP Protocol(Transaction Layer)
### TCP Header (20 bytes + options)

```
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |          Source Port          |       Destination Port        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                        Sequence Number                        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Acknowledgment Number                      |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |  Data |       |U|A|P|R|S|F|                                   |
   | Offset| Rsrvd |R|C|S|S|Y|I|            Window Size            |
   |  (4)  |  (6)  |G|K|H|T|N|N|                                   |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |           Checksum            |         Urgent Pointer        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                    Options (if any)                           |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

| Field | Bits | Description |
|-------|------|-------------|
| Source Port | 16 | Sender's port number |
| Destination Port | 16 | Receiver's port number |
| Sequence Number | 32 | Byte position in data stream |
| Acknowledgment Number | 32 | Next expected byte from sender |
| Data Offset | 4 | Header length in 32-bit words |
| Reserved | 6 | Reserved for future use |
| Window Size | 16 | Receive window size (flow control) |
| Checksum | 16 | Header + data error checking |
| Urgent Pointer | 16 | Offset to urgent data |

### TCP Flags

| Flag | Name | Description |
|------|------|-------------|
| **SYN** | Synchronize | Initiate connection (three-way handshake) |
| **ACK** | Acknowledgment | Confirms received data |
| **FIN** | Finish | Gracefully close connection |
| **RST** | Reset | Abort connection immediately |
| **PSH** | Push | Deliver data to application immediately |
| **URG** | Urgent | Urgent data present (use Urgent Pointer) |

3. UDP Protocol(Transaction Layer)
### UDP Header (8 bytes)

```
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |          Source Port          |       Destination Port        |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |            Length             |           Checksum            |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

| Field | Bits | Description |
|-------|------|-------------|
| Source Port | 16 | Sender's port number |
| Destination Port | 16 | Receiver's port number |
| Length | 16 | Header + data length in bytes |
| Checksum | 16 | Optional error checking (IPv4) |

### UDP Characteristics

| Feature | Description |
|---------|-------------|
| **Connectionless** | No handshake required, send immediately |
| **Unreliable** | No delivery guarantee, no retransmission |
| **No Congestion Control** | Sends at any rate regardless of network |
| **Low Overhead** | Only 8 bytes header (vs TCP's 20+ bytes) |
| **High Speed** | Ideal for real-time: DNS, VoIP, Gaming, Streaming |

## TCP vs UDP
|feature|TCP|UDP|
|---|---|---|
|connection|connection-oriented|connectionless|
|reliability|reliable connection|unreliable|
|sequence|Guaranteed order|out-of-order|
|Flow Control|Yes|No|
|Congestion Control|Yes|No|
|header size|20+ bytes|8 bytes|
|scenario|HTTP, file transactions|vedio,DNS|

## TCP three-way handshake
```mermaid
sequenceDiagram
    participant C as Client
    participant S as Server

    Note over C,S: Three-way Handshake (Connection Establishment)

    Note over C: CLOSED
    Note over S: LISTEN

    C->>S: SYN, seq=x
    Note over C: SYN_SENT
    Note right of C: Client requests connection

    S->>C: SYN+ACK, seq=y, ack=x+1
    Note over S: SYN_RECEIVED
    Note left of S: Server agrees

    C->>S: ACK, ack=y+1
    Note over C: ESTABLISHED
    Note over S: ESTABLISHED

    Note over C,S: Connection Established
```

## TCP four-way handshake
```mermaid
sequenceDiagram
    participant C as Client
    participant S as Server

    Note over C,S: Four-way Handshake (Connection Termination)

    Note over C: ESTABLISHED
    Note over S: ESTABLISHED

    C->>S: FIN, seq=u
    Note over C: FIN_WAIT_1
    Note right of C: Client initiates close

    S->>C: ACK, ack=u+1
    Note over C: FIN_WAIT_2
    Note over S: CLOSE_WAIT
    Note left of S: Server acknowledges

    Note over S: Server may continue sending data...

    S->>C: FIN, seq=v
    Note over S: LAST_ACK
    Note left of S: Server ready to close

    C->>S: ACK, ack=v+1
    Note over C: TIME_WAIT
    Note over S: CLOSED

    Note over C: Wait 2MSL...
    Note over C: CLOSED

    Note over C,S: Connection Closed
```

## TCP state machine
```mermaid
stateDiagram-v2
    [*] --> CLOSED

    state "Connection Establishment" as establish {
        CLOSED --> LISTEN : Passive open
        CLOSED --> SYN_SENT : Active open, send SYN
        LISTEN --> SYN_RCVD : Receive SYN, send SYN+ACK
        SYN_SENT --> ESTABLISHED : Receive SYN+ACK, send ACK
        SYN_RCVD --> ESTABLISHED : Receive ACK
    }

    state "Active Close" as active_close {
        ESTABLISHED --> FIN_WAIT_1 : Send FIN
        FIN_WAIT_1 --> FIN_WAIT_2 : Receive ACK
        FIN_WAIT_2 --> TIME_WAIT : Receive FIN, send ACK
        TIME_WAIT --> CLOSED : Wait 2MSL
    }

    state "Passive Close" as passive_close {
        ESTABLISHED --> CLOSE_WAIT : Receive FIN, send ACK
        CLOSE_WAIT --> LAST_ACK : Send FIN
        LAST_ACK --> CLOSED : Receive ACK
    }

    note right of TIME_WAIT
        Wait 2MSL (Maximum Segment Lifetime)
        to ensure final ACK is received
    end note
```

## TCP Flow Control

### Purpose
To prevent the sender from transmitting data too fast and overwhelming the receiver.
It is mainly an End-to-End mechanism limited by the receiver's processing and buffering capacity.

### Sliding Window Mechanism

```mermaid
sequenceDiagram
    participant S as Sender
    participant R as Receiver

    Note over S,R: Sliding Window Flow Control

    S->>R: Data (seq=1, 1000 bytes)
    R->>S: ACK=1001, rwnd=3000
    Note left of S: Can send up to 3000 bytes

    S->>R: Data (seq=1001, 1000 bytes)
    S->>R: Data (seq=2001, 1000 bytes)
    R->>S: ACK=3001, rwnd=1000
    Note left of S: Window shrinks, slow down

    S->>R: Data (seq=3001, 1000 bytes)
    R->>S: ACK=4001, rwnd=0
    Note over S: Zero Window! Stop sending

    Note over R: Receiver processes data...

    R->>S: ACK=4001, rwnd=2000
    Note over S: Window opens, resume sending
```

### Key Concepts

| Term | Description |
|------|-------------|
| rwnd | Receiver Window - advertised by receiver in ACK |
| Sliding Window | Sender can send data within window without waiting for ACK |
| Zero Window | rwnd=0, sender must stop and wait |
| Window Probe | Sender periodically checks if rwnd > 0 |

## TCP Congestion Control
### purpose
To manage network congestion and optimize throughput.

###state machine
```mermaid
stateDiagram-v2
    [*] --> SlowStart

    SlowStart --> SlowStart : ACK received
    SlowStart --> CongestionAvoidance : cwnd >= ssthresh
    SlowStart --> FastRecovery : 3 Dup ACKs

    CongestionAvoidance --> CongestionAvoidance : ACK received
    CongestionAvoidance --> FastRecovery : 3 Dup ACKs
    CongestionAvoidance --> SlowStart : Timeout

    FastRecovery --> CongestionAvoidance : New ACK
    FastRecovery --> SlowStart : Timeout

    SlowStart --> SlowStart : Timeout

    note right of SlowStart
        cwnd doubles per RTT
        ssthresh = cwnd/2 on timeout
        cwnd = 1 on timeout
    end note

    note right of CongestionAvoidance
        cwnd += 1 per RTT
    end note

    note right of FastRecovery
        ssthresh = cwnd/2
        cwnd = ssthresh + 3
    end note
```
