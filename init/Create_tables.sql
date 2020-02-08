create table users (
	user_id integer constraint users_PK primary key,
	token varchar(100) not null
);