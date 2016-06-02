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
-- Table structure for table `announcementconfig`
--

DROP TABLE IF EXISTS `announcementconfig`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `announcementconfig` (
  `id` bigint(20) NOT NULL,
  `title` varchar(2048) NOT NULL,
  `message` varchar(2048) DEFAULT NULL,
  `description` varchar(10240) DEFAULT NULL,
  `begintime` varchar(64) DEFAULT NULL,
  `endtime` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `announcementconfig`
--

LOCK TABLES `announcementconfig` WRITE;
/*!40000 ALTER TABLE `announcementconfig` DISABLE KEYS */;
INSERT INTO `announcementconfig` VALUES (1,'公告title','公告~~~~','欢迎进入电池人世界','2015/01/13 00:00','2015/01/13 00:00'),(2,'test title2','test message 2','test description 2','2015/01/13 00:00','2015/01/13 00:00'),(3,'公告','公告','欢迎进入电池人世界！','2015/06/13 00:00','2015/07/3 00:00'),(4,'更新公告','公告','亲爱的电池人<r>《电池快跑》将在10月12日上午11点进行一次版本更新，届时将有一段时间无法登陆游戏，预计3小时。维护完成后将无法使用旧版本登陆游戏，请玩家去下载新的app体验我们的电池世界。<r>维护内容： <r>优化登陆过程，如您在登陆后发现不是原先的关卡进度，请在设置内重新绑定您的微博账号，将获取原先的关卡进度！','2015/10/10 00:00','2015/10/14 00:00'),(5,'公告','公告','《电池快跑》新版本上线，请继续体验我们的电池世界！<r>如您在游戏里发现任何问题或建议，欢迎提出：<r>电池玩家群： 271434925<r>客服QQ    ： 4008500737','2015/10/10 00:00','2015/10/27 00:00'),(6,'欢迎进入电池人世界！','公告','酷酷的电池世界等你来主宰，bug、建议统统都到碗里来！任何意见或问题收到反馈都有可能在后续版本中更新！说出你喜欢的创意，让《电池快跑》更加酷炫。<r>【联系我们】<r>电池QQ群：271434925（加群即送100钻石，提bug送钻石！）<r>微信公众号: BatteryRun<r>电池官网：  run.737.com','2015/08/27 00:00','2015/8/27 00:00'),(7,'公告','公告','《电池快跑》新版本上线，请继续体验我们的电池世界！<r>如您在游戏里发现任何问题或建议，欢迎提出：<r>电池玩家群： 271434925<r>客服QQ    ： 4008500737','2015/11/5 00:00','2015/11/5 00:00'),(13,'公告','公告','《电池快跑》现已更新，主要内容如下：<r>1、更新有礼：打开邮件领取更新礼包！<r>2、新增关卡全球排名，你能超越多少人呢?<r>3、新增“闪电“道具，通关更简单！<r>4、新增奖励：每日分享可获极品道具“闪电”。<r><r>电池官方QQ群诚邀您加入：271434925','2015/12/2 00:00','2015/12/23 00:00');
/*!40000 ALTER TABLE `announcementconfig` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-02-04 18:12:37
