CREATE TABLE links
(
id serial,
url varchar UNIQUE,
subreddit varchar,
title text,
created_at timestamp without time zone
)
