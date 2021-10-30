CREATE DATABASE testdb;
go

use testdb;
go

-- DROP TABLE IF EXISTS parcels;
-- GO

-- CREATE TABLE parcels (
--  uid int not null Primary key IDENTITY(1,1),
--  owner_uid int,
--  owner_room_name nvarchar(30),
--  owner_ryosei_name nvarchar(50),
--  register_datetime datetime,
--  register_staff_uid int,
--  register_staff_room_name nvarchar(30),
--  register_staff_ryosei_name nvarchar(50),
--  placement int default 0,
--  fragile int default 0,
--  is_released bit DEFAULT 0,
--  release_datetime datetime,
--  release_staff_uid int,
--  release_staff_room_name nvarchar(30),
--  release_staff_ryosei_name nvarchar(50),
--  checked_count int DEFAULT 0,
--  is_lost bit DEFAULT 0,
--  lost_datetime datetime,
--  is_returned bit DEFAULT 0,
--  returned_datetime datetime,
--  is_operation_error bit DEFAULT 0,
--  operation_error_type int,
--  note nvarchar(100),
--  is_deleted bit DEFAULT 0,
-- );
-- GO

-- DROP TABLE IF EXISTS parcel_event;
-- GO

-- CREATE TABLE parcel_event(
--  uid int not null Primary key IDENTITY(1,1),
--  created_at datetime,
--  event_type int,
--  parcel_uid int,
--  ryosei_uid int,
--  room_name nvarchar(30),
--  ryosei_name nvarchar(50),
-- target_event_uid int,
--  note nvarchar(100),
--  is_finished bit DEFAULT 0,
--  is_deleted bit DEFAULT 0,
-- );
-- GO

DROP TABLE IF EXISTS ryosei;
GO

CREATE TABLE ryosei(
 uid int not null Primary key IDENTITY(1,1),
 room_name nvarchar(30),
 ryosei_name nvarchar(50),
 ryosei_name_kana nvarchar(50),
 ryosei_name_alphabet nvarchar(50),
 block_id int,
 slack_id nvarchar(15),
 status int DEFAULT 1,
 parcels_current_count int DEFAULT 0,
 parcels_total_count int DEFAULT 0,
 parcels_total_waittime nvarchar(30) DEFAULT '0:00:00',
 last_event_id int,
 last_event_datetime datetime,
); 
GO

BULK INSERT ryosei
FROM '/database/seeds/dbo.ryosei.csv'
WITH
(
   FIELDTERMINATOR = ',',
   ROWTERMINATOR = '\n'
);
GO