set host=192.168.1.195
set port=3306
set dbname=brdb02
set dbuser=xiaoyao
set dbpasswd=xiaoyao
set outdir=client_config

REM beforegamerandom
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select goodid as ID, value as Value from %dbname%.beforegamerandom" > "%outdir%/��Ϸǰ�����Ʒ������Ϣ.txt"

REM goods
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as Id, malltype as MallType, mallsubtype as MallSubType, posindex as PosIndex, discount as Discount, price as Price, iapid as IapID, items as Items, IFNULL(amountperuser, '') as AmountPerUser, IFNULL(amountperround,'') as AmountPerRound, if(bestdeal=0,'','TRUE') as BestDeal, if(tesell=0,'','TRUE') as TeSell, expireddate as ExpiredDate from %dbname%.goods" > "%outdir%/��Ʒ������Ϣ.txt"

REM jigsaw
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as PuzzleID, items as Items, unlockprop as UnlockProps, resname as ResName from %dbname%.jigsawconfig" > "%outdir%/ƴͼ������Ϣ.txt"

REM missionconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select tipid as TipId, tipdesc as TipDesc from %dbname%.missionconfig" > "%outdir%/����������Ϣ.txt"

REM prop
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as Id, nameid as NameID, iconresname as IconResName, type as Type, IFNULL(descriptionid, '') as DescriptionID, items as Items, resolvevalue as ResolveValue from %dbname%.prop" > "%outdir%/����������Ϣ.txt"

REM roleconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select id as RoleID, name as RoleName, talkcontent as TalkContent, entityid as EntityID, skilldescription as SkillDescribe, fullleveldescription as FullLevelDescribe, maxlevel as MaxLevel, IFNULL(jigsawid,'') as JigsawID from %dbname%.roleconfig" > "%outdir%/��ɫ��Ϣ.txt"

REM rolelevelconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select rolelevelid as RoleLevelID, IFNULL(propbonus, '') as PropBonus, IFNULL(hp, '') as Hp, IFNULL(helpbonus, '') as HelpBonus, IFNULL(goldbonus, '') as GoldBonus, IFNULL(scorebonus, '') as ScoreBonus, IFNULL(skillBonus, '') as SkillBonus, IFNULL(skillcasttime, '') as SkillCastTime, IFNULL(skillColdTime, '') as SkillColdTime, price as Price from %dbname%.rolelevelconfig" > "%outdir%/��ɫ������Ϣ.txt"

REM runeconfig
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% -e "select propid as ID, value as Value from %dbname%.runeconfig" > "%outdir%/ϵͳ������ֵ����.txt"
