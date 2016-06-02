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
-- Table structure for table `mailconfig`
--

DROP TABLE IF EXISTS `mailconfig`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mailconfig` (
  `mailid` int(11) NOT NULL,
  `title` varchar(512) DEFAULT NULL,
  `message` varchar(512) DEFAULT NULL,
  `description` varchar(2048) DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `propid` bigint(20) DEFAULT NULL,
  `starttime` varchar(128) DEFAULT NULL,
  `endtime` varchar(128) DEFAULT NULL,
  `remark` varchar(1024) DEFAULT NULL,
  PRIMARY KEY (`mailid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `mailconfig`
--

LOCK TABLES `mailconfig` WRITE;
/*!40000 ALTER TABLE `mailconfig` DISABLE KEYS */;
INSERT INTO `mailconfig` VALUES (111001,'登陆奖励','第一封test！','喜迎上线，钻石免费送，只要你敢来，我就敢送！<R><R>登陆就送100钻！',3,140160000,'2014/11/10 00:00','2014/11/10 00:00',''),(500001,'活动公告','首次充值双倍钻石！','首次充值双倍钻石！<R><R>冲的多，送的多！<R><R>冲200送500！<R><R>充500送1300！<R><R>心动了吗？钻石简直不要钱！',1,0,'2014/11/1 15:00','2014/11/1 15:00',''),(500002,'活动公告','连续登陆有奖！','连续登陆有奖哟！<R><R>记得每天上来看看！',1,0,'2014/11/1 15:00','2014/11/1 15:00',''),(500003,'版本更新通知','2015.10.12版本更新','亲爱的电池人<r>《电池快跑》将在10月12日上午11点进行一次版本更新，届时将有一段时间无法登陆游戏，预计3小时。维护完成后将无法使用旧版本登陆游戏，请玩家去下载新的app体验我们的电池世界。<r>维护内容： <r>1、优化登陆过程，如您在登陆后发现不是原先的关卡进度，请在设置内重新绑定您的微博账号，将获取原先的关卡进度！<r>2、调整部分关卡难度',2,0,'2015/10/10 00:00','2015/10/14 00:00',''),(500004,'版本更新礼包','2015.12.3版本更新','欢迎来到新的电池世界！里面有更精彩的内容在等着你！',3,140400000,'2015/12/3 00:00','2015/12/23 00:00',''),(800000,'系统发放','系统赏给你的','系统赏给你的',3,0,'','','运营专用，勿修改'),(900000,'微博绑定礼包','微博绑定奖励','微博绑定奖励',3,0,'2015/9/1 00:00','2025/9/1 00:00','[900000,1000000)为微博登录奖励动态邮件专用'),(1000000,'新玩家登录礼包','新玩家登录奖励','新玩家登录奖励',3,0,'2015/9/1 00:00','2025/9/1 00:00','[1000000,1100000)为游客登录奖励动态邮件专用');
/*!40000 ALTER TABLE `mailconfig` ENABLE KEYS */;
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
