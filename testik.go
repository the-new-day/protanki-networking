package main

import (
	"fmt"

	"github.com/the-new-day/probogo/pkg/modules/protection"
)

func main() {
	// Твои ключи
	keys := []byte{0x12, 0x34, 0x56, 0x78} // например

	// 1. Создаем защиту с ТЕМИ ЖЕ ключами
	protector := protection.NewXorProtection(false) // или true, смотри как надо
	protector.Activate(keys)

	// 2. Берем первые 10 байт зашифрованных данных
	first4Encrypted := []byte{}

	// 3. Расшифровываем ТОЛЬКО их
	first4Decrypted := protector.Decrypt(first4Encrypted)

	fmt.Printf("Первые 10 расшифрованных: %x\n", first4Decrypted)
}
