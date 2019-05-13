DROP DATABASE `whohastheball` IF EXISTS; 
CREATE DATABASE `whohastheball`;

CREATE TABLE `whohastheball`.`user` (
  `user_id` INT NOT NULL AUTO_INCREMENT,
  `twitter_id` VARCHAR(45) NOT NULL,
  `twitter_screenname` VARCHAR(15) NULL,
  `created_at` DATETIME NOT NULL DEFAULT NOW(),
  `catches` INT NULL DEFAULT 1,
  PRIMARY KEY (`user_id`, `twitter_id_str`));

CREATE TABLE `whohastheball`.`possession` (
  `possession_id` INT NOT NULL AUTO_INCREMENT,
  `user_id` INT(11) NOT NULL,
  `received_at` DATETIME NOT NULL DEFAULT NOW(),
  `taken_at` DATETIME NULL,
  `duration` INT NULL,
  PRIMARY KEY (`possession_id`),
  INDEX `FK_USER_ID_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `FK_USER_ID`
    FOREIGN KEY (`user_id`)
    REFERENCES `whohastheball`.`user` (`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);