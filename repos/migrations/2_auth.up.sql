create table auth_tokens
(
	atn_id serial not null
		constraint auth_requests_pk
			primary key,
	atn_pubkey varchar(55) not null,
    atn_data varchar(64) not null,
    atn_type varchar not null,
    atn_is_used boolean not null,
    atn_expires_at timestamp without time zone not null
);
