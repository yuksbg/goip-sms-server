package main

import (
	"encoding/csv"
	"fmt"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
)

// GoIPMessage represents a parsed GoIP message for RECEIVE
type GoIPMessage struct {
	ID         string
	Password   string
	SrcNum     string
	Message    string
	RemoteAddr string
}

func main() {
	address := "0.0.0.0:44444"

	// Set up logrus for structured logging
	log := logrus.New()

	// Create a UDP address to listen on
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.WithError(err).Fatal("ResolveUDPAddr failed")
	}

	// Create a UDP socket
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.WithError(err).Fatal("ListenUDP failed")
	}
	defer conn.Close()

	log.Info("Start listening for GOIP messages")

	buffer := make([]byte, 2048)
	for {
		// Read incoming data
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.WithError(err).Error("Error reading from UDP")
			continue
		}

		message := string(buffer[:n])
		log.WithFields(logrus.Fields{
			"remote_address": addr.String(),
			"message":        message,
		}).Info("Received message")

		// Handle the message
		if strings.HasPrefix(message, "RECEIVE:") {
			parsedMessage, err := parseGoIPMessage(message, addr.String())
			if err != nil {
				log.WithError(err).Error("Error parsing RECEIVE message")
				continue
			}
			log.WithFields(logrus.Fields{
				"id":          parsedMessage.ID,
				"srcnum":      parsedMessage.SrcNum,
				"message":     parsedMessage.Message,
				"remote_addr": parsedMessage.RemoteAddr,
			}).Info("Parsed RECEIVE message")
		} else if strings.HasPrefix(message, "req") {
			// Convert the message from colon-separated to semicolon-separated
			csvMessage := strings.ReplaceAll(message, ":", ";")
			reader := csv.NewReader(strings.NewReader(csvMessage))
			reader.Comma = ';'
			messageArray, err := reader.Read()
			if err != nil {
				log.WithError(err).Error("Error parsing CSV")
				continue
			}

			// Construct and send the response for keep-alive requests
			response := fmt.Sprintf("reg:%s;status:200;", messageArray[1])
			_, err = conn.WriteToUDP([]byte(response), addr)
			if err != nil {
				log.WithError(err).Error("Error sending response")
				continue
			}
			log.WithFields(logrus.Fields{
				"response": response,
			}).Info("Sent response")
		}
	}
}

// parseGoIPMessage parses a RECEIVE message and returns a GoIPMessage struct
func parseGoIPMessage(rawMessage, remoteAddr string) (*GoIPMessage, error) {
	// Removing "RECEIVE:" part
	parts := strings.SplitN(rawMessage, "RECEIVE:", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid message format")
	}

	// The main data part of the message
	dataPart := parts[1]

	// Splitting by semicolons to get individual fields
	fields := strings.Split(dataPart, ";")
	message := &GoIPMessage{
		RemoteAddr: remoteAddr,
	}

	// Iterating through the fields to populate the struct
	for _, field := range fields {
		parts := strings.Split(field, ":")
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "id":
			message.ID = value
		case "password":
			message.Password = value
		case "srcnum":
			message.SrcNum = value
		case "msg":
			message.Message = value
		}
	}

	return message, nil
}
