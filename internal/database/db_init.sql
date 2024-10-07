CREATE TABLE IF NOT EXISTS roles 
		(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE
		);
		
CREATE TABLE IF NOT EXISTS member_roles 
(
	memberid INTEGER,
	roleid INTEGER,
	PRIMARY KEY (memberid, roleid),
	FOREIGN KEY (memberid) REFERENCES members (id),
	FOREIGN KEY (roleid) REFERENCES roles (id)
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

-- DROP TABLE bots;
CREATE TABLE IF NOT EXISTS bots (
	botid integer primary key AUTOINCREMENT,
	botguild varchar(255),
    botname varchar(255),
	avatarurl varchar(255),
	bannerurl varchar(255),
    createdat varchar(255),
	author varchar(255),
    botstatus varchar(255),
    isSinger boolean

);

CREATE TABLE IF NOT EXISTS lines (
	lineid integer primary key AUTOINCREMENT,
	bid integer,
	phrase text,
	author varchar(255),
	toid varchar(255),
	ltype varchar(255),
	createdat DATETIME DEFAULT CURRENT_TIMESTAMP,
	foreign key (bid) references bots (botid)
);
