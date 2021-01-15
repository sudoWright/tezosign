create table msig.contracts
(
	address varchar(36) not null,
	id serial not null
		constraint contracts_pk
			primary key
);

create unique index contracts_address_uindex
	on msig.contracts (address);


create table msig.requests
(
	data text not null,
	status varchar default 'wait',
	ctr_id int not null
		constraint requests_contracts_id_fk
			references msig.contracts,
	id serial not null
		constraint requests_pk
			primary key
);

create table msig.signatures
(
	id serial not null
		constraint signatures_pk
			primary key,
	req_id int not null
		constraint signatures_requests_id_fk
			references msig.requests,
	sign varchar not null
);

create unique index signatures_sign_uindex
	on msig.signatures (sign);


