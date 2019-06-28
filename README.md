# Family tree

Project with Go and Neo4j.

## Keywords

---

[Go, Neo4j, Clean Architecture]

## Requirements

---

- Docker
- Docker-compose
- Go >= 1.11
- Neo4j

## Installation

---

`git clone https://github.com/larien/family-tree.git`

## Run application

---

### Option 1: docker-compose [unstable]

```bash
make compose
```

### Option 2: docker

```bash
make docker
```

## Running tests

---

```bash
cd backend/tests/
go test
```

## Usage

---

### Add People

- **URL**

  _localhost:8899/api/v1/person_

- **Method:**

  `POST`

- **URL Params**

  `none`

- **Data Params**

  ```json
  [
    {
      "name": "<name>",
      "parents": ["<parent>", "<parent>"],
      "children": ["<child>", "<child>"]
    }
  ]
  ```

- **Success Response:**

  - **Code:** 201 CREATED <br />
    **Content:** `{ "message": "People registered successfully!" }`

- **Error Response:**

  - **Description:** Invalid JSON
    **Code:** 401 Bad Request <br />
    **Content:** `{ "message": "Failed to parse json" }`

  - **Description:** Invalid JSON
    **Code:** 500 Internal Server Error <br />
    **Content:** `{ "message": "Failed to register people" }`

- **Sample Body:**

```json
[
  {
    "name": "Bruce",
    "parents": ["Mike", "Phoebe"]
  },
  {
    "name": "Dunny",
    "parents": ["Mike", "Phoebe"]
  }
]
```

### Get All People

- **URL**

  _localhost:8899/api/v1/person_

- **Method:**

  `GET`

- **URL Params**

  `none`

- **Data Params**

  `none`

- **Success Response:**

  - **Code:** 200 OK <br />
    **Content:**

    ```json
    [
      {
        "name": "Bruce",
        "parents": ["Mike", "Phoebe"],
        "children": null
      },
      {
        "name": "Mike",
        "parents": null,
        "children": ["Bruce", "Dunny"]
      },
      {
        "name": "Phoebe",
        "parents": null,
        "children": ["Bruce", "Dunny"]
      },
      {
        "name": "Dunny",
        "parents": ["Mike", "Phoebe"],
        "children": null
      }
    ]
    ```

  - **Code:** 204 No Content <br />
    **Content:** `{ "message": "No people were found" }`

- **Error Response:**

  - **Description:** Invalid JSON
    **Code:** 500 Internal Server Error <br />
    **Content:** `{ "message": "Failed to find all people" }`

- **Sample Body:**

  `none`
