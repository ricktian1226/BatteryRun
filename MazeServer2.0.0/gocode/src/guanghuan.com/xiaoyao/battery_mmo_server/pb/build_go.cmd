set PROTOC=..\protoc
set SOURCE=.\
set TARGET=E:\svn_debug\Savage\Program\MazeServer\gocode\src\guanghuan.com\xiaoyao\battery_mmo_server\pb
%PROTOC% --go_out=%SOURCE% --proto_path=. .\*.proto
copy /y %SOURCE%\*.* %TARGET%