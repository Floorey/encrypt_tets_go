package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func caesarEncrypt(text string, shift int) string {
	encrypted := ""
	for _, char := range text {
		if char >= 'a' && char <= 'z' {
			encrypted += string(((int(char) - 'a' + shift) % 26) + 'a')
		} else if char >= 'A' && char <= 'Z' {
			encrypted += string(((int(char) - 'A' + shift) % 26) + 'A')
		} else {
			encrypted += string(char)
		}
	}
	return encrypted
}
func rot13Encrypted(text string) string {
	encrypted := ""
	for _, char := range text {
		if char >= 'a' && char <= 'z' {
			encrypted += string(((int(char) - 'a' + 13) % 26) + 'a')
		} else if char >= 'A' && char <= 'Z' {
			encrypted += string(((int(char) - 'A' + 13) % 26) + 'A')
		} else {
			encrypted += string(char)
		}
	}
	return encrypted
}
func substitutionEncryption(text string, key string) string {
	encrypted := ""
	for _, char := range text {
		if char >= 'a' && char <= 'z' {
			encrypted += string(key[int(char)-'a'])
		} else if char >= 'A' && char <= 'Z' {
			encrypted += string(key[int(char)-'A'])
		} else {
			encrypted += string(char)
		}
	}
	return encrypted
}

func asciiEncrypted(text string) string {
	encrypted := ""
	for _, char := range text {
		encrypted += strconv.Itoa(int(char)) + " "
	}
	return encrypted
}

func writeToFile(filename string, content string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	writer.Flush()

	fmt.Println("Text successfully written to file:", filename)
}
func sendToAPI(text string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"text": text,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://localhost:8080", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var responseBody map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return "", err
	}
	return responseBody["encrypted_text"], nil

}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter text to encrypt: ")
	text, _ := reader.ReadString('\n')

	fmt.Print("Choose encryption method: ")
	fmt.Print("1. Caesar ")
	fmt.Print("2. rot13 ")
	fmt.Print("3. substitution ")
	fmt.Print("4. ASCII ")
	fmt.Print("Enter your choice:")
	choice, _ := reader.ReadString('\n')
	choice = strings.Trim(choice, "\r\n")

	switch choice {
	case "1":
		fmt.Print("Enter the shift amount for Caesar encryption")
		shiftStr, _ := reader.ReadString('\n')
		shiftStr = strings.TrimSpace(shiftStr)
		shift, _ := strconv.Atoi(shiftStr)

		caesarEncryptText := caesarEncrypt(text, shift)
		writeToFile("caesar_encrypted_text.txt", caesarEncryptText)
		fmt.Print("Encrypted text (Caesar):", caesarEncryptText)
	case "2":
		rot13EncryptedText := rot13Encrypted(text)
		writeToFile("rot13_encrypted_text.txt", rot13EncryptedText)
		fmt.Print("Encrypted text (Rot13):", rot13EncryptedText)
	case "3":
		fmt.Print("Enter the key for Substitution encryption:")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		substitutionEncryptedText := substitutionEncryption(text, key)
		writeToFile("substitution_encrypted_text.txt", substitutionEncryptedText)
		fmt.Print("Encrypted text (Substitution):", substitutionEncryptedText)
	case "4":
		asciiEncryptedText := asciiEncrypted(text)
		writeToFile("ascii_encrypted_text.txt", asciiEncryptedText)
		fmt.Print("Encrypted text (ASCII):", asciiEncryptedText)
	case "5":
		encryptText, err := sendToAPI(text)
		if err != nil {
			fmt.Println("Error sending to API:", err)
			return
		}
		writeToFile("api_encrypted_text.txt", encryptText)
		fmt.Print("Encrypted text (API):", encryptText)
	default:
		fmt.Println("Invalid choice:")
	}
}
