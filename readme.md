Certainly! Below is a README.md file for the provided code:

```markdown
# Simple Book Inventory API

This is a simple RESTful API built with Go (Golang) and Gin framework to manage a book inventory. The API allows you to perform CRUD (Create, Read, Update, Delete) operations on a collection of books stored in a MongoDB database.

## Features

- Create a new book
- Retrieve all books or a specific book by its ID
- Update the quantity of a book when checking out or returning
- Support for single and batch book creation

## Dependencies

- Go 1.16 or higher
- MongoDB

## Installation

1. Clone the repository:

```bash
git clone <repository_url>
```

2. Install dependencies:

```bash
go mod tidy
```

3. Set up MongoDB:

    - Install MongoDB on your system or use a cloud-based MongoDB service such as MongoDB Atlas.
    - Create a MongoDB database and collection for storing books.
    - Replace the MongoDB connection string in the code (`initMongoDB` function) with your own connection string.

4. Build and run the application:

```bash
go run main.go
```

## API Endpoints

### 1. Retrieve all books

```http
GET /books
```

### 2. Retrieve a book by ID

```http
GET /books/:id
```

### 3. Create a new book

```http
POST /book
```

Request Body (JSON):

```json
{
    "id": "unique_id",
    "title": "Book Title",
    "author": "Author Name",
    "quantity": 5,
    "max_present": 10
}
```

### 4. Create multiple books

```http
POST /books
```

Request Body (JSON):

```json
[
    {
        "id": "unique_id_1",
        "title": "Book 1 Title",
        "author": "Author Name",
        "quantity": 5,
        "max_present": 10
    },
    {
        "id": "unique_id_2",
        "title": "Book 2 Title",
        "author": "Author Name",
        "quantity": 8,
        "max_present": 15
    },
    ...
]
```

### 5. Update book quantity (Checkout)

```http
PATCH /checkout/:id
```

### 6. Update book quantity (Return)

```http
PATCH /return/:id
```

## Contributors

- [Adhayan Sharma](https://github.com/adhayan002)
```

Feel free to customize and expand upon this README.md according to your project's specific requirements and features.