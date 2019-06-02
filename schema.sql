DROP SCHEMA IF EXISTS `seizetheball`;
CREATE SCHEMA `seizetheball`;

CREATE TABLE `seizetheball`.`user` (
  `user_id` INT NOT NULL AUTO_INCREMENT,
  `twitter_id` VARCHAR(45) UNIQUE NOT NULL,
  `screen_name` VARCHAR(15) NULL,
  `created_at` DATETIME NOT NULL DEFAULT NOW(),
  PRIMARY KEY (`user_id`, `twitter_id`));

CREATE TABLE `seizetheball`.`possession` (
  `possession_id` INT NOT NULL AUTO_INCREMENT,
  `tweet_id` VARCHAR(45) UNIQUE NOT NULL,
  `user_id` INT(11) NOT NULL,
  `start` DATETIME NOT NULL DEFAULT NOW(),
  `end` DATETIME NULL,
  `duration` INT NOT NULL DEFAULT 0,
  PRIMARY KEY (`possession_id`),
  INDEX `FK_USER_ID_idx` (`user_id` ASC),
  CONSTRAINT `FK_USER_ID`
    FOREIGN KEY (`user_id`)
    REFERENCES `seizetheball`.`user` (`user_id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE);