create table users (
	user_id integer constraint users_PK primary key,
	token varchar(100) not null
);

create table user_sources (
	user_id integer not null,
	source varchar(100) not null,
	UNIQUE (user_id, source)
);