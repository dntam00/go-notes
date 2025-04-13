package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
)

// IPsecPacket represents a simplified IPsec encapsulated packet
type IPsecPacket struct {
	SPI           uint32 // Security Parameters Index
	SequenceNum   uint32 // Sequence Number
	IV            []byte // Initialization Vector
	EncryptedData []byte // Encrypted payload
	ICV           []byte // Integrity Check Value
}

// SecurityAssociation represents an IPsec SA
type SecurityAssociation struct {
	SPI         uint32
	SourceIP    net.IP
	DestIP      net.IP
	EncryptKey  []byte
	AuthKey     []byte
	EncryptAlgo string
	AuthAlgo    string
	Mode        string // "transport" or "tunnel"
}

// EncapsulateESP encapsulates data using ESP (Encapsulating Security Payload)
func EncapsulateESP(sa *SecurityAssociation, originalPacket []byte) (*IPsecPacket, error) {
	// Create a new IPsec packet
	packet := &IPsecPacket{
		SPI:         sa.SPI,
		SequenceNum: generateSequenceNumber(),
	}

	// Generate IV for AES-CBC
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}
	packet.IV = iv

	// Encrypt the payload
	encryptedData, err := encryptPayload(originalPacket, sa.EncryptKey, iv, sa.EncryptAlgo)
	if err != nil {
		return nil, err
	}
	packet.EncryptedData = encryptedData

	// Calculate integrity check value (simplified)
	packet.ICV = calculateICV(packet, sa.AuthKey, sa.AuthAlgo)

	return packet, nil
}

// DecapsulateESP decapsulates an ESP packet
func DecapsulateESP(sa *SecurityAssociation, packet *IPsecPacket) ([]byte, error) {
	// Verify integrity check value
	calculatedICV := calculateICV(packet, sa.AuthKey, sa.AuthAlgo)
	if !compareICV(packet.ICV, calculatedICV) {
		return nil, fmt.Errorf("integrity check failed")
	}

	// Decrypt the payload
	decryptedData, err := decryptPayload(packet.EncryptedData, sa.EncryptKey, packet.IV, sa.EncryptAlgo)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

// Helper functions
func generateSequenceNumber() uint32 {
	// In a real implementation, this would be tracked per SA
	var seqNum [4]byte
	rand.Read(seqNum[:])
	return uint32(seqNum[0])<<24 | uint32(seqNum[1])<<16 | uint32(seqNum[2])<<8 | uint32(seqNum[3])
}

func encryptPayload(data, key, iv []byte, algo string) ([]byte, error) {
	switch algo {
	case "aes-cbc":
		// Create AES cipher
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}

		// Pad the data to match block size (simplified)
		paddedData := padPKCS7(data, aes.BlockSize)

		// Encrypt
		ciphertext := make([]byte, len(paddedData))
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(ciphertext, paddedData)

		return ciphertext, nil
	default:
		return nil, fmt.Errorf("unsupported encryption algorithm: %s", algo)
	}
}

func decryptPayload(ciphertext, key, iv []byte, algo string) ([]byte, error) {
	switch algo {
	case "aes-cbc":
		// Create AES cipher
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}

		// Decrypt
		plaintext := make([]byte, len(ciphertext))
		mode := cipher.NewCBCDecrypter(block, iv)
		mode.CryptBlocks(plaintext, ciphertext)

		// Unpad
		unpaddedData := unpadPKCS7(plaintext)
		return unpaddedData, nil
	default:
		return nil, fmt.Errorf("unsupported encryption algorithm: %s", algo)
	}
}

func calculateICV(packet *IPsecPacket, key []byte, algo string) []byte {
	// Simplified ICV calculation - in real implementations,
	// this would use HMAC-SHA1, HMAC-SHA256, etc.
	icv := make([]byte, 16)
	// Simulate ICV calculation
	for i := 0; i < 16; i++ {
		icv[i] = byte(i) ^ byte(packet.SPI>>uint(i*2))
	}
	return icv
}

func compareICV(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func padPKCS7(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

func unpadPKCS7(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return nil
	}
	padding := int(data[length-1])
	return data[:length-padding]
}

func dumpPacket(packet *IPsecPacket) {
	fmt.Printf("IPsec Packet:\n")
	fmt.Printf("  SPI: %08x\n", packet.SPI)
	fmt.Printf("  Sequence Number: %d\n", packet.SequenceNum)
	fmt.Printf("  IV: %s\n", hex.EncodeToString(packet.IV))
	fmt.Printf("  Encrypted Data (%d bytes): %s...\n",
		len(packet.EncryptedData),
		hex.EncodeToString(packet.EncryptedData[:min(16, len(packet.EncryptedData))]))
	fmt.Printf("  ICV: %s\n", hex.EncodeToString(packet.ICV))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Example usage
	// Create a security association
	sa := &SecurityAssociation{
		SPI:         0x12345678,
		SourceIP:    net.ParseIP("192.168.1.1"),
		DestIP:      net.ParseIP("192.168.1.2"),
		EncryptKey:  []byte("0123456789abcdef0123456789abcdef"), // 32-byte key for AES-256
		AuthKey:     []byte("authenticationkey12345"),
		EncryptAlgo: "aes-cbc",
		AuthAlgo:    "hmac-sha1",
		Mode:        "transport",
	}

	// Original IP packet (simplified)
	originalPacket := []byte("This is a test IP packet payload that needs to be protected with IPsec")
	fmt.Printf("Original packet (%d bytes): %s\n\n", len(originalPacket), originalPacket)

	// Encapsulate using ESP
	ipsecPacket, err := EncapsulateESP(sa, originalPacket)
	if err != nil {
		log.Fatalf("Failed to encapsulate: %v", err)
	}

	// Display the ESP packet
	dumpPacket(ipsecPacket)
	fmt.Println()

	// Decapsulate the ESP packet
	decapsulated, err := DecapsulateESP(sa, ipsecPacket)
	if err != nil {
		log.Fatalf("Failed to decapsulate: %v", err)
	}

	// Verify the decapsulated packet
	fmt.Printf("Decapsulated packet (%d bytes): %s\n", len(decapsulated), decapsulated)

	// Verify that decapsulation produces the original packet
	if string(decapsulated) == string(originalPacket) {
		fmt.Println("\nVerification successful: Decapsulated packet matches original packet")
	} else {
		fmt.Println("\nVerification failed: Decapsulated packet does not match original packet")
	}
}
