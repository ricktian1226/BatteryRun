set host=192.168.1.195
set port=3306
set dbname=brdb02
set dbuser=xiaoyao
set dbpasswd=xiaoyao
set outdir=backup_config

mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% beforegamerandom > "%outdir%/beforegamerandom.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% goods > "%outdir%/goods.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% jigsawconfig > "%outdir%/jigsawconfig.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% lottoserialnumslotconfig > "%outdir%/lottoserialnumslotconfig.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% lottoslotitem > "%outdir%/lottoslotitem.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% lottoweight > "%outdir%/lottoweight.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% mailconfig > "%outdir%/mailconfig.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% missionconfig > "%outdir%/missionconfig.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% pickupweight > "%outdir%/pickupweight.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% prop > "%outdir%/prop.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% roleconfig > "%outdir%/roleconfig.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% rolelevelconfig > "%outdir%/rolelevelconfig.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% runeconfig > "%outdir%/runeconfig.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% signaward > "%outdir%/signaward.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% signinactivity > "%outdir%/signinactivity.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% advertisement > "%outdir%/advertisement.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% advertisementspace > "%outdir%/advertisementspace.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% announcementconfig > "%outdir%/announcement.sql"
mysqldump -h%host% -P%port% -u%dbuser% -p%dbpasswd% %dbname% tipconfig > "%outdir%/tipconfig.sql"

