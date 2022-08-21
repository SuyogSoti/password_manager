create TABLE Users (
		Email varchar(255) NOT NULL PRIMARY KEY,
		HashedPassword varchar(255) NOT NULL,
);

create TABLE Passwords (
		UserEmail varchar(255) NOT NULL,
		Site varchar(255) not NULL,
		SiteUserName varchar(255) NOT NULL,
		HashedPassword varchar(255) NOT NULL,
		PRIMARY KEY(UserEmail, Site, SiteUserName),
		FOREIGN KEY (UserEmail) REFERENCES Users(Email) ON DELETE CASCADE);
