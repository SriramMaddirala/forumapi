CREATE TABLE forum (
  postId         BIGINT NOT NULL PRIMARY KEY,
  posterId      VARCHAR(255) NOT NULL,
  postDate     VARCHAR(255) NOT NULL,
  commId      VARCHAR(255) NOT NULL,
  parentPostId      VARCHAR(255) NOT NULL,
  textContent      TEXT NOT NULL,
  mediaLinks      VARCHAR(255) NOT NULL,
  eventId         VARCHAR(255) NOT NULL
);
CREATE TABLE users (
  posterId      VARCHAR(255) NOT NULL PRIMARY KEY,
  joinDate     VARCHAR(255) NOT NULL,
  username    VARCHAR(255) NOT NULL,
  pword      VARCHAR(255) NOT NULL,
  email         VARCHAR(255) NOT NULL
);
