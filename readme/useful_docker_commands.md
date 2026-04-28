Useful Docker commands

rebuild & run docker api:

docker compose build api    
docker compose up -d api  

Build and run:

docker compose up -d --build

Stop:

docker compose down

Stop and remove DB volume:

docker compose down -v

See running containers:

docker compose ps

View logs:

docker compose logs -f api
docker compose logs -f worker
docker compose logs -f db


MY API:
https://dashboard.render.com/
https://pawpal-backend-a4i8.onrender.com/