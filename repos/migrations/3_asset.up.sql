create table msig.assets
(
	ast_id serial not null
		constraint asset_pk
			primary key,
    ast_name varchar not null,
    ast_contract_type varchar not null,
    ast_address varchar(36) not null,
    ast_dexter_address varchar(36) not null,
    ast_scale int not null,
    ast_ticker varchar not null
);