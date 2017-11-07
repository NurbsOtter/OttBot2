CREATE TABLE `chatUser` (
	`id` int PRIMARY KEY AUTO_INCREMENT NOT NULL,
	`userName` text CHARACTER SET utf8mb4,
	`tgID` int NOT NULL,
	`pingAllowed` tinyint(1) NOT NULL DEFAULT '1',
	`activeUser` tinyint(1) NOT NULL DEFAULT '1',
	`MarkovAskAllowed` tinyint(1) NOT NULL DEFAULT '1'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `aliases` (
	`id` int PRIMARY KEY AUTO_INCREMENT NOT NULL,
	`name` text CHARACTER SET utf8mb4 NOT NULL,
	`userID` int NOT NULL,
	`changeDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(`userID`) REFERENCES chatUser(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE `warning` (
	`id` int PRIMARY KEY AUTO_INCREMENT NOT NULL,
	`userID` int NOT NULL,
	`warningText` text CHARACTER SET utf8mb4 NOT NULL,
	`warnDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(`userID`) REFERENCES chatUser(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

