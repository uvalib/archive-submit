
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
   ("Websites", "(e.g., archived, content-based)", 1), 
   ("Social Media Content", "(e.g., downloaded data from Twitter, Facebook, or Instagram)",1),
   ("Text/Documents", "", 0),
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
   api_token varchar(25) NOT NULL DEFAULT "",
   created_at datetime NOT NULL,
   updated_at datetime NOT NULL,
   INDEX (verify_token),
   INDEX (api_token),
   UNIQUE KEY (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO users (last_name, first_name, email, title, university_affiliation, phone, verified, admin, created_at)
VALUES  
	('Foster', 'Lou', 'lf6f@virginia.edu', 'Software Engineer', 'UVA Library', '(434) 982-2812', 1, 1, NOW()),
   ('Work', 'Lauren', 'lw2cd@virginia.edu', 'Digital Preservation Librarian', 'UVA Library', '(434) 924-1348', 1, 1, NOW()),
   ('Anderson', 'Bethany', 'bga3d@virginia.edu', 'University Archivist', 'UVA Library', '', 1, 1, NOW());

--
-- Create table for general accession info
--
DROP TABLE IF EXISTS accessions;
CREATE TABLE accessions (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   identifier varchar(25) not null,
   user_id int(11) NOT NULL,
   description text DEFAULT NULL,
   activities text DEFAULT NULL,
   creator text DEFAULT NULL,
   accession_type varchar(25), 
   created_at datetime default CURRENT_TIMESTAMP,
   unique index(identifier),
   FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Create table for accession genres
--
DROP TABLE IF EXISTS accession_genres;
CREATE TABLE accession_genres (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   accession_id int(11) NOT NULL,
   genre_id int(11) NOT NULL,
   FOREIGN KEY (accession_id) REFERENCES accessions(id) ON DELETE CASCADE,
   FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Create table for digital accessions
--
DROP TABLE IF EXISTS digital_accessions;
CREATE TABLE digital_accessions (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   accession_id int(11) NOT NULL,
   upload_size int(11),
   date_range varchar(255),
   description text,
   FOREIGN KEY (accession_id) REFERENCES accessions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Create table for accession record types
--
DROP TABLE IF EXISTS accession_record_types;
CREATE TABLE accession_record_types (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   accession_id int(11) NOT NULL,
   accession_type varchar(10) not null,
   record_type_id int(11) NOT NULL,
   INDEX(accession_type),
   FOREIGN KEY (accession_id) REFERENCES accessions(id) ON DELETE CASCADE,
   FOREIGN KEY (record_type_id) REFERENCES record_types(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Create table for physical accessions
--
DROP TABLE IF EXISTS physical_accessions;
CREATE TABLE physical_accessions (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   accession_id int(11) NOT NULL,
   date_range varchar(255),
   box_info varchar(255),
   transfer_method_id int(11) NOT NULL,
   has_digital boolean not null default false,
   has_software varchar(10) not null default "no",
   tech_description text,
   media_counts varchar(255),
   FOREIGN KEY (transfer_method_id) REFERENCES transfer_methods(id) ON DELETE CASCADE,
   FOREIGN KEY (accession_id) REFERENCES accessions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Create table for physicl accession media carriers
--
DROP TABLE IF EXISTS physical_media_carriers;
CREATE TABLE physical_media_carriers (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   physical_accession_id int(11) NOT NULL,
   media_carrier_id int(11) NOT NULL,
   FOREIGN KEY (physical_accession_id) REFERENCES physical_accessions(id) ON DELETE CASCADE,
   FOREIGN KEY (media_carrier_id) REFERENCES media_carriers(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Create table for digital files
--
DROP TABLE IF EXISTS digital_files;
CREATE TABLE digital_files (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   digital_accession_id int(11) NOT NULL,
   filename varchar(255) NOT NULL,
   FOREIGN KEY (digital_accession_id) REFERENCES digital_accessions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Create table for inventory items
--
DROP TABLE IF EXISTS inventory_items;
CREATE TABLE inventory_items (
   id int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
   physical_accession_id int(11) NOT NULL,
   box_number varchar(100),
   record_group_number varchar(100),
   box_title varchar(255),
   description text,
   dates varchar(255),
   FOREIGN KEY (physical_accession_id) REFERENCES physical_accessions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
