-- MySQL dump 10.13  Distrib 5.6.17, for Win64 (x86_64)
--
-- Host: 192.168.1.195    Database: brdb02
-- ------------------------------------------------------
-- Server version	5.6.21-70.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `rolelevelconfig`
--

DROP TABLE IF EXISTS `rolelevelconfig`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rolelevelconfig` (
  `rolelevelid` bigint(20) NOT NULL,
  `propbonus` int(11) DEFAULT NULL,
  `hp` int(11) DEFAULT NULL,
  `helpbonus` int(11) DEFAULT NULL,
  `goldbonus` int(11) DEFAULT NULL,
  `scorebonus` int(11) DEFAULT NULL,
  `skillbonus` int(11) DEFAULT NULL,
  `skillcasttime` int(11) DEFAULT NULL,
  `skillcoldtime` int(11) DEFAULT NULL,
  `price` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`rolelevelid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rolelevelconfig`
--

LOCK TABLES `rolelevelconfig` WRITE;
/*!40000 ALTER TABLE `rolelevelconfig` DISABLE KEYS */;
INSERT INTO `rolelevelconfig` VALUES (150000000,150,5,0,50,50,50,0,0,'0'),(150010001,2,1,0,1,1,0,0,0,'0'),(150020001,2,1,0,1,1,0,0,0,'3:400;'),(150020002,4,1,0,2,2,0,0,0,'0:200;'),(150020003,6,1,0,3,3,0,0,0,'0:400;'),(150020004,8,1,0,4,4,0,0,0,'0:600;'),(150020005,10,1,0,5,5,0,0,0,'0:800;'),(150020006,12,1,0,6,6,0,0,0,'0:1000;'),(150020007,14,1,0,7,7,0,0,0,'0:1200;'),(150020008,16,1,0,8,8,0,0,0,'0:1400;'),(150020009,18,1,0,9,9,0,0,0,'0:1600;'),(150020010,20,2,0,10,10,0,0,0,'0:1800;'),(150020011,22,2,0,11,11,0,0,0,'0:2000;'),(150020012,24,2,0,12,12,0,0,0,'0:2200;'),(150020013,26,2,0,13,13,0,0,0,'0:2400;'),(150020014,28,2,0,14,14,0,0,0,'0:2600;'),(150020015,30,2,0,15,15,0,0,0,'0:2800;'),(150020016,32,2,0,16,16,0,0,0,'0:3000;'),(150020017,34,2,0,17,17,0,0,0,'0:3200;'),(150020018,36,2,0,18,18,0,0,0,'0:3400;'),(150020019,38,2,0,19,19,0,0,0,'0:3600;'),(150020020,40,2,0,20,20,0,0,0,'0:3800;'),(150020021,42,2,0,21,21,0,0,0,'0:4000;'),(150020022,44,2,0,22,22,0,0,0,'0:4200;'),(150020023,46,2,0,23,23,0,0,0,'0:4400;'),(150020024,48,2,0,24,24,0,0,0,'0:4600;'),(150020025,50,2,0,25,25,0,0,0,'0:4800;'),(150020026,52,2,0,26,26,0,0,0,'0:5000;'),(150020027,54,2,0,27,27,0,0,0,'0:5200;'),(150020028,56,2,0,28,28,0,0,0,'0:5400;'),(150020029,58,2,0,29,29,0,0,0,'0:5600;'),(150020030,60,2,0,30,30,0,0,0,'0:5800;'),(150020031,62,2,0,31,31,0,0,0,'0:6000;'),(150020032,64,2,0,32,32,0,0,0,'0:6200;'),(150020033,66,2,0,33,33,0,0,0,'0:6400;'),(150020034,68,2,0,34,34,0,0,0,'0:6600;'),(150020035,70,2,0,35,35,0,0,0,'0:6800;'),(150020036,72,2,0,36,36,0,0,0,'0:7000;'),(150020037,74,2,0,37,37,0,0,0,'0:7200;'),(150020038,76,2,0,38,38,0,0,0,'0:7400;'),(150020039,78,2,0,39,39,0,0,0,'0:7600;'),(150020040,80,3,0,40,40,0,0,0,'0:7800;'),(150030001,2,2,25,1,1,0,0,100,'3:1000;'),(150030002,4,2,25,2,2,0,0,100,'0:200;'),(150030003,6,2,25,3,3,0,0,100,'0:400;'),(150030004,8,2,25,4,4,0,0,100,'0:600;'),(150030005,10,2,25,5,5,0,0,100,'0:800;'),(150030006,12,2,25,6,6,0,0,100,'0:1000;'),(150030007,14,2,25,7,7,0,0,100,'0:1200;'),(150030008,16,2,25,8,8,0,0,100,'0:1400;'),(150030009,18,2,25,9,9,0,0,100,'0:1600;'),(150030010,20,2,25,10,10,0,0,100,'0:1800;'),(150030011,22,2,25,11,11,0,0,90,'0:2000;'),(150030012,24,2,25,12,12,0,0,90,'0:2200;'),(150030013,26,2,25,13,13,0,0,90,'0:2400;'),(150030014,28,2,25,14,14,0,0,90,'0:2600;'),(150030015,30,2,25,15,15,0,0,90,'0:2800;'),(150030016,32,2,25,16,16,0,0,90,'0:3000;'),(150030017,34,2,25,17,17,0,0,90,'0:3200;'),(150030018,36,2,25,18,18,0,0,90,'0:3400;'),(150030019,38,2,25,19,19,0,0,90,'0:3600;'),(150030020,40,3,25,20,20,0,0,90,'0:3800;'),(150030021,42,3,25,21,21,0,0,90,'0:4000;'),(150030022,44,3,25,22,22,0,0,90,'0:4200;'),(150030023,46,3,25,23,23,0,0,90,'0:4400;'),(150030024,48,3,25,24,24,0,0,90,'0:4600;'),(150030025,50,3,50,25,25,0,0,90,'0:4800;'),(150030026,52,3,50,26,26,0,0,90,'0:5000;'),(150030027,54,3,50,27,27,0,0,90,'0:5200;'),(150030028,56,3,50,28,28,0,0,90,'0:5400;'),(150030029,58,3,50,29,29,0,0,90,'0:5600;'),(150030030,60,3,50,30,30,25,0,80,'0:5800;'),(150030031,62,3,50,31,31,25,0,80,'0:6000;'),(150030032,64,3,50,32,32,25,0,80,'0:6200;'),(150030033,66,3,50,33,33,25,0,80,'0:6400;'),(150030034,68,3,50,34,34,25,0,80,'0:6600;'),(150030035,70,3,50,35,35,25,0,80,'0:6800;'),(150030036,72,3,50,36,36,25,0,80,'0:7000;'),(150030037,74,3,50,37,37,25,0,80,'0:7200;'),(150030038,76,3,50,38,38,25,0,80,'0:7400;'),(150030039,78,3,50,39,39,25,0,80,'0:7600;'),(150030040,80,3,75,40,40,25,0,80,'0:7800;'),(150030041,82,3,75,41,41,25,0,80,'0:8000;'),(150030042,84,3,75,42,42,25,0,80,'0:8200;'),(150030043,86,3,75,43,43,25,0,80,'0:8400;'),(150030044,88,3,75,44,44,25,0,80,'0:8600;'),(150030045,90,3,75,45,45,25,0,80,'0:8800;'),(150030046,92,3,75,46,46,25,0,80,'3:120;'),(150030047,94,3,75,47,47,25,0,80,'3:120;'),(150030048,96,3,75,48,48,25,0,80,'3:120;'),(150030049,98,3,75,49,49,25,0,80,'3:120;'),(150030050,100,3,100,50,50,50,0,70,'3:120;'),(150030051,102,4,100,51,51,50,0,70,'3:120;'),(150030052,104,4,100,52,52,50,0,70,'3:120;'),(150030053,106,4,100,53,53,50,0,70,'3:120;'),(150030054,108,4,100,54,54,50,0,70,'3:120;'),(150030055,110,4,100,55,55,50,0,70,'3:120;'),(150030056,112,4,100,56,56,50,0,70,'3:120;'),(150030057,114,4,100,57,57,50,0,70,'3:120;'),(150030058,116,4,100,58,58,50,0,70,'3:120;'),(150030059,118,4,100,59,59,50,0,70,'3:120;'),(150030060,120,4,100,60,60,50,0,70,'3:120;'),(150040001,2,2,100,1,1,0,0,100,'3:1000;'),(150040002,4,2,100,2,2,0,0,100,'0:200;'),(150040003,6,2,100,3,3,0,0,100,'0:400;'),(150040004,8,2,100,4,4,0,0,100,'0:600;'),(150040005,10,2,100,5,5,0,0,100,'0:800;'),(150040006,12,2,100,6,6,0,0,100,'0:1000;'),(150040007,14,2,100,7,7,0,0,100,'0:1200;'),(150040008,16,2,100,8,8,0,0,100,'0:1400;'),(150040009,18,2,100,9,9,0,0,100,'0:1600;'),(150040010,20,2,100,10,10,0,0,100,'0:1800;'),(150040011,22,2,100,11,11,0,0,90,'0:2000;'),(150040012,24,2,100,12,12,0,0,90,'0:2200;'),(150040013,26,2,100,13,13,0,0,90,'0:2400;'),(150040014,28,2,100,14,14,0,0,90,'0:2600;'),(150040015,30,2,100,15,15,0,0,90,'0:2800;'),(150040016,32,2,100,16,16,0,0,90,'0:3000;'),(150040017,34,2,100,17,17,0,0,90,'0:3200;'),(150040018,36,2,100,18,18,0,0,90,'0:3400;'),(150040019,38,2,100,19,19,0,0,90,'0:3600;'),(150040020,40,3,100,20,20,25,0,90,'0:3800;'),(150040021,42,3,100,21,21,25,0,90,'0:4000;'),(150040022,44,3,100,22,22,25,0,90,'0:4200;'),(150040023,46,3,100,23,23,25,0,90,'0:4400;'),(150040024,48,3,100,24,24,25,0,90,'0:4600;'),(150040025,50,3,100,25,25,25,0,90,'0:4800;'),(150040026,52,3,100,26,26,25,0,90,'0:5000;'),(150040027,54,3,100,27,27,25,0,90,'0:5200;'),(150040028,56,3,100,28,28,25,0,90,'0:5400;'),(150040029,58,3,100,29,29,25,0,90,'0:5600;'),(150040030,60,3,100,30,30,25,0,80,'0:5800;'),(150040031,62,3,100,31,31,25,0,80,'0:6000;'),(150040032,64,3,100,32,32,25,0,80,'0:6200;'),(150040033,66,3,100,33,33,25,0,80,'0:6400;'),(150040034,68,3,100,34,34,25,0,80,'0:6600;'),(150040035,70,3,100,35,35,25,0,80,'0:6800;'),(150040036,72,3,100,36,36,25,0,80,'0:7000;'),(150040037,74,3,100,37,37,25,0,80,'0:7200;'),(150040038,76,3,100,38,38,25,0,80,'0:7400;'),(150040039,78,3,100,39,39,25,0,80,'0:7600;'),(150040040,80,3,100,40,40,25,0,80,'0:7800;'),(150040041,82,3,100,41,41,25,0,80,'0:8000;'),(150040042,84,3,100,42,42,25,0,80,'0:8200;'),(150040043,86,3,100,43,43,25,0,80,'0:8400;'),(150040044,88,3,100,44,44,25,0,80,'0:8600;'),(150040045,90,3,100,45,45,25,0,80,'0:8800;'),(150040046,92,3,100,46,46,50,0,80,'3:120;'),(150040047,94,3,100,47,47,50,0,80,'3:120;'),(150040048,96,3,100,48,48,50,0,80,'3:120;'),(150040049,98,3,100,49,49,50,0,80,'3:120;'),(150040050,100,3,100,50,50,50,0,70,'3:120;'),(150040051,102,4,100,51,51,50,0,70,'3:120;'),(150040052,104,4,100,52,52,50,0,70,'3:120;'),(150040053,106,4,100,53,53,50,0,70,'3:120;'),(150040054,108,4,100,54,54,50,0,70,'3:120;'),(150040055,110,4,100,55,55,50,0,70,'3:120;'),(150040056,112,4,100,56,56,50,0,70,'3:120;'),(150040057,114,4,100,57,57,50,0,70,'3:120;'),(150040058,116,4,100,58,58,50,0,70,'3:120;'),(150040059,118,4,100,59,59,50,0,70,'3:120;'),(150040060,120,4,100,60,60,50,0,70,'3:120;'),(150050001,2,2,100,1,1,0,0,0,'3:500;'),(150050002,4,2,100,2,2,0,0,0,'0:200;'),(150050003,6,2,100,3,3,0,0,0,'0:400;'),(150050004,8,2,100,4,4,0,0,0,'0:600;'),(150050005,10,2,100,5,5,0,0,0,'0:800;'),(150050006,12,2,100,6,6,0,0,0,'0:1000;'),(150050007,14,2,100,7,7,0,0,0,'0:1200;'),(150050008,16,2,100,8,8,0,0,0,'0:1400;'),(150050009,18,2,100,9,9,0,0,0,'0:1600;'),(150050010,20,2,100,10,10,0,0,0,'0:1800;'),(150050011,22,2,100,11,11,0,0,0,'0:2000;'),(150050012,24,2,100,12,12,0,0,0,'0:2200;'),(150050013,26,2,100,13,13,0,0,0,'0:2400;'),(150050014,28,2,100,14,14,0,0,0,'0:2600;'),(150050015,30,2,100,15,15,0,0,0,'0:2800;'),(150050016,32,2,100,16,16,0,0,0,'0:3000;'),(150050017,34,2,100,17,17,0,0,0,'0:3200;'),(150050018,36,2,100,18,18,0,0,0,'0:3400;'),(150050019,38,2,100,19,19,0,0,0,'0:3600;'),(150050020,40,3,100,20,20,0,0,0,'0:3800;'),(150050021,42,3,100,21,21,0,0,0,'0:4000;'),(150050022,44,3,100,22,22,0,0,0,'0:4200;'),(150050023,46,3,100,23,23,0,0,0,'0:4400;'),(150050024,48,3,100,24,24,0,0,0,'0:4600;'),(150050025,50,3,100,25,25,0,0,0,'0:4800;'),(150050026,52,3,100,26,26,0,0,0,'0:5000;'),(150050027,54,3,100,27,27,0,0,0,'0:5200;'),(150050028,56,3,100,28,28,0,0,0,'0:5400;'),(150050029,58,3,100,29,29,0,0,0,'0:5600;'),(150050030,60,3,100,30,30,0,0,0,'0:5800;'),(150050031,62,3,100,31,31,0,0,0,'0:6000;'),(150050032,64,3,100,32,32,0,0,0,'0:6200;'),(150050033,66,3,100,33,33,0,0,0,'0:6400;'),(150050034,68,3,100,34,34,0,0,0,'0:6600;'),(150050035,70,3,100,35,35,0,0,0,'0:6800;'),(150050036,72,3,100,36,36,0,0,0,'0:7000;'),(150050037,74,3,100,37,37,0,0,0,'0:7200;'),(150050038,76,3,100,38,38,0,0,0,'0:7400;'),(150050039,78,3,100,39,39,0,0,0,'0:7600;'),(150050040,80,3,100,40,40,0,0,0,'0:7800;'),(150050041,82,3,100,41,41,0,0,0,'0:8000;'),(150050042,84,3,100,42,42,0,0,0,'0:8200;'),(150050043,86,3,100,43,43,0,0,0,'0:8400;'),(150050044,88,3,100,44,44,0,0,0,'0:8600;'),(150050045,90,3,100,45,45,0,0,0,'0:8800;'),(150050046,92,3,100,46,46,0,0,0,'3:120;'),(150050047,94,3,100,47,47,0,0,0,'3:120;'),(150050048,96,3,100,48,48,0,0,0,'3:120;'),(150050049,98,3,100,49,49,0,0,0,'3:120;'),(150050050,100,3,100,50,50,50,0,0,'3:120;');
/*!40000 ALTER TABLE `rolelevelconfig` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-02-04 18:12:36
