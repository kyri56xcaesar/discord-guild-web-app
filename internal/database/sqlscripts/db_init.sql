CREATE TABLE IF NOT EXISTS roles 
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	userid INTEGER,
	rolename TEXT,
	rolecolor TEXT,
	foreign key (userid) references members (id) 
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

CREATE TABLE IF NOT EXISTS messages (
	messageid	integer primary key AUTOINCREMENT,
	userid		integer,
	content		text,
	channel		text,
	createdat	text,
	foreign key (userid) references members (id)
);
