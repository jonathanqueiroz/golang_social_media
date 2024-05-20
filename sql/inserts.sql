INSERT INTO users (name, email, username, avatar_url, bio, birthdate, password) VALUES ('Alice', 'alice@mail.com', 'alice', 'https://i.pravatar.cc/150?img=1', 'I am a software engineer', '1990-01-01', '$2a$10$7CuLzispLK.g/YCdW4uRrOu3PnS0..Z8VkcnB0.xVEiFLSySmNwPW');

INSERT INTO users (name, email, username, avatar_url, bio, birthdate, password) VALUES ('Bob', 'bob@mail.com', 'bob', 'https://i.pravatar.cc/150?img=2', 'I am a software engineer', '1990-01-01', '$2a$10$7CuLzispLK.g/YCdW4uRrOu3PnS0..Z8VkcnB0.xVEiFLSySmNwPW');

INSERT INTO users (name, email, username, avatar_url, bio, birthdate, password) VALUES ('Charlie', 'charlie@mail.com', 'charlie', 'https://i.pravatar.cc/150?img=3', 'I am a software engineer', '1990-01-01', '$2a$10$7CuLzispLK.g/YCdW4uRrOu3PnS0..Z8VkcnB0.xVEiFLSySmNwPW');

INSERT INTO users (name, email, username, avatar_url, bio, birthdate, password) VALUES ('David', 'david@mail.com', 'david', 'https://i.pravatar.cc/150?img=4', 'I am a software engineer', '1990-01-01', '$2a$10$7CuLzispLK.g/YCdW4uRrOu3PnS0..Z8VkcnB0.xVEiFLSySmNwPW');

INSERT INTO users (name, email, username, avatar_url, bio, birthdate, password) VALUES ('Eve', 'eve@mail.com', 'eve', 'https://i.pravatar.cc/150?img=5', 'I am a software engineer', '1990-01-01', '$2a$10$7CuLzispLK.g/YCdW4uRrOu3PnS0..Z8VkcnB0.xVEiFLSySmNwPW');

INSERT INTO followers (user_id, follower_id) VALUES (1, 2);
INSERT INTO followers (user_id, follower_id) VALUES (1, 3);
INSERT INTO followers (user_id, follower_id) VALUES (1, 4);
INSERT INTO followers (user_id, follower_id) VALUES (2, 3);
INSERT INTO followers (user_id, follower_id) VALUES (2, 4);
INSERT INTO followers (user_id, follower_id) VALUES (2, 5);
INSERT INTO followers (user_id, follower_id) VALUES (3, 4);
INSERT INTO followers (user_id, follower_id) VALUES (3, 5);
INSERT INTO followers (user_id, follower_id) VALUES (4, 5);


INSERT INTO posts (author_id, content) VALUES (1, 'Content of the post');
INSERT INTO posts (author_id, content) VALUES (2, 'Another post content');
INSERT INTO posts (author_id, content) VALUES (3, 'Hello, World!');
INSERT INTO posts (author_id, content, parent_id) VALUES (1, 'Reply to the first post', 1);
INSERT INTO posts (author_id, content, parent_id) VALUES (2, 'Reply to the second post', 2);
INSERT INTO posts (author_id, content, parent_id) VALUES (3, 'Reply to the third post', 3);
INSERT INTO posts (author_id, content, parent_id) VALUES (4, 'Replying Post Another Time', 3);

INSERT INTO likes (post_id, user_id) VALUES (1, 2);
INSERT INTO likes (post_id, user_id) VALUES (1, 3);
INSERT INTO likes (post_id, user_id) VALUES (1, 4);
INSERT INTO likes (post_id, user_id) VALUES (2, 3);
INSERT INTO likes (post_id, user_id) VALUES (2, 4);
INSERT INTO likes (post_id, user_id) VALUES (2, 5);
INSERT INTO likes (post_id, user_id) VALUES (3, 4);
INSERT INTO likes (post_id, user_id) VALUES (3, 5);
INSERT INTO likes (post_id, user_id) VALUES (4, 5);

