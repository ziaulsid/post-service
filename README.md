# Post Service

### User Stories
- Allow user to make a post with caption and an image (only 1 image per post)
- Allow user to comment on a post
- Allow user to view all the posts including recent 2 comments

### Setup of the Project
- Run `docker-compose up --build` command to setup the database
- docker-compose.yml file will read the values from .env file to connect to database
- Execute `go build .` and `go run main.go` commands to build and run the application. 
- Your server should be ready to serve the endpoints

- Setup a .env file
```
DB_USER=root
DB_PASSWORD=
DB_NAME=post_db
DB_HOST=localhost:3306
DB_DRIVER=mysql

```

- Alternatively, if using Goland or any other IDE, follow this setup
<img width="1066" alt="image" src="https://github.com/user-attachments/assets/9135577d-3b32-48ae-afd1-5cb1efa11a26">
<img width="501" alt="image" src="https://github.com/user-attachments/assets/7ab98cf0-f975-49f7-a6ee-ee32db5540d8">



### API Endpoints
1.  Create a post
```
curl --location 'http://localhost:8080/posts' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": 1001,
    "caption": "chilling in Norway"
}'
```

2. Add image to a post
```
curl --location 'http://localhost:8080/posts/1/images' \
--form 'image=@/Users/siddiquizia/Downloads/oskar-smethurst-B1GtwanCbiw-unsplash.jpg'
```

3. Add comment to a post
```
curl --location --request POST 'http://localhost:8080/posts/1/comments' \
--header 'Content-Type: application/json' \
--data-raw '{
    "content": "Great post!",
    "user_id": 456
}'

```

- Get post of a user
```
curl --location --request GET 'http://localhost:8080/users/1001/posts' 
```

### Key Design Decisions
- For simplicity, images are stored directly in the relational database (MySQL). However, a more optimal approach would be to use blob storage (e.g., S3, GCS) for storing images and save the image URL in the database.
- There is a dedicated API for uploading images. In a production environment, the backend server would typically generate a signed URL for direct uploads to blob storage such as S3 or Azure Blob Storage.
- Image resizing (optimization) is performed asynchronously to avoid blocking the main request goroutines. This is managed by a WorkerPool with multiple workers. In a production environment, using message queues and a separate set of workers for this operation would be a more scalable approach.


**Note:** The user story to delete comment from post is not implemented. Also test cases are not comprehensive due to work.
