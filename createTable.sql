CREATE TABLE forum (
  postId         BIGINT NOT NULL PRIMARY KEY,
  posterId      VARCHAR(255) NOT NULL,
  postDate     VARCHAR(255) NOT NULL,
  commId      VARCHAR(255) NOT NULL,
  parentPostId      VARCHAR(255) NOT NULL,
  textContent      TEXT NOT NULL,
  mediaLinks      TEXT NOT NULL,
  eventId         VARCHAR(255) NOT NULL
);