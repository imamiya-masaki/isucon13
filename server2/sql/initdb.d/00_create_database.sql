CREATE DATABASE IF NOT EXISTS `isupipe`;

DROP USER IF EXISTS `isucon`@`%`;
CREATE USER `isucon`@`192.168.%` IDENTIFIED BY 'isucon';
GRANT ALL PRIVILEGES ON `isupipe`.* TO `isucon`@`192.168.%`;

CREATE DATABASE IF NOT EXISTS `isudns`;

DROP USER IF EXISTS `isudns`@`%`;
CREATE USER `isudns`@`192.168.%` IDENTIFIED BY 'isudns';
GRANT ALL PRIVILEGES ON `isudns`.* TO `isudns`@`192.168.%`;

