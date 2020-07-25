CREATE TABLE users (
	user_id SERIAL PRIMARY KEY
,	username VARCHAR (150) UNIQUE NOT NULL
,	created_at TIMESTAMP WITH TIME ZONE DEFAULT (now() at TIME ZONE 'utc') 
);

CREATE TABLE chats (
	chat_id SERIAL PRIMARY KEY
,	name VARCHAR(150) UNIQUE NOT NULL
,	created_at TIMESTAMP WITH TIME ZONE DEFAULT (now() at TIME ZONE 'utc') 
);

CREATE TABLE chat_user (
	chat_id INT REFERENCES chats (chat_id) ON UPDATE CASCADE,
	user_id INT REFERENCES users (user_id) ON UPDATE CASCADE
,	CONSTRAINT chat_user_pkey PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE messages (
	message_id SERIAL PRIMARY KEY
,	chat_id INT REFERENCES chats (chat_id) 
,	author_id INT REFERENCES users(user_id)
,	text TEXT NOT NULL
,	created_at TIMESTAMP WITH TIME ZONE DEFAULT (now() at TIME ZONE 'utc') 
);


