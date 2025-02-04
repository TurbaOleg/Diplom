CREATE TABLE IF NOT EXISTS rules (
	id INTEGER PRIMARY KEY NOT NULL,
	domain_pattern TEXT NOT NULL,
	is_secure INTEGER NOT NULL,
	is_http_only INTEGER NOT NULL,
	same_site INTEGER NOT NULL
);
insert into rules(domain_pattern, is_secure, is_http_only, same_site) values('%', 1, 1, 2);
