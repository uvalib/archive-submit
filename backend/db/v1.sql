
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Create table to track schema version
--
DROP TABLE IF EXISTS versions;
CREATE TABLE versions (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   version varchar(10) NOT NULL UNIQUE KEY,
   created_at datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
insert into versions(version, created_at) values ("v1", NOW());

--
-- Create table for users
--
DROP TABLE IF EXISTS genres;
CREATE TABLE genres (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   name varchar(255) DEFAULT NULL,
   UNIQUE KEY (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Add seed genres
insert into genres(name)
   values ("Correspondence"), ("Reports"), ("Publications"), ("Speeches/Lectures"), ("Research Data"), 
          ("Manuscripts"), ("Course Material"), ("Biographical Information"), ("Diaries/Journals");

--
-- Create table for users
--
DROP TABLE IF EXISTS users;
CREATE TABLE users (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   last_name varchar(255) DEFAULT NULL,
   first_name varchar(255) DEFAULT NULL,
   email varchar(255) NOT NULL,
   title varchar(50),
   university_affiliation varchar(50),
   phone varchar(20),
   created_at datetime NOT NULL,
   updated_at datetime NOT NULL,
   UNIQUE KEY (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Add starter user
insert into users(last_name,first_name,email,university_affiliation,created_at,updated_at)
   values ("Foster", "Lou", "lf6f@virginia.edu", "UVA", NOW(), NOW() );


--
-- Create table for generl accession info
--
DROP TABLE IF EXISTS accessions;
CREATE TABLE accessions (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   user_id int(11) NOT NULL,
   description text DEFAULT NULL,
   activities text DEFAULT NULL,
   creator text DEFAULT NULL,
   FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
