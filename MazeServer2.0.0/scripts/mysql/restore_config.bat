set host=192.168.1.195
set port=3306
set dbname=brdb02
set dbuser=xiaoyao
set dbpasswd=xiaoyao
set indir=backup_config
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/advertisement.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/advertisementspace.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/announcement.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/beforegamerandom.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/goods.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/jigsawconfig.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/lottoserialnumslotconfig.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/lottoslotitem.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/lottoweight.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/mailconfig.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/missionconfig.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/pickupweight.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/prop.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/roleconfig.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/rolelevelconfig.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/runeconfig.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/signaward.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/signinactivity.sql
mysql -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% < %indir%/tipconfig.sql

