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
-- Table structure for table `signaward`
--

DROP TABLE IF EXISTS `signaward`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `signaward` (
  `_id` int(11) NOT NULL AUTO_INCREMENT,
  `id` int(11) NOT NULL,
  `day` int(11) NOT NULL,
  `rewards` varchar(1024) DEFAULT NULL,
  PRIMARY KEY (`_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `signaward`
--

LOCK TABLES `signaward` WRITE;
/*!40000 ALTER TABLE `signaward` DISABLE KEYS */;
INSERT INTO `signaward` VALUES (1,1,0,'147010001:1:14'),(2,1,1,'160020000:10:16'),(3,1,2,'20010000:5:2'),(4,1,3,'40010000:20:4'),(5,1,4,'10010000:5000:1'),(6,1,5,'160030000:10:16'),(7,1,6,'10010000:8000:1'),(8,1,7,'160050000:10:16');
/*!40000 ALTER TABLE `signaward` ENABLE KEYS */;
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
