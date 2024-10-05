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

CREATE TABLE IF NOT EXISTS trigger_words (
    trigger_id integer primary key AUTOINCREMENT,
    bot_id integer foreign key REFERENCES bots (botid),
    phrase varchar(255),
    author char(20),
    data_datetime DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS typical_lines (
    line_id integer primary key AUTOINCREMENT,
    bot_id integer foreign key REFERENCES bots (botid) ,
    phrase varchar(255),
    author char(20),
	to_member varchar(255),
    data_datetime DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sad_lines (
    line_id integer primary key AUTOINCREMENT,
    bot_id integer foreign key references bots (botid),
    phrase varchar(255),
    author char(20),
	to_id varchar(255),
    data_datetime DATETIME DEFAULT CURRENT_TIMESTAMP
);