/*
 * forum
 *
 * Тестовое задание для реализации проекта \"Форумы\" на курсе по базам данных в Технопарке Mail.ru (https://park.mail.ru).
 *
 * API version: 0.1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package models

// Информация о пользователе.
type User struct {
	Id int `json:"id,omitempty"`

	// Имя пользователя (уникальное поле). Данное поле допускает только латиницу, цифры и знак подчеркивания. Сравнение имени регистронезависимо.
	Nickname string `json:"nickname,omitempty"`

	// Полное имя пользователя.
	Fullname string `json:"fullname"`

	// Описание пользователя.
	About string `json:"about"`

	// Почтовый адрес пользователя (уникальное поле).
	Email string `json:"email"`
}
