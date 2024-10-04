CREATE TABLE IF NOT EXISTS roles 
		(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE
		);
CREATE TABLE IF NOT EXISTS member_roles 
(
	member_id INTEGER,
	role_id INTEGER,
	PRIMARY KEY (member_id, role_id),
	FOREIGN KEY (member_id) REFERENCES members (id),
	FOREIGN KEY (role_id) REFERENCES roles (id)
);
CREATE TABLE IF NOT EXISTS members 
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	guild TEXT,
	username TEXT,
	nickname TEXT,
	avatarurl TEXT,
	displayavatarurl TEXT,
	bannerurl TEXT,
	displaybannerurl TEXT,
	usercolor TEXT,
	joinedat TEXT,
	userstatus TEXT,
	msgcount INTEGER
);