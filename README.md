# Person Service

1. **Клонировать репозиторий**  
   ```bash
   git clone https://github.com/Ravwvil/person-service
   cd person-service
   ```


2. **Запустить сервис**  
   ```bash
   go run cmd/main.go
   ```

3. **Открыть Swagger UI**  
   http://localhost:8080/swagger/index.html

---

## 🔌 API

| Метод       | Путь                  | Описание                             |
|-------------|-----------------------|--------------------------------------|
| **POST**    | `/people`             | Создать нового человека              |
| **GET**     | `/people`             | Список с фильтрами и пагинацией      |
| **PUT**     | `/people/{id}`        | Обновить данные по ID                |
| **DELETE**  | `/people/{id}`        | Удалить по ID                        |

---