# Clean Menu (Webapp for fast attention)

Clean Menu borns with a simple idea: fast attention at a low level for a restaurant, with powerful tools for the owner. The main goal is to provide a simple and clean interface for every employee. The idea is to have a webapp that can be used in a tablet, with a simple interface that can be used by anyone in the restaurant.

There will be three roles in the system: waiter, chef and owner. The waiter will be able to see the orders and mark them as done. The chef will be able to see the orders and mark them as done. The owner will be able to see the orders, mark them as done, see the statistics and manage the menu.

# MVP

The ideal flow work for the MVP is:

1. The owner creates an account and a the initial menu with prices and categories.
2. The waiter can take orders and mark them as done. Also the waiter send the orders to the kitchen. The waiter can see the orders that are ready to be served.
3. The chef can see the orders that are ready to be served and mark them as done, and put the time that took to make the order.
4. The owner can see the orders that are ready to be served, mark them as done, see the statistics and manage the menu.

## Roles 

- **Owner**: Can see the orders, mark them as done, see the statistics and manage the menu.
- **Waiter**: Can see the orders and mark them as done. The waiter can send the orders to the kitchen.
- **Chef**: Can see the orders and mark them as done. The chef can put the time that took to make the order.

# Technical details

## Frontend

The frontend will be using only templates in Go. The idea is to have a simple interface that can be used in a tablet. The frontend will be using the following technologies:

- HTML
- CSS
- JavaScript
- Go templates

## Backend

The backend will be using only Go and only his native libraries. The idea is to have a simple backend that can be used in a Raspberry Pi!.

## Database

The database will be using SQLite for a local database.

## Connection

The connection will be using WebSockets for real-time updates. Because of the simplicity of the project, the WebSockets will be enough for the MVP.

## Security

The security will be using JWT for the authentication. The idea is to have a simple and secure way to authenticate the users.

# Roadmap

The roadmap for the project is:

- [ x ] Create the project structure
- [ ] Create the database structure
- [ x ] Create the owner interface
- [ x ] Create the waiter interface
- [ x ] Create the chef interface
- [ ] Create the statistics interface
- [ ] Create the menu management interface
- [ ] Create the WebSocket Connection
- [ ] Create the JWT authentication
- [ ] Create the Dockerfile
- [ ] Create the Docker Compose
- [ ] Create the tests
- [ ] Create the documentation

# Database connection

```bash
# Create the container


# Create the database
docker exec -it pg-docker-clean-menu psql -U postgres -c "CREATE DATABASE clean_menu_db;"
```