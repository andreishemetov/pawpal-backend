// @title PawPal API
// @version 1.0
// @description API for managing pets in PawPal.
// @termsOfService https://example.com/terms/

// @contact.name API Support
// @contact.url https://example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/andreishemetov/pawpal/internal/app"
	"github.com/andreishemetov/pawpal/internal/config"
)

/*
HTTP layer  →  Handler
Handler     →  Service
Service     →  Data model
*/

func main() {
	cfg, err := config.LoadForAPI()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.RunAPI(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}

/*

swag init -g cmd/api/main.go -o docs

curl -X POST http://localhost:8080/pets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"id":1,"name":"Charlie","age":3,"visits":0}'

curl -X POST http://localhost:8080/pets \
  -H "Content-Type: application/json" \
  -d '{"id":1,"name":"Charlie","age":3,"visits":0}'


curl http://localhost:8080/pets

curl http://localhost:8080/pets \
  -H "Authorization: Bearer YOUR_TOKEN"

curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

  curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'

  docker compose exec db psql -U pawpal -d pawpal_dev -c "
UPDATE users
SET role = 'admin'
WHERE email = 'admin@example.com';


curl -X POST http://localhost:8080/reminders \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "pet_id": 1,
    "message": "Give medicine",
    "remind_at": "2026-04-02T12:00:00Z"
  }'



docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0001_init.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0002_add_indexes.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0003_users.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0004_add_user_id_to_pets.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0005_refresh_tokens.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0006_add_user_role.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0007_reminders.up.sql
docker compose exec -T db psql -U pawpal -d pawpal_dev < migrations/0008_reminder_delivery_fields.up.sql

*/
