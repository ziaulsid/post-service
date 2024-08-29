# Post Service

### User Stories
- Allow user to make a post with caption and an image (only 1 image per post)
- Allow user to comment on a post
- Allow user to view all the posts including recent 2 comments

### Setup of the Project
- Requirements: `docker-compose`.
- Run `docker-compose up --build` command to setup the database and run the app within the container itself
  <img width="966" alt="image" src="https://github.com/user-attachments/assets/1eba5414-d388-4bc9-9817-e22daf55b5cf">



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


**Note:** The user story to delete comment from post is not implemented. Also test cases are not comprehensive as i am travelling outside Singapore due to work.
