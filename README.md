
# GoIP SMS Server

This repository contains a simple SMS server built using Golang that interacts with a **GoIP GSM Gateway** device. The server listens for incoming messages from the GoIP device over UDP and processes SMS messages. It also responds to keep-alive requests from the GoIP device, ensuring the connection is maintained.

## About GoIP Devices

**GoIP** devices are GSM gateways that enable communication between the GSM and IP networks. These devices are widely used for applications such as:

- Sending and receiving SMS messages.
- Making and receiving calls via GSM using an IP-based network (VoIP).
- Providing gateway functionality in call centers and SMS gateways.

This SMS server is designed to interface with a GoIP device, handling SMS messages sent to the GSM gateway and processing them in real-time.


## Requirements

- Go 1.18 or later
- A **GoIP device** properly configured to communicate with this server over UDP.
- A working GSM SIM card installed in the GoIP device.

## Setup Instructions

Ensure your GoIP device is configured to send SMS messages to the server’s IP address and port. The server listens on UDP port `44444` by default, but you can modify this in the code.

### Build and Run the Server

To build and run the server:

```bash
go build -o goip-sms-server .
./goip-sms-server
```


### Sample Output

Upon receiving an SMS message, you will see structured logs in the terminal. For example:

```
INFO[0001] Start listening for GOIP messages
INFO[0005] Received message    remote_address=192.168.1.50:5060 message="RECEIVE:id:1;password:secret;srcnum:+1234567890;msg:Hello, World!"
INFO[0005] Parsed RECEIVE message    id=1 srcnum="+1234567890" message="Hello, World!" remote_addr="192.168.1.50:5060"
```


### Keep-Alive Responses

The server will also handle keep-alive (`req`) messages from the GoIP device. The device sends periodic requests to check if the server is online. The server responds with a status code `200`, indicating the connection is active.

### Modify the Server Address

If you need to change the server's listening address or port, modify the `address` variable in the `main()` function:

```go
address := "0.0.0.0:44444" // Change this to your preferred IP and port
```

### Logs and Debugging

This server uses `logrus` for logging, providing detailed logs for each event. You can modify the logging format and level by editing the logging setup in the `main()` function.


## Contributing

Feel free to submit issues or pull requests for any improvements or features you'd like to add.
