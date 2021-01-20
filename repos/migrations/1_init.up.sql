create table msig.contracts
(
	ctr_id serial not null
		constraint contracts_pk
			primary key,
	ctr_address varchar(36) not null
);

create unique index contracts_address_uindex
	on msig.contracts (ctr_address);

create table msig.requests
(
	req_id serial not null
		constraint requests_pk
			primary key,
	ctr_id int not null
		constraint requests_contracts_id_fk
			references msig.contracts,
    req_hash varchar(32) not null,
	req_status varchar default 'wait' not null,
	req_counter int not null,
	req_data text not null
);

create table msig.signatures
(
	sig_id serial not null
		constraint signatures_pk
			primary key,
	req_id int not null
		constraint signatures_requests_id_fk
			references msig.requests,
	sig_index int not null,
	sig_data varchar not null
);

create unique index signatures_sign_uindex
	on msig.signatures (sig_data);


