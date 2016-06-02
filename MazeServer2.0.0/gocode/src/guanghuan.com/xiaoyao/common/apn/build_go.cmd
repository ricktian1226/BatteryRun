set PROTOC=E:\svn_debug\Savage\Program\MazeNetMsgDef2.0.0\protoc
set SOURCE=.\
set TARGET=E:\svn_debug\Savage\Program\MazeServer2.0.0\gocode\src\guanghuan.com\xiaoyao\common\apn
%PROTOC% --go_out=%SOURCE% --proto_path=. .\*.proto
copy /y %SOURCE%\*.* %TARGET%