set host=192.168.1.195
set port=3306
set dbname=brdb02
set dbuser=xiaoyao
set dbpasswd=xiaoyao
set outdir=server_config

REM beforegamerandom
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select goodid as ID, weight as Weight, value as Value from %dbname%.beforegamerandom" > "%outdir%/游戏前随机商品配置信息.txt"

REM goods
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as Id, malltype as MallType, mallsubtype as MallSubType, posindex as PosIndex, discount as Discount, price as Price, iapid as IapID, items as Items, IFNULL(amountperuser, '') as AmountPerUser, IFNULL(amountperround,'') as AmountPerRound, if(bestdeal=0,'','TRUE') as BestDeal, if(tesell=0,'','TRUE') as TeSell, expireddate as ExpiredDate, if(valid=0, 'FALSE', 'TRUE')  as Valid from %dbname%.goods" > "%outdir%/商品配置信息.txt"

REM jigsaw
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as PuzzleID, items as Items,  unlockprop as UnlockProps  from %dbname%.jigsawconfig" > "%outdir%/拼图配置信息.txt"

REM lottoslotitem
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select slotid as SlotId, lottotype as LottoType,  propid as PackageId,  weight as Weight, IFNULL(stage, '') as Stage, if(valid=0,'FALSE','TRUE') as Valid from %dbname%.lottoslotitem" > "%outdir%/奖池配置信息.txt"

REM lottoweight
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select IFNULL(beginvalue,'') as BeginValue, IFNULL(endvalue,'') as EndValue,  weightlist as WeightList,  if(valid=0,'FALSE','TRUE') as Valid from %dbname%.lottoweight" > "%outdir%/系统抽奖格子权重信息.txt"

REM mailconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select mailid as MailID, title as Title, message as Message, description as Description, type as Type, IFNULL(propid,'') as PropID, starttime as StartTime, endtime as EndTime from %dbname%.mailconfig" > "%outdir%/系统基础邮件信息.txt"

REM missionconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as Id, type as Type, relatedmissions as RelatedMissions, relatedprops as RelatedProps, quotas as Quotas, rewards as Rewards, begintime as BeginTime, endtime as EndTime, if(autoCollect=0,'FALSE', 'TRUE') as AutoCollect, if(valid=0,'FALSE', 'TRUE') as Valid, IFNULL(tipid,'') as TipId, tipdesc as TipDesc, if(expiredrestart=0,'FALSE','TRUE') as ExpiredRestart, IFNULL(priority,'') as Priority from %dbname%.missionconfig" > "%outdir%/任务配置信息.txt"

REM pickupweight
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select checkpointid as CheckPointId, proptype as PropType,  propid as PropId,  weight as Weight from %dbname%.pickupweight" > "%outdir%/收集物配置信息.txt"

REM prop
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as Id, type as Type, items as Items, resolvevalue as ResolveValue, IFNULL(lottovalue,'') as LottoValue, if(valid=0,'FALSE','TRUE') as Valid from %dbname%.prop" > "%outdir%/道具配置信息.txt"

REM roleconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as RoleID, maxlevel as MaxLevel, IFNULL(jigsawid,'') as JigsawID,  if(isdefaultown=0,'FALSE','TRUE') as IsDefaultOwn from %dbname%.roleconfig" > "%outdir%/角色信息.txt"

REM rolelevelconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select rolelevelid as RoleLevelID, goldbonus as GoldBonus, scorebonus as ScoreBonus, price as Price from %dbname%.rolelevelconfig" > "%outdir%/角色升级信息.txt"

REM runeconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select propid as ID, value as Value from %dbname%.runeconfig" > "%outdir%/系统道具数值配置.txt"
												
REM signaward
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as Id, day as Day, rewards as Rewards from %dbname%.signaward" > "%outdir%/签到活动奖励.txt"

REM signinactivity
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as Id, type as Type, goalvalue as GoalValue, begintime as BeginTime, endtime as EndTime, if(autocollect=0, 'FALSE', 'TRUE') as AutoCollect, if(valid=0,'FALSE','TRUE') as Valid from %dbname%.signinactivity" > "%outdir%/签到活动.txt"

REM announcement
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as Id, title as Title, message as Message, description as Description, begintime as BeginTime, endtime as EndTime from %dbname%.announcementconfig" > "%outdir%/公告.txt"

REM lottoserialnumslotconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select serialnum as Num, proplist as PropList, selected as Selected, if(valid=0,'FALSE','TRUE') as Valid from %dbname%.lottoserialnumslotconfig" > "%outdir%/特殊抽奖配置信息.txt"

REM advertisement
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as Id, viewurl as ViewUrl, materialurl as MaterialUrl, clickurl as ClickUrl from %dbname%.advertisement" > "%outdir%/广告配置信息.txt"

REM advertisementspace
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as Id, items as Items, if(enable=0,'FALSE','TRUE') as Enable, flags as Flags from %dbname%.advertisementspace" > "%outdir%/广告位配置信息.txt"

REM tip
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as Id, language as Language, title as Title, content as Content from %dbname%.tipconfig" > "%outdir%/提示配置信息.txt"

REM leaguelevelconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select level as Level, humanamount as HumanAmount, groupamount as GroupAmount, promoteamount as PromoteAmount, demoteamount as DemoteAmount from %dbname%.leaguelevelconfig" > "%outdir%/联赛等级配置信息.txt"

REM leaguelevelawardconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select level as Level, items as Items from %dbname%.leaguelevelawardconfig" > "%outdir%/联赛等级奖励配置信息.txt"

REM leaguerankawardconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select level as Level, beginrank as BeginRank, endrank as EndRank, items as Items from %dbname%.leaguerankawardconfig" > "%outdir%/联赛排名奖励配置信息.txt"

REM leaguedefaultaiscoreconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select rank as Rank, score as Score from %dbname%.leaguedefaultaiscoreconfig" > "%outdir%/联赛默认ai分数配置信息.txt"

REM  newaccountprop
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select  source as Source, propitems as PropItems, type as Type, mailid as MailID from %dbname%.newaccountprop" > "%outdir%/首次登录礼包配置信息.txt"

REM 备份数据库
backup_config.bat