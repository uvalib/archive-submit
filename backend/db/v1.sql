
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
-- Create table for genres
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
          ("Manuscripts"), ("Course Material"), ("Biographical Information"), ("Diaries/Journals"),
          ("Social Media Content"), ("Websites");

--
-- Create table for record types
--
DROP TABLE IF EXISTS object_types;
DROP TABLE IF EXISTS record_types;
CREATE TABLE record_types (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   name varchar(30) NOT NULL UNIQUE KEY,
   description varchar(255) DEFAULT NULL,
   digital boolean DEFAULT TRUE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Add seed types
insert into record_types(name,description,digital)
   values ("Text", "(e.g., reports, contracts, email)", 1), ("Images", "(e.g., TIFFs, PDFs)", 1), 
   ("Video", "(e.g., How-to-videos, event recordings)", 1), ("Audio/Sound Recordings", "(e.g., interviews)", 1),
   ("Software/Multimedia", "(e.g., SVG, Python)", 1), ("Databases/Data", "(e.g., relational databases, research data)", 1),
   ("Websites", "(e.g., archived, content-based)", 1), ("Text/Documents", "", 0),
   ("Photographs/Still Images", "", 0), ("Audio/Video Media", "", 0);

--
-- Create table for transfer_methods
--
DROP TABLE IF EXISTS transfer_methods;
CREATE TABLE transfer_methods (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   name varchar(255) DEFAULT NULL,
   UNIQUE KEY (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Add seed transfer_methods
insert into transfer_methods(name)
   values ("UVA campus mail"), ("Personal delivery"), ("Pickup by UVA Facilities");

--
-- Create table for media_carriers
--
DROP TABLE IF EXISTS media_carriers;
CREATE TABLE media_carriers (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   name varchar(255) DEFAULT NULL,
   UNIQUE KEY (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- Add seed media_carriers
insert into media_carriers(name)
   values ('5 1/4" Disk'), ('3 1/2" Disk'), ("CD-ROM"), ("DVD-R/W"), ("USB Flash Drive"), ("Hard Drive");

--
-- Create table for users
--
DROP TABLE IF EXISTS users;
CREATE TABLE users (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   last_name varchar(255) DEFAULT NULL,
   first_name varchar(255) DEFAULT NULL,
   email varchar(255) NOT NULL,
   title varchar(50) DEFAULT NULL,
   university_affiliation varchar(50) DEFAULT NULL,
   phone varchar(20) DEFAULT NULL,
   verified boolean default FALSE,
   verify_token varchar(25), 
   admin boolean default FALSE,
   created_at datetime NOT NULL,
   updated_at datetime NOT NULL,
   INDEX (verify_token),
   UNIQUE KEY (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

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

--
-- Create table for inventory items
--
DROP TABLE IF EXISTS inventory_items;
CREATE TABLE inventory_items (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   accession_id int(11) NOT NULL,
   box_number varchar(100),
   record_group_number varchar(100),
   box_title varchar(255),
   description text,
   dates varchar(255),
   FOREIGN KEY (accession_id) REFERENCES accession(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
