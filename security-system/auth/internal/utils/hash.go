package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) { // Функция для хэширования пароля с использованием bcrypt. Она принимает пароль в виде строки и возвращает хэшированный пароль и ошибку, если она возникла.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Генерируем хэш пароля с использованием bcrypt и стандартной стоимости.
	return string(bytes), err                                                       // Возвращаем хэшированный пароль в виде строки и ошибку, если она возникла.
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) // Сравниваем хэшированный пароль с введенным паролем. Если они совпадают, возвращаем true, иначе false.
	return err == nil
}
